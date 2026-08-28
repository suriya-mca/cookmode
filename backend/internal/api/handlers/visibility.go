package handlers

import "cookmode/internal/models"

// isRecipeVisible determines whether a recipe is visible to a user.
// It returns true for published recipes and recipes owned by the specified user.
func isRecipeVisible(recipe *models.Recipe, userID string) bool {
	return recipe.Status == models.StatusPublished || recipe.UserID == userID
}
