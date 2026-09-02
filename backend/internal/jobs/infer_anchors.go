package jobs

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"

	"cookmode/internal/models"
)

type InferAnchorsArgs struct {
	RecipeID string `json:"recipe_id"`
}

func (InferAnchorsArgs) Kind() string { return "infer_anchors" }

type InferAnchorsWorker struct {
	river.WorkerDefaults[InferAnchorsArgs]
	Pool *pgxpool.Pool
}

func (w *InferAnchorsWorker) Work(ctx context.Context, job *river.Job[InferAnchorsArgs]) error {
	var stepsJSON []byte
	var duration float64
	err := w.Pool.QueryRow(ctx, `SELECT steps, video_duration_sec FROM recipes WHERE id=$1`, job.Args.RecipeID).Scan(&stepsJSON, &duration)
	if err != nil {
		return err
	}
	var steps []models.Step
	if err := json.Unmarshal(stepsJSON, &steps); err != nil {
		return err
	}
	if len(steps) == 0 {
		return nil
	}
	if duration == 0 {
		duration = 120
	}
	for i := range steps {
		steps[i].AnchorSeconds = duration * float64(i) / float64(len(steps))
	}
	b, _ := json.Marshal(steps)
	_, err = w.Pool.Exec(ctx, `UPDATE recipes SET steps=$1 WHERE id=$2`, string(b), job.Args.RecipeID)
	return err
}
