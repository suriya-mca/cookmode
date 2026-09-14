package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"cookmode/internal/api/middleware"
	"cookmode/internal/db/queries"
	"cookmode/internal/httpx"
)

func RegisterAuth(rg *gin.RouterGroup, pool *pgxpool.Pool, auth *middleware.Auth) {
	h := &authHandler{db: queries.New(pool), auth: auth}
	g := rg.Group("/auth")
	{
		g.POST("/signup", h.signup)
		g.POST("/login", h.login)
	}
}

type authHandler struct {
	db   *queries.DB
	auth *middleware.Auth
}

func (h *authHandler) signup(c *gin.Context) {
	var req struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrBadRequest(c, "invalid body")
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	if len(req.Username) < 3 || len(req.Username) > 30 {
		httpx.ErrBadRequest(c, "username must be 3-30 characters")
		return
	}
	if len(req.Password) < 8 {
		httpx.ErrBadRequest(c, "password must be at least 8 characters")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		httpx.ErrInternal(c, "failed to hash password")
		return
	}
	id := uuid.NewString()
	user, err := h.db.CreateAuthUser(c.Request.Context(), id, req.Username, string(hash), req.DisplayName)
	if err != nil {
		if errors.Is(err, queries.ErrUsernameTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
			return
		}
		httpx.ErrInternal(c, "failed to create user")
		return
	}
	token, err := h.auth.IssueToken(user.ID)
	if err != nil {
		httpx.ErrInternal(c, "failed to issue token")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}

func (h *authHandler) login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrBadRequest(c, "invalid body")
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		httpx.ErrBadRequest(c, "username and password are required")
		return
	}
	user, hash, err := h.db.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		httpx.ErrInternal(c, "failed to lookup user")
		return
	}
	if hash == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := h.auth.IssueToken(user.ID)
	if err != nil {
		httpx.ErrInternal(c, "failed to issue token")
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}
