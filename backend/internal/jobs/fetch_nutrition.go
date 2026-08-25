package jobs

import (
	"context"

	"github.com/riverqueue/river"
)

// FetchNutritionArgs computes per-serving nutrition from the ingredient
// list (Edamam primary, USDA fallback) and writes it to the nutrition JSONB.
type FetchNutritionArgs struct {
	RecipeID string `json:"recipe_id"`
}

func (FetchNutritionArgs) Kind() string { return "fetch_nutrition" }

type FetchNutritionWorker struct {
	river.WorkerDefaults[FetchNutritionArgs]
}

func (w *FetchNutritionWorker) Work(ctx context.Context, job *river.Job[FetchNutritionArgs]) error {
	// TODO: call NutritionService.Compute, persist nutrition JSONB.
	return nil
}
