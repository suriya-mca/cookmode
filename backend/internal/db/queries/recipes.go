package queries

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"cookmode/internal/models"
)

func encodeCursor(t time.Time, id string) string {
	raw := fmt.Sprintf("%s|%s", t.Format(time.RFC3339Nano), id)
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

func decodeCursor(cursor string) (time.Time, string, error) {
	b, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return time.Time{}, "", err
	}
	parts := string(b)
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == '|' {
			ts, err := time.Parse(time.RFC3339Nano, parts[:i])
			if err != nil {
				return time.Time{}, "", err
			}
			return ts, parts[i+1:], nil
		}
	}
	return time.Time{}, "", fmt.Errorf("invalid cursor")
}

func scanRecipe(row pgx.Row) (*models.Recipe, error) {
	var r models.Recipe
	var ingredients, steps, substitutions, nutrition []byte
	var dietaryTags []string
	var createdAt, updatedAt time.Time
	err := row.Scan(
		&r.ID, &r.UserID, &r.Title, &r.Description, &r.Cuisine,
		&r.PrepTimeMin, &r.CookTimeMin, &r.Servings, &r.Difficulty,
		&dietaryTags,
		&ingredients, &steps, &substitutions, &nutrition,
		&r.VideoUID, &r.VideoHLSURL, &r.VideoThumbnailURL, &r.VideoDurationSec,
		&r.Views, &r.Saves, &r.Status,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	r.DietaryTags = dietaryTags
	r.CreatedAt = createdAt
	r.UpdatedAt = updatedAt
	if err := json.Unmarshal(ingredients, &r.Ingredients); err != nil {
		return nil, fmt.Errorf("unmarshal ingredients: %w", err)
	}
	if err := json.Unmarshal(steps, &r.Steps); err != nil {
		return nil, fmt.Errorf("unmarshal steps: %w", err)
	}
	if err := json.Unmarshal(substitutions, &r.Substitutions); err != nil {
		return nil, fmt.Errorf("unmarshal substitutions: %w", err)
	}
	if len(nutrition) > 0 && string(nutrition) != "null" {
		var n models.Nutrition
		if err := json.Unmarshal(nutrition, &n); err != nil {
			return nil, fmt.Errorf("unmarshal nutrition: %w", err)
		}
		r.Nutrition = &n
	}
	if r.Ingredients == nil {
		r.Ingredients = []models.Ingredient{}
	}
	if r.Steps == nil {
		r.Steps = []models.Step{}
	}
	if r.Substitutions == nil {
		r.Substitutions = []models.Substitution{}
	}
	if r.DietaryTags == nil {
		r.DietaryTags = []string{}
	}
	return &r, nil
}

func (db *DB) CreateRecipe(ctx context.Context, r *models.Recipe) (*models.Recipe, error) {
	ingredients, _ := json.Marshal(r.Ingredients)
	steps, _ := json.Marshal(r.Steps)
	substitutions, _ := json.Marshal(r.Substitutions)
	var nutrition []byte
	if r.Nutrition != nil {
		nutrition, _ = json.Marshal(r.Nutrition)
	}
	if ingredients == nil {
		ingredients = []byte("[]")
	}
	if steps == nil {
		steps = []byte("[]")
	}
	if substitutions == nil {
		substitutions = []byte("[]")
	}
	if r.DietaryTags == nil {
		r.DietaryTags = []string{}
	}
	if r.Servings == 0 {
		r.Servings = 1
	}
	if r.Status == "" {
		r.Status = models.StatusDraft
	}
	row := db.Pool.QueryRow(ctx, `
		INSERT INTO recipes (user_id, title, description, cuisine, prep_time_min, cook_time_min, servings, difficulty, dietary_tags, ingredients, steps, substitutions, nutrition, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING id, user_id, title, description, cuisine, prep_time_min, cook_time_min, servings, difficulty, dietary_tags, ingredients, steps, substitutions, nutrition, video_uid, video_hls_url, video_thumbnail_url, video_duration_sec, views, saves, status, created_at, updated_at
	`, r.UserID, r.Title, r.Description, r.Cuisine, r.PrepTimeMin, r.CookTimeMin, r.Servings, string(r.Difficulty), r.DietaryTags, ingredients, steps, substitutions, nutrition, string(r.Status))
	return scanRecipe(row)
}

func (db *DB) GetRecipe(ctx context.Context, id string) (*models.Recipe, error) {
	row := db.Pool.QueryRow(ctx, `
		SELECT id, user_id, title, description, cuisine, prep_time_min, cook_time_min, servings, difficulty, dietary_tags, ingredients, steps, substitutions, nutrition, video_uid, video_hls_url, video_thumbnail_url, video_duration_sec, views, saves, status, created_at, updated_at
		FROM recipes WHERE id = $1
	`, id)
	return scanRecipe(row)
}

func (db *DB) ListRecipes(ctx context.Context, limit int, cursor string) ([]*models.Recipe, string, error) {
	var rows pgx.Rows
	var err error
	if cursor == "" {
		rows, err = db.Pool.Query(ctx, `
			SELECT id, user_id, title, description, cuisine, prep_time_min, cook_time_min, servings, difficulty, dietary_tags, ingredients, steps, substitutions, nutrition, video_uid, video_hls_url, video_thumbnail_url, video_duration_sec, views, saves, status, created_at, updated_at
			FROM recipes WHERE status = 'published'
			ORDER BY created_at DESC, id DESC LIMIT $1
		`, limit+1)
	} else {
		ts, cid, err := decodeCursor(cursor)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor")
		}
		rows, err = db.Pool.Query(ctx, `
			SELECT id, user_id, title, description, cuisine, prep_time_min, cook_time_min, servings, difficulty, dietary_tags, ingredients, steps, substitutions, nutrition, video_uid, video_hls_url, video_thumbnail_url, video_duration_sec, views, saves, status, created_at, updated_at
			FROM recipes WHERE status = 'published' AND (created_at, id) < ($1, $2)
			ORDER BY created_at DESC, id DESC LIMIT $3
		`, ts, cid, limit+1)
	}
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	var out []*models.Recipe
	for rows.Next() {
		r, err := scanRecipe(rows)
		if err != nil {
			return nil, "", err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, "", err
	}
	var nextCursor string
	if len(out) > limit {
		last := out[limit-1]
		nextCursor = encodeCursor(last.CreatedAt, last.ID)
		out = out[:limit]
	}
	return out, nextCursor, nil
}

func (db *DB) UpdateRecipe(ctx context.Context, r *models.Recipe) (*models.Recipe, error) {
	ingredients, _ := json.Marshal(r.Ingredients)
	steps, _ := json.Marshal(r.Steps)
	substitutions, _ := json.Marshal(r.Substitutions)
	var nutrition []byte
	if r.Nutrition != nil {
		nutrition, _ = json.Marshal(r.Nutrition)
	}
	if ingredients == nil {
		ingredients = []byte("[]")
	}
	if steps == nil {
		steps = []byte("[]")
	}
	if substitutions == nil {
		substitutions = []byte("[]")
	}
	if r.DietaryTags == nil {
		r.DietaryTags = []string{}
	}
	row := db.Pool.QueryRow(ctx, `
		UPDATE recipes SET title=$1, description=$2, cuisine=$3, prep_time_min=$4, cook_time_min=$5, servings=$6, difficulty=$7, dietary_tags=$8, ingredients=$9, steps=$10, substitutions=$11, nutrition=$12
		WHERE id=$13
		RETURNING id, user_id, title, description, cuisine, prep_time_min, cook_time_min, servings, difficulty, dietary_tags, ingredients, steps, substitutions, nutrition, video_uid, video_hls_url, video_thumbnail_url, video_duration_sec, views, saves, status, created_at, updated_at
	`, r.Title, r.Description, r.Cuisine, r.PrepTimeMin, r.CookTimeMin, r.Servings, string(r.Difficulty), r.DietaryTags, ingredients, steps, substitutions, nutrition, r.ID)
	return scanRecipe(row)
}

func (db *DB) ArchiveRecipe(ctx context.Context, id string) (*models.Recipe, error) {
	row := db.Pool.QueryRow(ctx, `
		UPDATE recipes SET status='archived' WHERE id=$1
		RETURNING id, user_id, title, description, cuisine, prep_time_min, cook_time_min, servings, difficulty, dietary_tags, ingredients, steps, substitutions, nutrition, video_uid, video_hls_url, video_thumbnail_url, video_duration_sec, views, saves, status, created_at, updated_at
	`, id)
	return scanRecipe(row)
}

func (db *DB) SetRecipeVideo(ctx context.Context, id, videoUID, status string) (*models.Recipe, error) {
	row := db.Pool.QueryRow(ctx, `
		UPDATE recipes SET video_uid=$1, status=$2 WHERE id=$3
		RETURNING id, user_id, title, description, cuisine, prep_time_min, cook_time_min, servings, difficulty, dietary_tags, ingredients, steps, substitutions, nutrition, video_uid, video_hls_url, video_thumbnail_url, video_duration_sec, views, saves, status, created_at, updated_at
	`, videoUID, status, id)
	return scanRecipe(row)
}
