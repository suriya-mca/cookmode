package handlers

import "cookmode/internal/models"

// isRecipeVisible determines whether a recipe is visible to a user.
// isRecipeVisible reports whether a recipe is visible to a user. It returns true when the recipe is published or owned by the specified user, and false otherwise.
func isRecipeVisible(recipe *models.Recipe, userID string) bool {
	return recipe.Status == models.StatusPublished || recipe.UserID == userID
}
