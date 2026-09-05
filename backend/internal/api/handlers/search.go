package handlers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"cookmode/internal/httpx"
	"cookmode/internal/services"
)

// RegisterSearch registers the GET /search endpoint on rg using indexer.
// The endpoint supports search filters and pagination, defaults to a limit of 20
// and an offset of 0, returns a bad-request response for invalid time values,
// RegisterSearch registers a GET /search endpoint that returns filtered, paginated search results as JSON.
// It supports query, ingredient, dietary, difficulty, cuisine, and maximum-time filters, with default
// pagination of 20 results and offset 0. Invalid maximum-time values produce a 400 response, and search
// failures produce a 500 response.
func RegisterSearch(rg *gin.RouterGroup, indexer *services.SearchIndexer) {
	rg.GET("/search", func(c *gin.Context) {
		q := c.Query("q")
		ingredients := splitCSV(c.Query("ingredients"))
		dietary := splitCSV(c.Query("dietary"))
		if len(dietary) == 0 {
			dietary = splitCSV(c.Query("dietary_tags"))
		}
		difficulty := c.Query("difficulty")
		cuisine := c.Query("cuisine")
		var maxTime *int
		if v := c.Query("max_total_time"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil || n < 0 {
				httpx.ErrBadRequest(c, "max_total_time must be a non-negative integer")
				return
			}
			maxTime = &n
		} else if v := c.Query("max_time"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil || n < 0 {
				httpx.ErrBadRequest(c, "max_time must be a non-negative integer")
				return
			}
			maxTime = &n
		}
		limit := 20
		if v := c.Query("limit"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 50 {
				limit = n
			}
		}
		offset := 0
		if v := c.Query("offset"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n >= 0 {
				offset = n
			}
		}
		res, err := indexer.Search(c.Request.Context(), services.SearchParams{
			Query:       q,
			Ingredients: ingredients,
			Dietary:     dietary,
			Difficulty:  difficulty,
			Cuisine:     cuisine,
			MaxTime:     maxTime,
			Limit:       limit,
			Offset:      offset,
		})
		if err != nil {
			httpx.ErrInternal(c, "search failed")
			return
		}
		c.JSON(200, gin.H{
			"hits":               res.Hits,
			"query":              res.Query,
			"processingTimeMs":   res.ProcessingTimeMs,
			"limit":              res.Limit,
			"offset":             res.Offset,
			"estimatedTotalHits": res.EstimatedTotalHits,
		})
	})
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
