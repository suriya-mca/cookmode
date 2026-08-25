package services

// NutritionService computes per-serving nutrition from an ingredient list.
// Primary source: Edamam; fallback: USDA FoodData Central. No computer vision.
type NutritionService struct {
	edamamAppID  string
	edamamAppKey string
	usdaAPIKey   string
}

func NewNutritionService(edamamAppID, edamamAppKey, usdaAPIKey string) *NutritionService {
	return &NutritionService{
		edamamAppID:  edamamAppID,
		edamamAppKey: edamamAppKey,
		usdaAPIKey:   usdaAPIKey,
	}
}

// TODO: Compute(ingredients []models.Ingredient, servings int) (*models.Nutrition, error)
// TODO: Scale(nutrition *models.Nutrition, fromServings, toServings int)
