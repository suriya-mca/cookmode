package queries

import (
	"context"

	"github.com/jackc/pgx/v5"

	"cookmode/internal/models"
)

// scanPost scans a database row into a post and returns the post or a scan error.
func scanPost(row pgx.Row) (*models.Post, error) {
	var p models.Post
	err := row.Scan(&p.ID, &p.UserID, &p.RecipeID, &p.PhotoURL, &p.Caption, &p.Status, &p.RatingValue, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (db *DB) CreatePost(ctx context.Context, p *models.Post) (*models.Post, error) {
	if err := db.ensureUser(ctx, p.UserID); err != nil {
		return nil, err
	}
	row := db.Pool.QueryRow(ctx, `
		INSERT INTO posts (user_id, recipe_id, photo_url, caption, rating_value)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, user_id, recipe_id, photo_url, caption, status, rating_value, created_at
	`, p.UserID, p.RecipeID, p.PhotoURL, p.Caption, p.RatingValue)
	return scanPost(row)
}

func (db *DB) GetPost(ctx context.Context, id string) (*models.Post, error) {
	row := db.Pool.QueryRow(ctx, `SELECT id, user_id, recipe_id, photo_url, caption, status, rating_value, created_at FROM posts WHERE id=$1`, id)
	return scanPost(row)
}

func (db *DB) ListPosts(ctx context.Context, recipeID, userID string) ([]*models.Post, error) {
	var rows pgx.Rows
	var err error
	if recipeID != "" {
		rows, err = db.Pool.Query(ctx, `SELECT id, user_id, recipe_id, photo_url, caption, status, rating_value, created_at FROM posts WHERE recipe_id=$1 AND status='visible' ORDER BY created_at DESC`, recipeID)
	} else if userID != "" {
		rows, err = db.Pool.Query(ctx, `SELECT id, user_id, recipe_id, photo_url, caption, status, rating_value, created_at FROM posts WHERE user_id=$1 AND status='visible' ORDER BY created_at DESC`, userID)
	} else {
		rows, err = db.Pool.Query(ctx, `SELECT id, user_id, recipe_id, photo_url, caption, status, rating_value, created_at FROM posts WHERE status='visible' ORDER BY created_at DESC LIMIT 50`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Post
	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (db *DB) DeletePost(ctx context.Context, id, userID string) error {
	res, err := db.Pool.Exec(ctx, `DELETE FROM posts WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
