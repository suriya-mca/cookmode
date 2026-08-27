package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"cookmode/internal/api/middleware"
	"cookmode/internal/db/queries"
	"cookmode/internal/httpx"
)

func RegisterFollows(rg *gin.RouterGroup, pool *pgxpool.Pool, auth *middleware.Auth) {
	h := &followHandler{db: queries.New(pool)}
	f := rg.Group("/follows", auth.Middleware())
	{
		f.POST("", h.follow)
		f.DELETE("/:id", h.unfollow)
	}
	u := rg.Group("/users")
	{
		u.GET("/:id/followers", h.followers)
		u.GET("/:id/following", h.following)
	}
}

type followHandler struct {
	db *queries.DB
}

func (h *followHandler) follow(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	var req struct {
		FollowingID string `json:"following_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.FollowingID == "" {
		httpx.ErrBadRequest(c, "following_id is required")
		return
	}
	if req.FollowingID == userID {
		httpx.ErrBadRequest(c, "cannot follow yourself")
		return
	}
	if err := h.db.Follow(c.Request.Context(), userID, req.FollowingID); err != nil {
		httpx.ErrInternal(c, "failed to follow")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"following_id": req.FollowingID})
}

func (h *followHandler) unfollow(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	id := c.Param("id")
	if err := h.db.Unfollow(c.Request.Context(), userID, id); err != nil {
		httpx.ErrInternal(c, "failed to unfollow")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *followHandler) followers(c *gin.Context) {
	id := c.Param("id")
	list, err := h.db.ListFollowers(c.Request.Context(), id)
	if err != nil {
		httpx.ErrInternal(c, "failed to list followers")
		return
	}
	if list == nil {
		list = []string{}
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *followHandler) following(c *gin.Context) {
	id := c.Param("id")
	list, err := h.db.ListFollowing(c.Request.Context(), id)
	if err != nil {
		httpx.ErrInternal(c, "failed to list following")
		return
	}
	if list == nil {
		list = []string{}
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}
