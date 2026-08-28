package queries

import "context"

func (db *DB) SaveRecipe(ctx context.Context, userID, recipeID string) error {
	if err := db.ensureUser(ctx, userID); err != nil {
		return err
	}
	_, err := db.Pool.Exec(ctx, `INSERT INTO saves (user_id, recipe_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, userID, recipeID)
	if err != nil {
		return err
	}
	_, _ = db.Pool.Exec(ctx, `UPDATE recipes SET saves = (SELECT COUNT(*) FROM saves WHERE recipe_id=$1) WHERE id=$1`, recipeID)
	return nil
}

func (db *DB) UnsaveRecipe(ctx context.Context, userID, recipeID string) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM saves WHERE user_id=$1 AND recipe_id=$2`, userID, recipeID)
	if err != nil {
		return err
	}
	_, _ = db.Pool.Exec(ctx, `UPDATE recipes SET saves = (SELECT COUNT(*) FROM saves WHERE recipe_id=$1) WHERE id=$1`, recipeID)
	return nil
}

func (db *DB) IsSaved(ctx context.Context, userID, recipeID string) (bool, error) {
	var exists bool
	err := db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM saves WHERE user_id=$1 AND recipe_id=$2)`, userID, recipeID).Scan(&exists)
	return exists, err
}
