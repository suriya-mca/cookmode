package services

import (
	sdk "github.com/meilisearch/meilisearch-go"
)

// SearchIndexer wraps the Meilisearch client. Recipes are indexed on publish
// (async via River) with filterable attributes for ingredients, time, and
// dietary tags.
type SearchIndexer struct {
	client sdk.ServiceManager
}

func NewSearchIndexer(host, apiKey string) *SearchIndexer {
	return &SearchIndexer{client: sdk.New(host, sdk.WithAPIKey(apiKey))}
}

func (s *SearchIndexer) Client() sdk.ServiceManager {
	return s.client
}

// TODO: IndexRecipe(recipe), DeleteRecipe(id), ConfigureIndexes()
