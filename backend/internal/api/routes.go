package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"cookmode/internal/api/handlers"
	"cookmode/internal/api/middleware"
	"cookmode/internal/services"
)

// Deps carries everything the route tree needs.
type Deps struct {
	Pool             *pgxpool.Pool
	SearchIndexer    *services.SearchIndexer
	VideoService     *services.VideoService
	NutritionService *services.NutritionService
	JWTSecret        string
}

// NewRouter builds the Gin engine with all /api/v1 routes.
func NewRouter(d Deps) *gin.Engine {
	router := gin.Default()

	auth := middleware.NewAuth(d.JWTSecret)

	v1 := router.Group("/api/v1")
	handlers.RegisterHealth(v1)
	handlers.RegisterRecipes(v1, d.Pool, d.SearchIndexer, d.VideoService, auth)
	handlers.RegisterCollections(v1, d.Pool, auth)
	handlers.RegisterFollows(v1, d.Pool, auth)
	handlers.RegisterPosts(v1, d.Pool, auth)
	handlers.RegisterShoppingLists(v1, d.Pool, d.NutritionService, auth)
	handlers.RegisterSearch(v1, d.SearchIndexer)

	return router
}
