package queries

import "context"

func (db *DB) Follow(ctx context.Context, followerID, followingID string) error {
	if err := db.ensureUser(ctx, followerID); err != nil {
		return err
	}
	if err := db.ensureUser(ctx, followingID); err != nil {
		return err
	}
	_, err := db.Pool.Exec(ctx, `INSERT INTO follows (follower_id, following_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, followerID, followingID)
	return err
}

func (db *DB) Unfollow(ctx context.Context, followerID, followingID string) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM follows WHERE follower_id=$1 AND following_id=$2`, followerID, followingID)
	return err
}

func (db *DB) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	var exists bool
	err := db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id=$1 AND following_id=$2)`, followerID, followingID).Scan(&exists)
	return exists, err
}

func (db *DB) ListFollowers(ctx context.Context, userID string) ([]string, error) {
	rows, err := db.Pool.Query(ctx, `SELECT follower_id FROM follows WHERE following_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func (db *DB) ListFollowing(ctx context.Context, userID string) ([]string, error) {
	rows, err := db.Pool.Query(ctx, `SELECT following_id FROM follows WHERE follower_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}
