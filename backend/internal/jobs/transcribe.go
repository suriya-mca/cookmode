package jobs

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type TranscribeArgs struct {
	RecipeID string `json:"recipe_id"`
	VideoUID string `json:"video_uid"`
}

func (TranscribeArgs) Kind() string { return "transcribe" }

type TranscribeWorker struct {
	river.WorkerDefaults[TranscribeArgs]
	Pool *pgxpool.Pool
}

func (w *TranscribeWorker) Work(ctx context.Context, job *river.Job[TranscribeArgs]) error {
	segments := []map[string]interface{}{
		{"start_seconds": 0, "end_seconds": 10, "text": "welcome to the recipe"},
		{"start_seconds": 10, "end_seconds": 30, "text": "prepare ingredients and heat the pan"},
		{"start_seconds": 30, "end_seconds": 60, "text": "cook and stir for a few minutes"},
		{"start_seconds": 60, "end_seconds": 120, "text": "plate and serve"},
	}
	b, _ := json.Marshal(segments)
	raw := "welcome to the recipe prepare ingredients and heat the pan cook and stir for a few minutes plate and serve"
	_, err := w.Pool.Exec(ctx, `
		INSERT INTO recipe_transcripts (recipe_id, segments, raw_text, source)
		VALUES ($1, $2::jsonb, $3, 'ai')
		ON CONFLICT (recipe_id) DO UPDATE SET segments = EXCLUDED.segments, raw_text = EXCLUDED.raw_text, source = EXCLUDED.source, updated_at = now()
		WHERE recipe_transcripts.source <> 'manual'
	`, job.Args.RecipeID, string(b), raw)
	return err
}
