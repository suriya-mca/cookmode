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
	// Flip first so the subsequent index sees published; IndexRecipe
	// deletes non-published docs, so indexing before the flip would
	// delete then never re-index.
	_, err := q.Pool.Exec(ctx, `UPDATE recipes SET status='published' WHERE id=$1 AND status='processing'`, job.Args.RecipeID)
	if err != nil {
		return err
	}
	recipe, err := q.GetRecipe(ctx, job.Args.RecipeID)
	if err != nil {
		return err
	}
	return w.Indexer.IndexRecipe(ctx, recipe)
}
