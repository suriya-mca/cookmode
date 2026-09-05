package jobs

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type ProcessVideoArgs struct {
	RecipeID string `json:"recipe_id"`
	VideoUID string `json:"video_uid"`
}

func (ProcessVideoArgs) Kind() string { return "process_video" }

type ProcessVideoWorker struct {
	river.WorkerDefaults[ProcessVideoArgs]
	Pool *pgxpool.Pool
}

func (w *ProcessVideoWorker) Work(ctx context.Context, job *river.Job[ProcessVideoArgs]) error {
	_, err := w.Pool.Exec(ctx, `
		UPDATE recipes SET
			video_hls_url = $1,
			video_thumbnail_url = $2,
			video_duration_sec = $3
		WHERE id = $4
	`,
		fmt.Sprintf("https://dev.local/hls/%s/manifest.m3u8", job.Args.VideoUID),
		fmt.Sprintf("https://dev.local/thumb/%s.jpg", job.Args.VideoUID),
		120.0,
		job.Args.RecipeID,
	)
	if err != nil {
		return err
	}
	if client := river.ClientFromContext[pgx.Tx](ctx); client != nil {
		_, err = client.Insert(ctx, TranscribeArgs{RecipeID: job.Args.RecipeID, VideoUID: job.Args.VideoUID}, nil)
		return err
	}
	return nil
}
