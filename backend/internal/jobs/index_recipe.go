package jobs

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"

	"cookmode/internal/db/queries"
	"cookmode/internal/services"
)

// IndexRecipeArgs upserts a published recipe into Meilisearch so it is
// searchable by ingredients, time, and dietary needs.
type IndexRecipeArgs struct {
	RecipeID string `json:"recipe_id"`
}

func (IndexRecipeArgs) Kind() string { return "index_recipe" }

type IndexRecipeWorker struct {
	river.WorkerDefaults[IndexRecipeArgs]
	Pool    *pgxpool.Pool
	Indexer *services.SearchIndexer
}

func (w *IndexRecipeWorker) Work(ctx context.Context, job *river.Job[IndexRecipeArgs]) error {
	q := queries.New(w.Pool)
	recipe, err := q.GetRecipe(ctx, job.Args.RecipeID)
	if err != nil {
		return err
	}
	if err := w.Indexer.IndexRecipe(ctx, recipe); err != nil {
		return err
	}
	if recipe.Status == "processing" {
		_, _ = q.Pool.Exec(ctx, `UPDATE recipes SET status='published' WHERE id=$1`, recipe.ID)
	}
	return nil
}
