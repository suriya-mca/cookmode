package queries

import (
	"context"

	"github.com/jackc/pgx/v5"

	"cookmode/internal/models"
)

func scanUser(row pgx.Row) (*models.User, error) {
	var u models.User
	var username, displayName, avatarURL, bio *string
	err := row.Scan(&u.ID, &username, &displayName, &avatarURL, &bio, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if username != nil {
		u.Username = *username
	}
	if displayName != nil {
		u.DisplayName = *displayName
	}
	if avatarURL != nil {
		u.AvatarURL = *avatarURL
	}
	if bio != nil {
		u.Bio = *bio
	}
	return &u, nil
}

func (db *DB) GetUser(ctx context.Context, id string) (*models.User, error) {
	row := db.Pool.QueryRow(ctx, `SELECT id, username, display_name, avatar_url, bio, created_at, updated_at FROM users WHERE id=$1`, id)
	return scanUser(row)
}

func (db *DB) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	row := db.Pool.QueryRow(ctx, `SELECT id, username, display_name, avatar_url, bio, created_at, updated_at FROM users WHERE username=$1`, username)
	return scanUser(row)
}

func (db *DB) PatchUser(ctx context.Context, id string, username, displayName, avatarURL, bio *string) (*models.User, error) {
	row := db.Pool.QueryRow(ctx, `
		INSERT INTO users (id, username, display_name, avatar_url, bio)
		VALUES ($1, NULLIF($2::text, ''), COALESCE($3::text, ''), COALESCE($4::text, ''), COALESCE($5::text, ''))
		ON CONFLICT (id) DO UPDATE SET
			username = CASE WHEN $6 THEN EXCLUDED.username ELSE users.username END,
			display_name = CASE WHEN $7 THEN EXCLUDED.display_name ELSE users.display_name END,
			avatar_url = CASE WHEN $8 THEN EXCLUDED.avatar_url ELSE users.avatar_url END,
			bio = CASE WHEN $9 THEN EXCLUDED.bio ELSE users.bio END
		RETURNING id, username, display_name, avatar_url, bio, created_at, updated_at
	`, id, username, displayName, avatarURL, bio, username != nil, displayName != nil, avatarURL != nil, bio != nil)
	return scanUser(row)
}

func (db *DB) UpdateUser(ctx context.Context, u *models.User) (*models.User, error) {
	row := db.Pool.QueryRow(ctx, `
		UPDATE users SET username=$2, display_name=$3, avatar_url=$4, bio=$5
		WHERE id=$1
		RETURNING id, username, display_name, avatar_url, bio, created_at, updated_at
	`, u.ID, nilIfEmpty(u.Username), u.DisplayName, u.AvatarURL, u.Bio)
	return scanUser(row)
}

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
