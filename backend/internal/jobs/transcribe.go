package jobs

import (
	"context"

	"github.com/riverqueue/river"
)

// TranscribeArgs runs speech-to-text on a published video and stores
// timestamped segments for anchor-point inference.
type TranscribeArgs struct {
	RecipeID string `json:"recipe_id"`
	VideoUID string `json:"video_uid"`
}

func (TranscribeArgs) Kind() string { return "transcribe" }

type TranscribeWorker struct {
	river.WorkerDefaults[TranscribeArgs]
}

func (w *TranscribeWorker) Work(ctx context.Context, job *river.Job[TranscribeArgs]) error {
	// TODO: call TranscriptionService, store segments, enqueue InferAnchorsArgs.
	return nil
}
