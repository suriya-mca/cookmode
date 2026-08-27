package httpx

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimit = 20
	maxLimit     = 50
)

type Pagination struct {
	Limit  int
	Cursor string
}

func ParsePagination(c *gin.Context) Pagination {
	p := Pagination{Limit: defaultLimit, Cursor: c.Query("cursor")}
	if raw := c.Query("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			if n > maxLimit {
				n = maxLimit
			}
			p.Limit = n
		}
	}
	return p
}
