package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"cookmode/internal/api/middleware"
	"cookmode/internal/db/queries"
	"cookmode/internal/httpx"
	"cookmode/internal/models"
	"cookmode/internal/services"
)

func RegisterRecipes(
	rg *gin.RouterGroup,
	pool *pgxpool.Pool,
	indexer *services.SearchIndexer,
	video *services.VideoService,
	auth *middleware.Auth,
) {
	h := &recipeHandler{db: queries.New(pool), pool: pool, indexer: indexer, video: video}

	recipes := rg.Group("/recipes")
	{
		recipes.GET("", h.list)
		recipes.GET("/:id", auth.Optional(), h.get)

		protected := recipes.Group("", auth.Middleware())
		{
			protected.POST("", h.create)
			protected.POST("/upload-url", h.uploadURL)
			protected.PATCH("/:id", h.update)
			protected.DELETE("/:id", h.archive)
			protected.POST("/:id/publish", h.publish)
			protected.POST("/:id/saves", h.save)
			protected.DELETE("/:id/saves", h.unsave)
			protected.POST("/:id/fork", h.fork)
		}
	}
}

type recipeHandler struct {
	db      *queries.DB
	pool    *pgxpool.Pool
	indexer *services.SearchIndexer
	video   *services.VideoService
}

func (h *recipeHandler) list(c *gin.Context) {
	p := httpx.ParsePagination(c)
	recipes, nextCursor, err := h.db.ListRecipes(c.Request.Context(), p.Limit, p.Cursor)
	if err != nil {
		httpx.ErrInternal(c, "failed to list recipes")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": recipes, "next_cursor": nextCursor})
}

func (h *recipeHandler) get(c *gin.Context) {
	id := c.Param("id")
	recipe, err := h.db.GetRecipe(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "recipe not found")
			return
		}
		httpx.ErrInternal(c, "failed to get recipe")
		return
	}
	if recipe.Status != models.StatusPublished {
		userID := middleware.UserIDFrom(c)
		if userID == "" || userID != recipe.UserID {
			httpx.ErrNotFound(c, "recipe not found")
			return
		}
	}
	go func() {
		_, _ = h.pool.Exec(context.Background(), `UPDATE recipes SET views = views + 1 WHERE id=$1`, id)
	}()
	recipe.Views++
	c.JSON(http.StatusOK, recipe)
}

func (h *recipeHandler) create(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	var req models.Recipe
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrBadRequest(c, "invalid request body")
		return
	}
	req.UserID = userID
	req.Status = models.StatusDraft
	if msg := httpx.ValidateRecipe(&req); msg != "" {
		httpx.ErrBadRequest(c, msg)
		return
	}
	recipe, err := h.db.CreateRecipe(c.Request.Context(), &req)
	if err != nil {
		httpx.ErrInternal(c, "failed to create recipe")
		return
	}
	c.JSON(http.StatusCreated, recipe)
}

func (h *recipeHandler) update(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	existing, err := h.db.GetRecipe(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "recipe not found")
			return
		}
		httpx.ErrInternal(c, "failed to get recipe")
		return
	}
	if existing.UserID != userID {
		httpx.ErrForbidden(c, "not your recipe")
		return
	}
	if existing.Status == models.StatusArchived {
		httpx.ErrBadRequest(c, "archived recipes cannot be edited")
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		httpx.ErrBadRequest(c, "invalid request body")
		return
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		httpx.ErrBadRequest(c, "invalid request body")
		return
	}
	var req models.Recipe
	if err := json.Unmarshal(body, &req); err != nil {
		httpx.ErrBadRequest(c, "invalid request body")
		return
	}
	if _, ok := raw["title"]; !ok {
		req.Title = existing.Title
	}
	if _, ok := raw["description"]; !ok {
		req.Description = existing.Description
	}
	if _, ok := raw["cuisine"]; !ok {
		req.Cuisine = existing.Cuisine
	}
	if _, ok := raw["prep_time_min"]; !ok {
		req.PrepTimeMin = existing.PrepTimeMin
	}
	if _, ok := raw["cook_time_min"]; !ok {
		req.CookTimeMin = existing.CookTimeMin
	}
	if _, ok := raw["servings"]; !ok {
		req.Servings = existing.Servings
	}
	if _, ok := raw["difficulty"]; !ok {
		req.Difficulty = existing.Difficulty
	}
	if _, ok := raw["dietary_tags"]; !ok {
		req.DietaryTags = existing.DietaryTags
	}
	if _, ok := raw["ingredients"]; !ok {
		req.Ingredients = existing.Ingredients
	}
	if _, ok := raw["steps"]; !ok {
		req.Steps = existing.Steps
	}
	if _, ok := raw["substitutions"]; !ok {
		req.Substitutions = existing.Substitutions
	}
	if _, ok := raw["nutrition"]; !ok {
		req.Nutrition = existing.Nutrition
	}
	req.ID = id
	req.UserID = existing.UserID
	req.Status = existing.Status
	req.VideoUID = existing.VideoUID
	req.VideoHLSURL = existing.VideoHLSURL
	req.VideoThumbnailURL = existing.VideoThumbnailURL
	req.VideoDurationSec = existing.VideoDurationSec
	req.Views = existing.Views
	req.Saves = existing.Saves
	req.CreatedAt = existing.CreatedAt
	if msg := httpx.ValidateRecipe(&req); msg != "" {
		httpx.ErrBadRequest(c, msg)
		return
	}
	updated, err := h.db.UpdateRecipe(c.Request.Context(), &req)
	if err != nil {
		httpx.ErrInternal(c, "failed to update recipe")
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *recipeHandler) archive(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	existing, err := h.db.GetRecipe(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "recipe not found")
			return
		}
		httpx.ErrInternal(c, "failed to get recipe")
		return
	}
	if existing.UserID != userID {
		httpx.ErrForbidden(c, "not your recipe")
		return
	}
	if existing.Status == models.StatusArchived {
		c.JSON(http.StatusOK, existing)
		return
	}
	archived, err := h.db.ArchiveRecipe(c.Request.Context(), id)
	if err != nil {
		httpx.ErrInternal(c, "failed to archive recipe")
		return
	}
	c.JSON(http.StatusOK, archived)
}

func (h *recipeHandler) save(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	recipe, err := h.db.GetRecipe(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "recipe not found")
			return
		}
		httpx.ErrInternal(c, "failed to get recipe")
		return
	}
	if recipe.Status == models.StatusArchived {
		httpx.ErrBadRequest(c, "cannot save archived recipe")
		return
	}
	if recipe.Status != models.StatusPublished && recipe.UserID != userID {
		httpx.ErrNotFound(c, "recipe not found")
		return
	}
	if err := h.db.SaveRecipe(c.Request.Context(), userID, id); err != nil {
		httpx.ErrInternal(c, "failed to save")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"saved": true})
}

func (h *recipeHandler) unsave(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	if err := h.db.UnsaveRecipe(c.Request.Context(), userID, id); err != nil {
		httpx.ErrInternal(c, "failed to unsave")
		return
	}
	c.JSON(http.StatusOK, gin.H{"saved": false})
}

func (h *recipeHandler) uploadURL(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	var req struct {
		RecipeID string `json:"recipe_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.RecipeID == "" {
		httpx.ErrBadRequest(c, "recipe_id is required")
		return
	}
	recipe, err := h.db.GetRecipe(c.Request.Context(), req.RecipeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "recipe not found")
			return
		}
		httpx.ErrInternal(c, "failed to get recipe")
		return
	}
	if recipe.UserID != userID {
		httpx.ErrForbidden(c, "not your recipe")
		return
	}
	if recipe.Status == models.StatusArchived {
		httpx.ErrBadRequest(c, "archived recipes cannot get upload url")
		return
	}
	uid, url, err := h.video.CreateDirectUploadURL(c.Request.Context())
	if err != nil {
		httpx.ErrInternal(c, "failed to create upload url")
		return
	}
	status := recipe.Status
	if status != models.StatusPublished {
		status = models.StatusProcessing
	}
	updated, err := h.db.SetRecipeVideo(c.Request.Context(), req.RecipeID, uid, string(status))
	if err != nil {
		httpx.ErrInternal(c, "failed to update recipe")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"video_uid":  uid,
		"upload_url": url,
		"recipe":     updated,
	})
}

func (h *recipeHandler) fork(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	recipe, err := h.db.GetRecipe(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "recipe not found")
			return
		}
		httpx.ErrInternal(c, "failed to get recipe")
		return
	}
	if recipe.Status == models.StatusArchived {
		httpx.ErrBadRequest(c, "cannot fork archived recipe")
		return
	}
	if recipe.Status != models.StatusPublished && recipe.UserID != userID {
		httpx.ErrNotFound(c, "recipe not found")
		return
	}
	forked, err := h.db.ForkRecipe(c.Request.Context(), id, userID)
	if err != nil {
		httpx.ErrInternal(c, "failed to fork")
		return
	}
	c.JSON(http.StatusCreated, forked)
}

func (h *recipeHandler) publish(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	recipe, err := h.db.GetRecipe(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "recipe not found")
			return
		}
		httpx.ErrInternal(c, "failed to get recipe")
		return
	}
	if recipe.UserID != userID {
		httpx.ErrForbidden(c, "not your recipe")
		return
	}
	if recipe.Status == models.StatusPublished {
		c.JSON(http.StatusOK, recipe)
		return
	}
	if recipe.Status == models.StatusArchived {
		httpx.ErrBadRequest(c, "archived recipes cannot be published")
		return
	}
	if msg := httpx.ValidateRecipe(recipe); msg != "" {
		httpx.ErrBadRequest(c, "cannot publish: "+msg)
		return
	}
	published, err := h.db.PublishRecipe(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrBadRequest(c, "cannot publish in current status")
			return
		}
		httpx.ErrInternal(c, "failed to publish")
		return
	}
	if h.indexer != nil {
		_ = h.indexer.IndexRecipe(c.Request.Context(), published)
	}
	c.JSON(http.StatusOK, published)
}
