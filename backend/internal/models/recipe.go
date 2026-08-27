package models

import "time"

type RecipeStatus string

const (
	StatusDraft      RecipeStatus = "draft"
	StatusProcessing RecipeStatus = "processing"
	StatusPublished  RecipeStatus = "published"
	StatusArchived   RecipeStatus = "archived"
)

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type Ingredient struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Optional bool    `json:"optional,omitempty"`
}

type Step struct {
	Order         int     `json:"order"`
	Text          string  `json:"text"`
	AnchorSeconds float64 `json:"anchor_seconds"`
	DurationHint  int     `json:"duration_hint,omitempty"`
}

type Substitution struct {
	ForIngredient string   `json:"for_ingredient"`
	Alternatives  []string `json:"alternatives"`
}

type Nutrition struct {
	Calories      float64 `json:"calories"`
	ProteinGrams  float64 `json:"protein_grams"`
	CarbGrams     float64 `json:"carb_grams"`
	FatGrams      float64 `json:"fat_grams"`
	FiberGrams    float64 `json:"fiber_grams,omitempty"`
	ServingWeight int     `json:"serving_weight,omitempty"`
	Source        string  `json:"source"`
}

type Recipe struct {
	ID                string         `json:"id"`
	UserID            string         `json:"user_id"`
	Title             string         `json:"title"`
	Description       string         `json:"description"`
	Cuisine           string         `json:"cuisine"`
	PrepTimeMin       int            `json:"prep_time_min"`
	CookTimeMin       int            `json:"cook_time_min"`
	Servings          int            `json:"servings"`
	Difficulty        Difficulty     `json:"difficulty"`
	DietaryTags       []string       `json:"dietary_tags"`
	Ingredients       []Ingredient   `json:"ingredients"`
	Steps             []Step         `json:"steps"`
	Substitutions     []Substitution `json:"substitutions"`
	Nutrition         *Nutrition     `json:"nutrition"`
	VideoUID          string         `json:"video_uid"`
	VideoHLSURL       string         `json:"video_hls_url"`
	VideoThumbnailURL string         `json:"video_thumbnail_url"`
	VideoDurationSec  float64        `json:"video_duration_sec"`
	Views             int            `json:"views"`
	Saves             int            `json:"saves"`
	Status            RecipeStatus   `json:"status"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

type TranscriptSegment struct {
	StartSeconds float64 `json:"start_seconds"`
	EndSeconds   float64 `json:"end_seconds"`
	Text         string  `json:"text"`
}

type RecipeTranscript struct {
	RecipeID  string              `json:"recipe_id"`
	Segments  []TranscriptSegment `json:"segments"`
	RawText   string              `json:"raw_text"`
	Source    string              `json:"source"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}
