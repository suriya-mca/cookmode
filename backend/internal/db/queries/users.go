package queries

import (
	"context"

	"github.com/jackc/pgx/v5"

	"cookmode/internal/models"
)

// scanUser converts a database row into a user model, preserving zero values for nullable fields.
// It returns any row scanning error unchanged.
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

func (db *DB) UpsertUser(ctx context.Context, u *models.User) (*models.User, error) {
	row := db.Pool.QueryRow(ctx, `
		INSERT INTO users (id, username, display_name, avatar_url, bio)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (id) DO UPDATE SET username=EXCLUDED.username, display_name=EXCLUDED.display_name, avatar_url=EXCLUDED.avatar_url, bio=EXCLUDED.bio
		RETURNING id, username, display_name, avatar_url, bio, created_at, updated_at
	`, u.ID, nilIfEmpty(u.Username), u.DisplayName, u.AvatarURL, u.Bio)
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

// nilIfEmpty returns nil for an empty string and the original string otherwise.
func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
