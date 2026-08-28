package httpx

import (
	"net/http/httptest"
	"testing"

	"cookmode/internal/models"

	"github.com/gin-gonic/gin"
)

func TestValidateRecipe(t *testing.T) {
	tests := []struct {
		name   string
		recipe models.Recipe
		want   string
	}{
		{"valid minimal", models.Recipe{Title: "Soup", Servings: 2}, ""},
		{"empty title", models.Recipe{Title: ""}, "title is required"},
		{"title too long", models.Recipe{Title: string(make([]byte, 201))}, "title must be at most 200 characters"},
		{"description too long", models.Recipe{Title: "T", Description: string(make([]byte, 5001))}, "description must be at most 5000 characters"},
		{"negative servings", models.Recipe{Title: "T", Servings: -1}, "servings must be >= 0"},
		{"negative prep", models.Recipe{Title: "T", PrepTimeMin: -1}, "prep_time and cook_time must be >= 0"},
		{"invalid difficulty", models.Recipe{Title: "T", Difficulty: "extreme"}, "difficulty must be easy, medium, or hard"},
		{"empty dietary tag", models.Recipe{Title: "T", DietaryTags: []string{"vegan", ""}}, "dietary_tags must not contain empty values"},
		{"empty ingredient name", models.Recipe{Title: "T", Ingredients: []models.Ingredient{{Name: ""}}}, "ingredients[0].name is required"},
		{"negative quantity", models.Recipe{Title: "T", Ingredients: []models.Ingredient{{Name: "salt", Quantity: -1}}}, "ingredients[0].quantity must be >= 0"},
		{"empty step text", models.Recipe{Title: "T", Steps: []models.Step{{Text: ""}}}, "steps[0].text is required"},
		{"negative anchor", models.Recipe{Title: "T", Steps: []models.Step{{Text: "do", AnchorSeconds: -1}}}, "steps[0].anchor_seconds must be >= 0"},
		{"valid difficulty", models.Recipe{Title: "T", Difficulty: models.DifficultyEasy}, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ValidateRecipe(&tc.recipe); got != tc.want {
				t.Errorf("ValidateRecipe() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name  string
		query string
		wantL int
		wantC string
	}{
		{"defaults", "", 20, ""},
		{"limit 10", "limit=10", 10, ""},
		{"limit over max", "limit=100", 50, ""},
		{"limit invalid", "limit=abc", 20, ""},
		{"cursor", "cursor=abc123", 20, "abc123"},
		{"both", "limit=5&cursor=xyz", 5, "xyz"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/?"+tc.query, nil)
			got := ParsePagination(c)
			if got.Limit != tc.wantL || got.Cursor != tc.wantC {
				t.Errorf("ParsePagination() = %+v, want Limit=%d Cursor=%q", got, tc.wantL, tc.wantC)
			}
		})
	}
}
