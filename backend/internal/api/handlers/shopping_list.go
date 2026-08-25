package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"cookmode/internal/api/middleware"
	"cookmode/internal/services"
)

// RegisterShoppingLists mounts shopping list CRUD with recipe scaling
// and nutrition lookup via the nutrition service.
func RegisterShoppingLists(
	rg *gin.RouterGroup,
	pool *pgxpool.Pool,
	nutrition *services.NutritionService,
	auth *middleware.Auth,
) {
	h := &shoppingListHandler{db: pool, nutrition: nutrition}

	lists := rg.Group("/shopping-lists", auth.Middleware())
	{
		lists.GET("", h.list)
		lists.POST("", h.create)
		lists.POST("/:id/items", h.addItems)
		lists.PATCH("/:id/items/:itemIndex", h.checkItem)
	}
}

type shoppingListHandler struct {
	db        *pgxpool.Pool
	nutrition *services.NutritionService
}

func (h *shoppingListHandler) list(c *gin.Context)      { c.Status(501) } // TODO
func (h *shoppingListHandler) create(c *gin.Context)    { c.Status(501) } // TODO: build from recipe(s), scale quantities to servings
func (h *shoppingListHandler) addItems(c *gin.Context)  { c.Status(501) } // TODO
func (h *shoppingListHandler) checkItem(c *gin.Context) { c.Status(501) } // TODO: toggle checked
