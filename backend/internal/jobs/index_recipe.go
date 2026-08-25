package jobs

import (
	"context"

	"github.com/riverqueue/river"
)

// IndexRecipeArgs upserts a published recipe into Meilisearch so it is
// searchable by ingredients, time, and dietary needs.
type IndexRecipeArgs struct {
	RecipeID string `json:"recipe_id"`
}

func (IndexRecipeArgs) Kind() string { return "index_recipe" }

type IndexRecipeWorker struct {
	river.WorkerDefaults[IndexRecipeArgs]
}

func (w *IndexRecipeWorker) Work(ctx context.Context, job *river.Job[IndexRecipeArgs]) error {
	// TODO: load recipe, upsert into Meilisearch recipes index,
	// flip status processing -> published.
	return nil
}
