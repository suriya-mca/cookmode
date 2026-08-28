package queries

import (
	"context"
	"encoding/json"

	"cookmode/internal/models"
)

func (db *DB) ForkRecipe(ctx context.Context, parentID, userID string) (*models.Recipe, error) {
	if err := db.ensureUser(ctx, userID); err != nil {
		return nil, err
	}
	parent, err := db.GetRecipe(ctx, parentID)
	if err != nil {
		return nil, err
	}
	ingredients, _ := json.Marshal(parent.Ingredients)
	steps, _ := json.Marshal(parent.Steps)
	substitutions, _ := json.Marshal(parent.Substitutions)
	var nutrition []byte
	if parent.Nutrition != nil {
		nutrition, _ = json.Marshal(parent.Nutrition)
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
	if parent.DietaryTags == nil {
		parent.DietaryTags = []string{}
	}
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	row := tx.QueryRow(ctx, `
		INSERT INTO recipes (user_id, title, description, cuisine, prep_time_min, cook_time_min, servings, difficulty, dietary_tags, ingredients, steps, substitutions, nutrition, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,'draft')
		RETURNING id, user_id, title, description, cuisine, prep_time_min, cook_time_min, servings, difficulty, dietary_tags, ingredients, steps, substitutions, nutrition, video_uid, video_hls_url, video_thumbnail_url, video_duration_sec, views, saves, status, created_at, updated_at
	`, userID, parent.Title, parent.Description, parent.Cuisine, parent.PrepTimeMin, parent.CookTimeMin, parent.Servings, string(parent.Difficulty), parent.DietaryTags, ingredients, steps, substitutions, nutrition)
	child, err := scanRecipe(row)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(ctx, `INSERT INTO recipe_forks (parent_recipe_id, child_recipe_id, user_id) VALUES ($1,$2,$3)`, parentID, child.ID, userID)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return child, nil
}
