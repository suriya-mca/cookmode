package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type VideoService struct {
	accountID string
	apiToken  string
	client    *http.Client
}

func NewVideoService(accountID, apiToken string) *VideoService {
	return &VideoService{
		accountID: accountID,
		apiToken:  apiToken,
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

type directUploadResponse struct {
	Result struct {
		UID       string `json:"uid"`
		UploadURL string `json:"uploadURL"`
	} `json:"result"`
	Success bool `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (s *VideoService) CreateDirectUploadURL(ctx context.Context) (string, string, error) {
	if s.accountID == "" || s.apiToken == "" || s.accountID == "your-cloudflare-account-id" {
		uid := uuid.NewString()
		return uid, fmt.Sprintf("https://dev.local/upload/%s", uid), nil
	}

	body, _ := json.Marshal(map[string]any{
		"maxDurationSeconds": 3600,
		"expiry":             "2030-01-01T00:00:00Z",
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/stream/direct_upload", s.accountID),
		bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("stream api: %w", err)
	}
	defer resp.Body.Close()

	var out directUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", "", fmt.Errorf("stream decode: %w", err)
	}
	if !out.Success {
		msg := "stream direct_upload failed"
		if len(out.Errors) > 0 {
			msg = out.Errors[0].Message
		}
		return "", "", fmt.Errorf("%s", msg)
	}
	return out.Result.UID, out.Result.UploadURL, nil
}
