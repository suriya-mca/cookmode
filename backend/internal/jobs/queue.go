package jobs

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

// Workers returns all River workers for registration. The pipeline order is:
//
//	process_video -> transcribe -> infer_anchors -> fetch_nutrition + index_recipe
func Workers() *river.Workers {
	workers := river.NewWorkers()
	river.AddWorker(workers, &ProcessVideoWorker{})
	river.AddWorker(workers, &TranscribeWorker{})
	river.AddWorker(workers, &InferAnchorsWorker{})
	river.AddWorker(workers, &FetchNutritionWorker{})
	river.AddWorker(workers, &IndexRecipeWorker{})
	return workers
}

// NewClient creates the River client on top of the shared pgx pool.
func NewClient(pool *pgxpool.Pool) (*river.Client[pgx.Tx], error) {
	return river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			"default": {MaxWorkers: 10},
		},
		Workers: Workers(),
	})
}
