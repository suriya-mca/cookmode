package models

import "time"

// User mirrors the Supabase auth.users identity; profile data lives here.
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

type PostStatus string

const (
	PostVisible  PostStatus = "visible"
	PostArchived PostStatus = "archived"
)

// Post is an "I made this" photo attached to a recipe.
type Post struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	RecipeID    string     `json:"recipe_id"`
	PhotoURL    string     `json:"photo_url"` // served from R2
	Caption     string     `json:"caption"`
	Status      PostStatus `json:"status"`
	RatingValue int        `json:"rating_value,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type ShoppingListItem struct {
	IngredientName string  `json:"ingredient_name"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	Checked        bool    `json:"checked"`
}

type ShoppingList struct {
	ID        string             `json:"id"`
	UserID    string             `json:"user_id"`
	Title     string             `json:"title"`
	Items     []ShoppingListItem `json:"items"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
