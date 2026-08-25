package cloudflarestream

import (
	"net/http"
)

// Client wraps the Cloudflare Stream REST API: direct uploads,
// processing status, HLS playback URLs.
type Client struct {
	accountID string
	apiToken  string
	baseURL   string
	http      *http.Client
}

func New(accountID, apiToken string) *Client {
	return &Client{
		accountID: accountID,
		apiToken:  apiToken,
		baseURL:   "https://api.cloudflare.com/client/v4/accounts/" + accountID + "/stream",
		http:      &http.Client{},
	}
}
