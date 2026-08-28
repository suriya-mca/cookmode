package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"cookmode/internal/api/middleware"
	"cookmode/internal/db/queries"
	"cookmode/internal/httpx"
	"cookmode/internal/models"
)

func RegisterPosts(rg *gin.RouterGroup, pool *pgxpool.Pool, auth *middleware.Auth) {
	h := &postHandler{db: queries.New(pool)}
	posts := rg.Group("/posts", auth.Middleware())
	{
		posts.GET("", h.list)
		posts.POST("", h.create)
		posts.DELETE("/:id", h.delete)
	}
}

type postHandler struct {
	db *queries.DB
}

func (h *postHandler) list(c *gin.Context) {
	recipeID := c.Query("recipe_id")
	userID := c.Query("user_id")
	if recipeID == "" && userID == "" {
		httpx.ErrBadRequest(c, "recipe_id or user_id is required")
		return
	}
	posts, err := h.db.ListPosts(c.Request.Context(), recipeID, userID)
	if err != nil {
		httpx.ErrInternal(c, "failed to list posts")
		return
	}
	if posts == nil {
		posts = []*models.Post{}
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func (h *postHandler) create(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	var req struct {
		RecipeID    string `json:"recipe_id"`
		PhotoURL    string `json:"photo_url"`
		Caption     string `json:"caption"`
		RatingValue int    `json:"rating_value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.RecipeID == "" || req.PhotoURL == "" {
		httpx.ErrBadRequest(c, "recipe_id and photo_url are required")
		return
	}
	if req.RatingValue < 0 || req.RatingValue > 5 {
		httpx.ErrBadRequest(c, "rating_value must be 0-5")
		return
	}
	if len(req.Caption) > 500 {
		httpx.ErrBadRequest(c, "caption must be at most 500 characters")
		return
	}
	if _, err := h.db.GetRecipe(c.Request.Context(), req.RecipeID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "recipe not found")
			return
		}
		httpx.ErrInternal(c, "failed to get recipe")
		return
	}
	post, err := h.db.CreatePost(c.Request.Context(), &models.Post{
		UserID:      userID,
		RecipeID:    req.RecipeID,
		PhotoURL:    req.PhotoURL,
		Caption:     req.Caption,
		RatingValue: req.RatingValue,
	})
	if err != nil {
		httpx.ErrInternal(c, "failed to create post")
		return
	}
	c.JSON(http.StatusCreated, post)
}

func (h *postHandler) delete(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	if err := h.db.DeletePost(c.Request.Context(), id, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "post not found")
			return
		}
		httpx.ErrInternal(c, "failed to delete post")
		return
	}
	c.Status(http.StatusNoContent)
}
