package jobs

import (
	"context"

	"github.com/riverqueue/river"
)

// InferAnchorsArgs maps each recipe step to a video anchor point
// (anchor_seconds) by aligning step text with transcript segments.
type InferAnchorsArgs struct {
	RecipeID string `json:"recipe_id"`
}

func (InferAnchorsArgs) Kind() string { return "infer_anchors" }

type InferAnchorsWorker struct {
	river.WorkerDefaults[InferAnchorsArgs]
}

func (w *InferAnchorsWorker) Work(ctx context.Context, job *river.Job[InferAnchorsArgs]) error {
	// TODO: align steps with transcript, write anchor_seconds into steps JSONB,
	// enqueue FetchNutritionArgs and IndexRecipeArgs.
	return nil
}
