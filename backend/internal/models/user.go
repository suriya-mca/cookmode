package models

import "time"

type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	AvatarURL   string    `json:"avatar_url"`
	Bio         string    `json:"bio"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PostStatus string

const (
	PostVisible  PostStatus = "visible"
	PostArchived PostStatus = "archived"
)

type Post struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	RecipeID    string     `json:"recipe_id"`
	PhotoURL    string     `json:"photo_url"`
	Caption     string     `json:"caption"`
	Status      PostStatus `json:"status"`
	RatingValue *int       `json:"rating_value,omitempty"`
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

type Follow struct {
	FollowerID  string    `json:"follower_id"`
	FollowingID string    `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Save struct {
	UserID    string    `json:"user_id"`
	RecipeID  string    `json:"recipe_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Collection struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CollectionRecipe struct {
	CollectionID string    `json:"collection_id"`
	RecipeID     string    `json:"recipe_id"`
	AddedAt      time.Time `json:"added_at"`
}

type RecipeFork struct {
	ID             string    `json:"id"`
	ParentRecipeID string    `json:"parent_recipe_id"`
	ChildRecipeID  string    `json:"child_recipe_id"`
	UserID         string    `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
}
