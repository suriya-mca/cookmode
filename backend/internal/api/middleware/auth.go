package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// tokenTTL is how long issued tokens stay valid. There are no refresh
// tokens — clients re-login when the token expires.
const tokenTTL = 7 * 24 * time.Hour

// Auth verifies locally-issued HS256 JWTs (signed with the configured
// secret — see config.Validate) and injects the user id into the request
// context.
type Auth struct {
	secret []byte
}

// NewAuth creates a token issuer and verifier using the configured signing secret.
func NewAuth(secret string) *Auth {
	return &Auth{secret: []byte(secret)}
}

// parse validates a bearer token and returns its user id. The signing
// method is pinned to HS256, expiry is required, and sub must be a UUID.
func (a *Auth) parse(token string) (string, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return a.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil || !parsed.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	sub, _ := claims["sub"].(string)
	if _, err := uuid.Parse(sub); err != nil {
		return "", errors.New("invalid sub claim")
	}
	return sub, nil
}

// Middleware requires a valid bearer token and stores its user ID in the context.
func (a *Auth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		sub, err := a.parse(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("user_id", sub)
		c.Next()
	}
}

// Optional sets user_id when a valid bearer token is supplied. Unlike the
// strict middleware it lets requests without any Authorization header
// through as anonymous — but a present-yet-invalid token is rejected, so
// stale credentials surface as 401 instead of silently anonymous behavior.
func (a *Auth) Optional() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || token == "" {
			c.Next()
			return
		}
		sub, err := a.parse(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("user_id", sub)
		c.Next()
	}
}

// IssueToken creates a HS256 JWT for the given userID using the same secret
// that Middleware() verifies. This keeps hand-minted tokens compatible.
func (a *Auth) IssueToken(userID string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": jwt.NewNumericDate(now),
		"exp": jwt.NewNumericDate(now.Add(tokenTTL)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secret)
}

// UserIDFrom returns the authenticated user id set by Auth.Middleware.
func UserIDFrom(c *gin.Context) string {
	v, _ := c.Get("user_id")
	id, _ := v.(string)
	return id
}

// CORS allows cross-origin requests from the configured allowlist (from
// Config.CORSAllowOrigins). A single "*" entry allows any origin — dev
// only. Preflight OPTIONS requests are answered directly.
func CORS(allowed []string) gin.HandlerFunc {
	wildcard := len(allowed) == 1 && allowed[0] == "*"
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		// Always vary on Origin so shared caches don't serve one origin's
		// allow-origin header to another.
		c.Header("Vary", "Origin")
		allow := ""
		if wildcard {
			allow = "*"
		} else {
			for _, o := range allowed {
				if o != "" && o == origin {
					allow = origin
					break
				}
			}
		}
		if allow != "" {
			c.Header("Access-Control-Allow-Origin", allow)
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
