package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"cookmode/internal/api/middleware"
	"cookmode/internal/services"
)

// RegisterRecipes mounts /recipes routes. Recipe creation enqueues the async
// pipeline (transcribe -> infer anchors -> fetch nutrition -> index).
func RegisterRecipes(
	rg *gin.RouterGroup,
	pool *pgxpool.Pool,
	indexer *services.SearchIndexer,
	video *services.VideoService,
	auth *middleware.Auth,
) {
	h := &recipeHandler{db: pool, indexer: indexer, video: video}

	recipes := rg.Group("/recipes")
	{
		recipes.GET("", h.list)
		recipes.GET("/:id", h.get)

		protected := recipes.Group("", auth.Middleware())
		{
			protected.POST("", h.create)
			protected.PATCH("/:id", h.update)
			protected.DELETE("/:id", h.archive)
			protected.POST("/:id/saves", h.save)
			protected.POST("/upload-url", h.uploadURL)
		}
	}
}

type recipeHandler struct {
	db      *pgxpool.Pool
	indexer *services.SearchIndexer
	video   *services.VideoService
}

func (h *recipeHandler) list(c *gin.Context) { c.Status(501) } // TODO
func (h *recipeHandler) get(c *gin.Context)  { c.Status(501) } // TODO

func (h *recipeHandler) create(c *gin.Context) { c.Status(501) } // TODO: insert draft, request video upload, enqueue processing jobs

func (h *recipeHandler) update(c *gin.Context) { c.Status(501) } // TODO

func (h *recipeHandler) archive(c *gin.Context) { c.Status(501) } // TODO

func (h *recipeHandler) save(c *gin.Context) { c.Status(501) } // TODO: increment saves counter

func (h *recipeHandler) uploadURL(c *gin.Context) { c.Status(501) } // TODO: return Cloudflare Stream direct-upload URL
