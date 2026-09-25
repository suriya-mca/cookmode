package httpx

import (
	"regexp"
	"strings"

	"cookmode/internal/models"
)

// usernamePattern is the single username policy for signup, login, and
// profile updates: lowercase ASCII letters, digits, underscores.
var usernamePattern = regexp.MustCompile(`^[a-z0-9_]+$`)

// NormalizeUsername trims and lowercases a raw username and enforces the
// policy. It returns the normalized value, or an error message for the API.
func NormalizeUsername(raw string) (string, string) {
	u := strings.ToLower(strings.TrimSpace(raw))
	if len(u) < 3 || len(u) > 30 || !usernamePattern.MatchString(u) {
		return "", "username must be 3-30 lowercase letters, digits, or underscores"
	}
	return u, ""
}

// ValidatePassword enforces the password policy. The 72-char cap is bcrypt's
// input limit — reject it with a 400 instead of a 500 from the hasher.
func ValidatePassword(pw string) string {
	if len(pw) < 8 || len(pw) > 72 {
		return "password must be 8-72 characters"
	}
	return ""
}

func ValidateRecipe(r *models.Recipe) string {
	if strings.TrimSpace(r.Title) == "" {
		return "title is required"
	}
	if len(r.Title) > 200 {
		return "title must be at most 200 characters"
	}
	if len(r.Description) > 5000 {
		return "description must be at most 5000 characters"
	}
	if r.Servings < 0 {
		return "servings must be >= 0"
	}
	if r.PrepTimeMin < 0 || r.CookTimeMin < 0 {
		return "prep_time and cook_time must be >= 0"
	}
	if r.Difficulty != "" && r.Difficulty != models.DifficultyEasy && r.Difficulty != models.DifficultyMedium && r.Difficulty != models.DifficultyHard {
		return "difficulty must be easy, medium, or hard"
	}
	for _, tag := range r.DietaryTags {
		if strings.TrimSpace(tag) == "" {
			return "dietary_tags must not contain empty values"
		}
	}
	for i, ing := range r.Ingredients {
		if strings.TrimSpace(ing.Name) == "" {
			return "ingredients[" + itoa(i) + "].name is required"
		}
		if ing.Quantity < 0 {
			return "ingredients[" + itoa(i) + "].quantity must be >= 0"
		}
	}
	for i, s := range r.Steps {
		if strings.TrimSpace(s.Text) == "" {
			return "steps[" + itoa(i) + "].text is required"
		}
		if s.AnchorSeconds < 0 {
			return "steps[" + itoa(i) + "].anchor_seconds must be >= 0"
		}
	}
	return ""
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
