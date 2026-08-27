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

func RegisterCollections(rg *gin.RouterGroup, pool *pgxpool.Pool, auth *middleware.Auth) {
	h := &collectionHandler{db: queries.New(pool)}
	g := rg.Group("/collections", auth.Middleware())
	{
		g.GET("", h.list)
		g.POST("", h.create)
		g.GET("/:id", h.get)
		g.PATCH("/:id", h.update)
		g.DELETE("/:id", h.delete)
		g.GET("/:id/recipes", h.listRecipes)
		g.POST("/:id/recipes", h.addRecipe)
		g.DELETE("/:id/recipes/:recipeId", h.removeRecipe)
	}
}

type collectionHandler struct {
	db *queries.DB
}

func (h *collectionHandler) list(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	cols, err := h.db.ListCollections(c.Request.Context(), userID)
	if err != nil {
		httpx.ErrInternal(c, "failed to list collections")
		return
	}
	if cols == nil {
		cols = []*models.Collection{}
	}
	c.JSON(http.StatusOK, gin.H{"data": cols})
}

func (h *collectionHandler) create(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		IsPublic    *bool  `json:"is_public"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Title == "" {
		httpx.ErrBadRequest(c, "title is required")
		return
	}
	if len(req.Title) > 100 {
		httpx.ErrBadRequest(c, "title must be at most 100 characters")
		return
	}
	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}
	col, err := h.db.CreateCollection(c.Request.Context(), &models.Collection{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		IsPublic:    isPublic,
	})
	if err != nil {
		httpx.ErrInternal(c, "failed to create collection")
		return
	}
	c.JSON(http.StatusCreated, col)
}

func (h *collectionHandler) get(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	col, err := h.db.GetCollection(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "collection not found")
			return
		}
		httpx.ErrInternal(c, "failed to get collection")
		return
	}
	if !col.IsPublic && col.UserID != userID {
		httpx.ErrNotFound(c, "collection not found")
		return
	}
	c.JSON(http.StatusOK, col)
}

func (h *collectionHandler) update(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	col, err := h.db.GetCollection(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "collection not found")
			return
		}
		httpx.ErrInternal(c, "failed to get collection")
		return
	}
	if col.UserID != userID {
		httpx.ErrForbidden(c, "not your collection")
		return
	}
	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		IsPublic    *bool   `json:"is_public"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrBadRequest(c, "invalid body")
		return
	}
	if req.Title != nil {
		if *req.Title == "" || len(*req.Title) > 100 {
			httpx.ErrBadRequest(c, "title must be 1-100 characters")
			return
		}
		col.Title = *req.Title
	}
	if req.Description != nil {
		col.Description = *req.Description
	}
	if req.IsPublic != nil {
		col.IsPublic = *req.IsPublic
	}
	updated, err := h.db.UpdateCollection(c.Request.Context(), col)
	if err != nil {
		httpx.ErrInternal(c, "failed to update collection")
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *collectionHandler) delete(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	col, err := h.db.GetCollection(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "collection not found")
			return
		}
		httpx.ErrInternal(c, "failed to get collection")
		return
	}
	if col.UserID != userID {
		httpx.ErrForbidden(c, "not your collection")
		return
	}
	if err := h.db.DeleteCollection(c.Request.Context(), id); err != nil {
		httpx.ErrInternal(c, "failed to delete collection")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *collectionHandler) listRecipes(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	col, err := h.db.GetCollection(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "collection not found")
			return
		}
		httpx.ErrInternal(c, "failed to get collection")
		return
	}
	if !col.IsPublic && col.UserID != userID {
		httpx.ErrNotFound(c, "collection not found")
		return
	}
	recipes, err := h.db.ListCollectionRecipes(c.Request.Context(), id)
	if err != nil {
		httpx.ErrInternal(c, "failed to list recipes")
		return
	}
	if recipes == nil {
		recipes = []*models.Recipe{}
	}
	c.JSON(http.StatusOK, gin.H{"data": recipes})
}

func (h *collectionHandler) addRecipe(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	col, err := h.db.GetCollection(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "collection not found")
			return
		}
		httpx.ErrInternal(c, "failed to get collection")
		return
	}
	if col.UserID != userID {
		httpx.ErrForbidden(c, "not your collection")
		return
	}
	var req struct {
		RecipeID string `json:"recipe_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.RecipeID == "" {
		httpx.ErrBadRequest(c, "recipe_id is required")
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
	if err := h.db.AddRecipeToCollection(c.Request.Context(), id, req.RecipeID); err != nil {
		httpx.ErrInternal(c, "failed to add recipe")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"added": true})
}

func (h *collectionHandler) removeRecipe(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	recipeID := c.Param("recipeId")
	col, err := h.db.GetCollection(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "collection not found")
			return
		}
		httpx.ErrInternal(c, "failed to get collection")
		return
	}
	if col.UserID != userID {
		httpx.ErrForbidden(c, "not your collection")
		return
	}
	if err := h.db.RemoveRecipeFromCollection(c.Request.Context(), id, recipeID); err != nil {
		httpx.ErrInternal(c, "failed to remove recipe")
		return
	}
	c.Status(http.StatusNoContent)
}
