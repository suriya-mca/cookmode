package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"cookmode/internal/api/middleware"
	"cookmode/internal/db/queries"
	"cookmode/internal/httpx"
	"cookmode/internal/models"
)

// RegisterUsers registers authenticated current-user endpoints and a public user lookup endpoint.
func RegisterUsers(rg *gin.RouterGroup, pool *pgxpool.Pool, auth *middleware.Auth) {
	h := &userHandler{db: queries.New(pool)}
	me := rg.Group("/users/me", auth.Middleware())
	{
		me.GET("", h.getMe)
		me.PATCH("", h.updateMe)
	}
	rg.GET("/users/:id", h.getByID)
}

type userHandler struct {
	db *queries.DB
}

func (h *userHandler) getMe(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	u, err := h.db.GetUser(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			u = &models.User{ID: userID}
		} else {
			httpx.ErrInternal(c, "failed to get user")
			return
		}
	}
	c.JSON(http.StatusOK, u)
}

func (h *userHandler) getByID(c *gin.Context) {
	id := c.Param("id")
	u, err := h.db.GetUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.ErrNotFound(c, "user not found")
			return
		}
		httpx.ErrInternal(c, "failed to get user")
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *userHandler) updateMe(c *gin.Context) {
	userID := middleware.UserIDFrom(c)
	var req struct {
		Username    *string `json:"username"`
		DisplayName *string `json:"display_name"`
		AvatarURL   *string `json:"avatar_url"`
		Bio         *string `json:"bio"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrBadRequest(c, "invalid body")
		return
	}
	if req.Username != nil {
		// An empty or invalid username is a 400: NULLing the username
		// would permanently lock the account out of password login.
		u, msg := httpx.NormalizeUsername(*req.Username)
		if msg != "" {
			httpx.ErrBadRequest(c, msg)
			return
		}
		req.Username = &u
	}
	if req.Bio != nil {
		if len(*req.Bio) > 500 {
			httpx.ErrBadRequest(c, "bio must be at most 500 characters")
			return
		}
	}
	updated, err := h.db.PatchUser(c.Request.Context(), userID, req.Username, req.DisplayName, req.AvatarURL, req.Bio)
	if err != nil {
		log.Printf("upsert user error: %v", err)
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
			httpx.ErrBadRequest(c, "username already taken")
			return
		}
		httpx.ErrInternal(c, "failed to update user")
		return
	}
	c.JSON(http.StatusOK, updated)
}
