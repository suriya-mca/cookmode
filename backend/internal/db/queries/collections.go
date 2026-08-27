package queries

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"cookmode/internal/models"
)

func scanCollection(row pgx.Row) (*models.Collection, error) {
	var c models.Collection
	err := row.Scan(&c.ID, &c.UserID, &c.Title, &c.Description, &c.IsPublic, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (db *DB) CreateCollection(ctx context.Context, c *models.Collection) (*models.Collection, error) {
	if err := db.ensureUser(ctx, c.UserID); err != nil {
		return nil, err
	}
	row := db.Pool.QueryRow(ctx, `
		INSERT INTO collections (user_id, title, description, is_public)
		VALUES ($1,$2,$3,$4)
		RETURNING id, user_id, title, description, is_public, created_at, updated_at
	`, c.UserID, c.Title, c.Description, c.IsPublic)
	return scanCollection(row)
}

func (db *DB) GetCollection(ctx context.Context, id string) (*models.Collection, error) {
	row := db.Pool.QueryRow(ctx, `SELECT id, user_id, title, description, is_public, created_at, updated_at FROM collections WHERE id=$1`, id)
	return scanCollection(row)
}

func (db *DB) ListCollections(ctx context.Context, userID string) ([]*models.Collection, error) {
	rows, err := db.Pool.Query(ctx, `SELECT id, user_id, title, description, is_public, created_at, updated_at FROM collections WHERE user_id=$1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Collection
	for rows.Next() {
		c, err := scanCollection(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (db *DB) UpdateCollection(ctx context.Context, c *models.Collection) (*models.Collection, error) {
	row := db.Pool.QueryRow(ctx, `
		UPDATE collections SET title=$1, description=$2, is_public=$3 WHERE id=$4
		RETURNING id, user_id, title, description, is_public, created_at, updated_at
	`, c.Title, c.Description, c.IsPublic, c.ID)
	return scanCollection(row)
}

func (db *DB) DeleteCollection(ctx context.Context, id string) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM collections WHERE id=$1`, id)
	return err
}

func (db *DB) AddRecipeToCollection(ctx context.Context, collectionID, recipeID string) error {
	_, err := db.Pool.Exec(ctx, `INSERT INTO collection_recipes (collection_id, recipe_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, collectionID, recipeID)
	return err
}

func (db *DB) RemoveRecipeFromCollection(ctx context.Context, collectionID, recipeID string) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM collection_recipes WHERE collection_id=$1 AND recipe_id=$2`, collectionID, recipeID)
	return err
}

func (db *DB) ListCollectionRecipes(ctx context.Context, collectionID string) ([]*models.Recipe, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT r.id, r.user_id, r.title, r.description, r.cuisine, r.prep_time_min, r.cook_time_min, r.servings, r.difficulty, r.dietary_tags, r.ingredients, r.steps, r.substitutions, r.nutrition, r.video_uid, r.video_hls_url, r.video_thumbnail_url, r.video_duration_sec, r.views, r.saves, r.status, r.created_at, r.updated_at
		FROM recipes r JOIN collection_recipes cr ON cr.recipe_id = r.id
		WHERE cr.collection_id=$1 ORDER BY cr.added_at DESC
	`, collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Recipe
	for rows.Next() {
		r, err := scanRecipe(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

var _ = time.Now
