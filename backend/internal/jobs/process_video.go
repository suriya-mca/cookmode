package jobs

import (
	"context"

	"github.com/riverqueue/river"
)

// ProcessVideoArgs fires when Cloudflare Stream finishes processing an
// upload; it moves the recipe from draft to processing and kicks off
// transcription.
type ProcessVideoArgs struct {
	RecipeID string `json:"recipe_id"`
	VideoUID string `json:"video_uid"`
}

func (ProcessVideoArgs) Kind() string { return "process_video" }

type ProcessVideoWorker struct {
	river.WorkerDefaults[ProcessVideoArgs]
}

func (w *ProcessVideoWorker) Work(ctx context.Context, job *river.Job[ProcessVideoArgs]) error {
	// TODO: fetch HLS/thumbnail/duration from Stream, update recipe row,
	// enqueue TranscribeArgs.
	return nil
}
