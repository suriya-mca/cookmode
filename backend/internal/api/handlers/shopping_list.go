package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"cookmode/internal/api/middleware"
	"cookmode/internal/db/queries"
	"cookmode/internal/httpx"
	"cookmode/internal/models"
	"cookmode/internal/services"
)

func RegisterShoppingLists(
	rg *gin.RouterGroup,
	pool *pgxpool.Pool,
	nutrition *services.NutritionService,
	auth *middleware.Auth,
) {
	h := &shoppingListHandler{db: queries.New(pool), pool: pool, nutrition: nutrition}
	lists := rg.Group("/shopping-lists", auth.Middleware())
	{
		lists.GET("", h.list)
		lists.POST("", h.create)
		lists.GET("/:id", h.get)
		lists.DELETE("/:id", h.delete)
		lists.POST("/:id/items", h.addItems)
		lists.PATCH("/:id/items/:itemIndex", h.checkItem)
	}
}

type shoppingListHandler struct {
	db        *queries.DB
	pool      *pgxpool.Pool
	nutrition *services.NutritionService
}

func (h *shoppingListHandler) list(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	lists, err := h.db.ListShoppingLists(c.Request.Context(), userID)
	if err != nil {
		httpx.ErrInternal(c, "failed to list")
		return
	}
	if lists == nil {
		lists = []*models.ShoppingList{}
	}
	c.JSON(http.StatusOK, gin.H{"data": lists})
}

func (h *shoppingListHandler) get(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	l, err := h.db.GetShoppingList(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "shopping list not found")
			return
		}
		httpx.ErrInternal(c, "failed to get")
		return
	}
	c.JSON(http.StatusOK, l)
}

func (h *shoppingListHandler) delete(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	if err := h.db.DeleteShoppingList(c.Request.Context(), id, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "shopping list not found")
			return
		}
		httpx.ErrInternal(c, "failed to delete")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *shoppingListHandler) create(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	var req struct {
		Title     string   `json:"title"`
		RecipeIDs []string `json:"recipe_ids"`
		Servings  *int     `json:"servings"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrBadRequest(c, "invalid body")
		return
	}
	var items []models.ShoppingListItem
	for _, rid := range req.RecipeIDs {
		recipe, err := h.db.GetRecipe(c.Request.Context(), rid)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				httpx.ErrNotFound(c, "recipe not found: "+rid)
				return
			}
			httpx.ErrInternal(c, "failed to get recipe")
			return
		}
		factor := 1.0
		if req.Servings != nil && recipe.Servings > 0 {
			factor = float64(*req.Servings) / float64(recipe.Servings)
		}
		for _, ing := range recipe.Ingredients {
			items = append(items, models.ShoppingListItem{
				IngredientName: ing.Name,
				Quantity:       ing.Quantity * factor,
				Unit:           ing.Unit,
			})
		}
	}
	title := req.Title
	if title == "" {
		title = "Shopping List"
	}
	list, err := h.db.CreateShoppingList(c.Request.Context(), &models.ShoppingList{
		UserID: userID,
		Title:  title,
		Items:  items,
	})
	if err != nil {
		httpx.ErrInternal(c, "failed to create")
		return
	}
	c.JSON(http.StatusCreated, list)
}

func (h *shoppingListHandler) addItems(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	var req struct {
		Items []models.ShoppingListItem `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Items) == 0 {
		httpx.ErrBadRequest(c, "items are required")
		return
	}
	for i, it := range req.Items {
		if it.IngredientName == "" {
			httpx.ErrBadRequest(c, "items["+strconv.Itoa(i)+"].ingredient_name is required")
			return
		}
		if it.Quantity < 0 {
			httpx.ErrBadRequest(c, "quantity must be >= 0")
			return
		}
	}
	list, err := h.db.GetShoppingList(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "shopping list not found")
			return
		}
		httpx.ErrInternal(c, "failed to get")
		return
	}
	list.Items = append(list.Items, req.Items...)
	updated, err := h.db.UpdateShoppingListItems(c.Request.Context(), id, userID, list.Items)
	if err != nil {
		httpx.ErrInternal(c, "failed to update")
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *shoppingListHandler) checkItem(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	idx, err := strconv.Atoi(c.Param("itemIndex"))
	if err != nil {
		httpx.ErrBadRequest(c, "invalid itemIndex")
		return
	}
	var req struct {
		Checked *bool `json:"checked"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Checked == nil {
		httpx.ErrBadRequest(c, "checked is required")
		return
	}
	list, err := h.db.GetShoppingList(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "shopping list not found")
			return
		}
		httpx.ErrInternal(c, "failed to get")
		return
	}
	if idx < 0 || idx >= len(list.Items) {
		httpx.ErrBadRequest(c, "itemIndex out of range")
		return
	}
	list.Items[idx].Checked = *req.Checked
	updated, err := h.db.UpdateShoppingListItems(c.Request.Context(), id, userID, list.Items)
	if err != nil {
		httpx.ErrInternal(c, "failed to update")
		return
	}
	c.JSON(http.StatusOK, updated)
}
