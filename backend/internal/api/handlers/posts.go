package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"cookmode/internal/api/middleware"
)

// RegisterPosts mounts "I made this" photo posts (photos stored in R2).
func RegisterPosts(rg *gin.RouterGroup, pool *pgxpool.Pool, auth *middleware.Auth) {
	h := &postHandler{db: pool}

	posts := rg.Group("/posts", auth.Middleware())
	{
		posts.GET("", h.list)
		posts.POST("", h.create)
		posts.DELETE("/:id", h.delete)
	}
}

type postHandler struct {
	db *pgxpool.Pool
}

func (h *postHandler) list(c *gin.Context)   { c.Status(501) } // TODO: list by recipe_id or user_id
func (h *postHandler) create(c *gin.Context) { c.Status(501) } // TODO: accept R2 photo key + recipe_id
func (h *postHandler) delete(c *gin.Context) { c.Status(501) }
