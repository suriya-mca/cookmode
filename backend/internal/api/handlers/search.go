package handlers

import (
	"github.com/gin-gonic/gin"

	"cookmode/internal/services"
)

// RegisterSearch mounts ingredient/time/dietary search backed by Meilisearch.
func RegisterSearch(rg *gin.RouterGroup, indexer *services.SearchIndexer) {
	rg.GET("/search", func(c *gin.Context) {
		c.Status(501) // TODO: query Meilisearch recipes index (ingredients, time, dietary_tags)
	})
}
