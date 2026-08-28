package handlers

import "cookmode/internal/models"

func isRecipeVisible(recipe *models.Recipe, userID string) bool {
	return recipe.Status == models.StatusPublished || recipe.UserID == userID
}
