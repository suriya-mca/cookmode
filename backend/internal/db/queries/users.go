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

func (db *DB) GetUserByUsername(ctx context.Context, username string) (*models.User, string, error) {
	var u models.User
	var dbUsername, displayName, avatarURL, bio, hash *string
	err := db.Pool.QueryRow(ctx, `SELECT id, username, display_name, avatar_url, bio, password_hash, created_at, updated_at FROM users WHERE username=$1`, username).Scan(
		&u.ID, &dbUsername, &displayName, &avatarURL, &bio, &hash, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, "", err
	}
	if dbUsername != nil {
		u.Username = *dbUsername
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
	var h string
	if hash != nil {
		h = *hash
	}
	return &u, h, nil
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

func (db *DB) CreateAuthUser(ctx context.Context, id, username, hash, displayName string) (*models.User, error) {
	var u models.User
	var dbUsername, dbDisplayName, dbAvatarURL, dbBio *string
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO users (id, username, password_hash, display_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, display_name, avatar_url, bio, created_at, updated_at
	`, id, username, hash, displayName).Scan(&u.ID, &dbUsername, &dbDisplayName, &dbAvatarURL, &dbBio, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrUsernameTaken
		}
		return nil, err
	}
	if dbUsername != nil {
		u.Username = *dbUsername
	}
	if dbDisplayName != nil {
		u.DisplayName = *dbDisplayName
	}
	if dbAvatarURL != nil {
		u.AvatarURL = *dbAvatarURL
	}
	if dbBio != nil {
		u.Bio = *dbBio
	}
	return &u, nil
}

var ErrUsernameTaken = func() error {
	return &usernameTakenError{}
}()

type usernameTakenError struct{}

func (e *usernameTakenError) Error() string { return "username already taken" }

func isUniqueViolation(err error) bool {
	return err != nil && (contains(err.Error(), "duplicate key") || contains(err.Error(), "unique"))
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool {
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()
}
