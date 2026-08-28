package queries

import "context"

func (db *DB) SaveRecipe(ctx context.Context, userID, recipeID string) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `INSERT INTO users (id) VALUES ($1) ON CONFLICT DO NOTHING`, userID); err != nil {
		return err
	}
	tag, err := tx.Exec(ctx, `INSERT INTO saves (user_id, recipe_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, userID, recipeID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 1 {
		if _, err := tx.Exec(ctx, `UPDATE recipes SET saves = saves + 1 WHERE id=$1`, recipeID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (db *DB) UnsaveRecipe(ctx context.Context, userID, recipeID string) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	tag, err := tx.Exec(ctx, `DELETE FROM saves WHERE user_id=$1 AND recipe_id=$2`, userID, recipeID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 1 {
		if _, err := tx.Exec(ctx, `UPDATE recipes SET saves = GREATEST(saves - 1, 0) WHERE id=$1`, recipeID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (db *DB) IsSaved(ctx context.Context, userID, recipeID string) (bool, error) {
	var exists bool
	err := db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM saves WHERE user_id=$1 AND recipe_id=$2)`, userID, recipeID).Scan(&exists)
	return exists, err
}
