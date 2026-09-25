package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"cookmode/internal/api/middleware"
	"cookmode/internal/db/queries"
	"cookmode/internal/httpx"
)

// dummyHash is compared against on unknown-user / passwordless logins so a
// bcrypt comparison (~DefaultCost) always runs and the failure path takes
// the same time as a real password check.
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("cookmode-nonexistent-user-credential"), bcrypt.DefaultCost)

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

// signup validates credentials, creates a local user, and returns a session token.
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
	username, msg := httpx.NormalizeUsername(req.Username)
	if msg != "" {
		httpx.ErrBadRequest(c, msg)
		return
	}
	if msg := httpx.ValidatePassword(req.Password); msg != "" {
		httpx.ErrBadRequest(c, msg)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		httpx.ErrInternal(c, "failed to hash password")
		return
	}
	id := uuid.NewString()
	user, err := h.db.CreateAuthUser(c.Request.Context(), id, username, string(hash), req.DisplayName)
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

// login verifies a local password and returns a token for valid credentials.
func (h *authHandler) login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrBadRequest(c, "invalid body")
		return
	}
	username, msg := httpx.NormalizeUsername(req.Username)
	if msg != "" {
		httpx.ErrBadRequest(c, msg)
		return
	}
	if req.Password == "" {
		httpx.ErrBadRequest(c, "username and password are required")
		return
	}
	user, hash, err := h.db.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Flatten timing: burn one bcrypt comparison before the 401.
			_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(req.Password))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		httpx.ErrInternal(c, "failed to lookup user")
		return
	}
	if hash == "" {
		// Legacy/imported rows have no local password — same 401, same cost.
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(req.Password))
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
