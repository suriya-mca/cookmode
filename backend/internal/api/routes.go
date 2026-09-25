package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"

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
	River            *river.Client[pgx.Tx]
	JWTSecret        string
	// CORSAllowOrigins lists web origins allowed to call the API
	// (e.g. Expo web dev server). Empty = no cross-origin access.
	CORSAllowOrigins []string
}

// NewRouter builds the Gin engine with all /api/v1 routes.
func NewRouter(d Deps) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORS(d.CORSAllowOrigins))

	auth := middleware.NewAuth(d.JWTSecret)

	v1 := router.Group("/api/v1")
	handlers.RegisterHealth(v1)
	handlers.RegisterAuth(v1, d.Pool, auth)
	handlers.RegisterUsers(v1, d.Pool, auth)
	handlers.RegisterRecipes(v1, d.Pool, d.SearchIndexer, d.VideoService, d.River, auth)
	handlers.RegisterCollections(v1, d.Pool, auth)
	handlers.RegisterFollows(v1, d.Pool, auth)
	handlers.RegisterPosts(v1, d.Pool, auth)
	handlers.RegisterShoppingLists(v1, d.Pool, d.NutritionService, auth)
	handlers.RegisterSearch(v1, d.SearchIndexer)

	return router
}
