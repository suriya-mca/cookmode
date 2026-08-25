package meilisearch

import (
	sdk "github.com/meilisearch/meilisearch-go"
)

// Client is a thin wrapper around the Meilisearch SDK.
type Client struct {
	inner sdk.ServiceManager
}

func New(host, apiKey string) *Client {
	return &Client{inner: sdk.New(host, sdk.WithAPIKey(apiKey))}
}

func (c *Client) Inner() sdk.ServiceManager {
	return c.inner
}
