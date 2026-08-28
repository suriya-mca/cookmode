package queries

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"cookmode/internal/models"
)

func scanShoppingList(row pgx.Row) (*models.ShoppingList, error) {
	var l models.ShoppingList
	var items []byte
	err := row.Scan(&l.ID, &l.UserID, &l.Title, &items, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(items, &l.Items); err != nil {
		return nil, err
	}
	if l.Items == nil {
		l.Items = []models.ShoppingListItem{}
	}
	return &l, nil
}

func (db *DB) CreateShoppingList(ctx context.Context, l *models.ShoppingList) (*models.ShoppingList, error) {
	if err := db.ensureUser(ctx, l.UserID); err != nil {
		return nil, err
	}
	items, _ := json.Marshal(l.Items)
	if items == nil {
		items = []byte("[]")
	}
	row := db.Pool.QueryRow(ctx, `
		INSERT INTO shopping_lists (user_id, title, items)
		VALUES ($1,$2,$3)
		RETURNING id, user_id, title, items, created_at, updated_at
	`, l.UserID, l.Title, items)
	return scanShoppingList(row)
}

func (db *DB) GetShoppingList(ctx context.Context, id, userID string) (*models.ShoppingList, error) {
	row := db.Pool.QueryRow(ctx, `SELECT id, user_id, title, items, created_at, updated_at FROM shopping_lists WHERE id=$1 AND user_id=$2`, id, userID)
	return scanShoppingList(row)
}

func (db *DB) ListShoppingLists(ctx context.Context, userID string) ([]*models.ShoppingList, error) {
	rows, err := db.Pool.Query(ctx, `SELECT id, user_id, title, items, created_at, updated_at FROM shopping_lists WHERE user_id=$1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.ShoppingList
	for rows.Next() {
		l, err := scanShoppingList(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (db *DB) UpdateShoppingListItems(ctx context.Context, id, userID string, items []models.ShoppingListItem) (*models.ShoppingList, error) {
	b, _ := json.Marshal(items)
	row := db.Pool.QueryRow(ctx, `UPDATE shopping_lists SET items=$1 WHERE id=$2 AND user_id=$3 RETURNING id, user_id, title, items, created_at, updated_at`, b, id, userID)
	return scanShoppingList(row)
}

func (db *DB) AppendShoppingListItems(ctx context.Context, id, userID string, newItems []models.ShoppingListItem) (*models.ShoppingList, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var l models.ShoppingList
	var items []byte
	err = tx.QueryRow(ctx, `SELECT id, user_id, title, items, created_at, updated_at FROM shopping_lists WHERE id=$1 AND user_id=$2 FOR UPDATE`, id, userID).Scan(&l.ID, &l.UserID, &l.Title, &items, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(items, &l.Items); err != nil {
		return nil, err
	}
	if l.Items == nil {
		l.Items = []models.ShoppingListItem{}
	}
	l.Items = append(l.Items, newItems...)
	b, _ := json.Marshal(l.Items)
	row := tx.QueryRow(ctx, `UPDATE shopping_lists SET items=$1 WHERE id=$2 AND user_id=$3 RETURNING id, user_id, title, items, created_at, updated_at`, b, id, userID)
	res, err := scanShoppingList(row)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return res, nil
}

func (db *DB) ToggleShoppingListItem(ctx context.Context, id, userID string, idx int, checked bool) (*models.ShoppingList, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var l models.ShoppingList
	var items []byte
	err = tx.QueryRow(ctx, `SELECT id, user_id, title, items, created_at, updated_at FROM shopping_lists WHERE id=$1 AND user_id=$2 FOR UPDATE`, id, userID).Scan(&l.ID, &l.UserID, &l.Title, &items, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(items, &l.Items); err != nil {
		return nil, err
	}
	if l.Items == nil {
		l.Items = []models.ShoppingListItem{}
	}
	if idx < 0 || idx >= len(l.Items) {
		return nil, &IndexOutOfRangeError{Idx: idx, Len: len(l.Items)}
	}
	l.Items[idx].Checked = checked
	b, _ := json.Marshal(l.Items)
	row := tx.QueryRow(ctx, `UPDATE shopping_lists SET items=$1 WHERE id=$2 AND user_id=$3 RETURNING id, user_id, title, items, created_at, updated_at`, b, id, userID)
	res, err := scanShoppingList(row)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return res, nil
}

type IndexOutOfRangeError struct {
	Idx int
	Len int
}

func (e *IndexOutOfRangeError) Error() string { return "itemIndex out of range" }

func (db *DB) DeleteShoppingList(ctx context.Context, id, userID string) error {
	res, err := db.Pool.Exec(ctx, `DELETE FROM shopping_lists WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
