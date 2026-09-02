package jobs

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"

	"cookmode/internal/services"
)

// Workers returns all River workers for registration. The pipeline order is:
//
// Workers creates a River worker registry containing the video processing, transcription, anchor inference, nutrition fetching, and recipe indexing workers.
func Workers(pool *pgxpool.Pool, indexer *services.SearchIndexer) *river.Workers {
	workers := river.NewWorkers()
	river.AddWorker(workers, &ProcessVideoWorker{Pool: pool})
	river.AddWorker(workers, &TranscribeWorker{Pool: pool})
	river.AddWorker(workers, &InferAnchorsWorker{Pool: pool})
	river.AddWorker(workers, &FetchNutritionWorker{})
	river.AddWorker(workers, &IndexRecipeWorker{Pool: pool, Indexer: indexer})
	return workers
}

// NewClient creates the River client on top of the shared pgx pool.
func NewClient(pool *pgxpool.Pool, indexer *services.SearchIndexer) (*river.Client[pgx.Tx], error) {
	return river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			"default": {MaxWorkers: 10},
		},
		Workers: Workers(pool, indexer),
	})
}

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	migrator, err := rivermigrate.New(riverpgxv5.New(pool), nil)
	if err != nil {
		return err
	}
	_, err = migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	return err
}
