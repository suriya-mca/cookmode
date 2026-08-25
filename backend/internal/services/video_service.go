package services

import (
	"net/http"
)

// VideoService talks to Cloudflare Stream: create direct-upload URLs,
// poll processing status, fetch HLS playback URLs.
type VideoService struct {
	accountID string
	apiToken  string
	client    *http.Client
}

func NewVideoService(accountID, apiToken string) *VideoService {
	return &VideoService{
		accountID: accountID,
		apiToken:  apiToken,
		client:    &http.Client{},
	}
}

// TODO: CreateDirectUploadURL() (uid, uploadURL, error)
// TODO: GetVideo(uid) -> status, hls_url, thumbnail_url, duration
// TODO: OnUploadComplete -> enqueue transcription job
