package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const testSecret = "test-secret-that-is-long-enough-32-chars"

// init sets Gin to test mode before middleware tests create contexts.
func init() {
	gin.SetMode(gin.TestMode)
}

// runThrough feeds a request through a single middleware and returns the
// recorder plus the gin context (to inspect user_id).
func runThrough(h gin.HandlerFunc, bearer string) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	c.Request = req
	h(c)
	return w, c
}

// mint signs a test token with the requested algorithm, secret, and claims.
func mint(t *testing.T, method jwt.SigningMethod, secret string, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(method, claims)
	s, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

// validClaims returns unexpired JWT claims for a given subject.
func validClaims(sub string) jwt.MapClaims {
	return jwt.MapClaims{
		"sub": sub,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
}

// TestIssueTokenRoundTrip checks that an issued token parses to its user ID.
func TestIssueTokenRoundTrip(t *testing.T) {
	a := NewAuth(testSecret)
	id := uuid.NewString()
	token, err := a.IssueToken(id)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}
	got, err := a.parse(token)
	if err != nil {
		t.Fatalf("parse issued token: %v", err)
	}
	if got != id {
		t.Fatalf("got sub %q, want %q", got, id)
	}
}

// TestIssueTokenTTL checks the expiry interval encoded in an issued token.
func TestIssueTokenTTL(t *testing.T) {
	a := NewAuth(testSecret)
	token, err := a.IssueToken(uuid.NewString())
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}
	parsed, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("ParseUnverified: %v", err)
	}
	claims := parsed.Claims.(jwt.MapClaims)
	ttl := int64(claims["exp"].(float64)) - int64(claims["iat"].(float64))
	if ttl != int64(tokenTTL.Seconds()) {
		t.Fatalf("token ttl %d, want %d", ttl, int64(tokenTTL.Seconds()))
	}
}

// TestMiddlewareStrict checks acceptance of valid tokens and rejection of invalid ones.
func TestMiddlewareStrict(t *testing.T) {
	a := NewAuth(testSecret)
	id := uuid.NewString()
	good := mint(t, jwt.SigningMethodHS256, testSecret, validClaims(id))

	cases := []struct {
		name       string
		bearer     string // "" = no Authorization header
		wantCode   int
		wantUserID string
	}{
		{"missing header", "", http.StatusUnauthorized, ""},
		{"valid", good, http.StatusOK, id},
		{"wrong secret", mint(t, jwt.SigningMethodHS256, "another-secret-that-is-also-long-enough", validClaims(id)), http.StatusUnauthorized, ""},
		{"hs512 rejected", mint(t, jwt.SigningMethodHS512, testSecret, validClaims(id)), http.StatusUnauthorized, ""},
		{"hs384 rejected", mint(t, jwt.SigningMethodHS384, testSecret, validClaims(id)), http.StatusUnauthorized, ""},
		{"missing exp", mint(t, jwt.SigningMethodHS256, testSecret, jwt.MapClaims{"sub": id}), http.StatusUnauthorized, ""},
		{"expired", mint(t, jwt.SigningMethodHS256, testSecret, jwt.MapClaims{"sub": id, "exp": time.Now().Add(-time.Hour).Unix()}), http.StatusUnauthorized, ""},
		{"non-uuid sub", mint(t, jwt.SigningMethodHS256, testSecret, validClaims("not-a-uuid")), http.StatusUnauthorized, ""},
		{"garbage", "this.is.not-a-token", http.StatusUnauthorized, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w, c := runThrough(a.Middleware(), tc.bearer)
			if w.Code != tc.wantCode {
				t.Fatalf("code %d, want %d (body %s)", w.Code, tc.wantCode, w.Body.String())
			}
			if got := UserIDFrom(c); got != tc.wantUserID {
				t.Fatalf("user_id %q, want %q", got, tc.wantUserID)
			}
		})
	}
}

// TestOptional checks anonymous access and rejection of invalid supplied tokens.
func TestOptional(t *testing.T) {
	a := NewAuth(testSecret)
	id := uuid.NewString()
	good := mint(t, jwt.SigningMethodHS256, testSecret, validClaims(id))
	bad := mint(t, jwt.SigningMethodHS256, testSecret, validClaims("not-a-uuid"))

	// No header → anonymous, passes through.
	w, c := runThrough(a.Optional(), "")
	if w.Code != http.StatusOK || UserIDFrom(c) != "" {
		t.Fatalf("no header: code=%d user=%q, want 200/empty", w.Code, UserIDFrom(c))
	}
	// Valid → user set.
	w, c = runThrough(a.Optional(), good)
	if w.Code != http.StatusOK || UserIDFrom(c) != id {
		t.Fatalf("valid: code=%d user=%q, want 200/%q", w.Code, UserIDFrom(c), id)
	}
	// Present-but-invalid → 401 (fail closed, not silent anonymous).
	w, _ = runThrough(a.Optional(), bad)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("invalid bearer: code=%d, want 401", w.Code)
	}
	w, _ = runThrough(a.Optional(), "garbage-token")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("garbage bearer: code=%d, want 401", w.Code)
	}
}

// TestCORS checks allowed and unlisted origins, preflight, and wildcard access.
func TestCORS(t *testing.T) {
	h := CORS([]string{"http://localhost:8081"})

	// Allowed origin gets headers and passes through.
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Origin", "http://localhost:8081")
	c.Request = req
	h(c)
	if c.IsAborted() {
		t.Fatalf("allowed origin request was aborted")
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:8081" {
		t.Fatalf("missing allow-origin header: %v", w.Header())
	}

	// Unknown origin gets no CORS headers.
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	req2 := httptest.NewRequest(http.MethodGet, "/x", nil)
	req2.Header.Set("Origin", "https://evil.example")
	c2.Request = req2
	h(c2)
	if v := w2.Header().Get("Access-Control-Allow-Origin"); v != "" {
		t.Fatalf("unexpected allow-origin for unknown origin: %q", v)
	}

	// Preflight is answered directly.
	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)
	req3 := httptest.NewRequest(http.MethodOptions, "/x", nil)
	req3.Header.Set("Origin", "http://localhost:8081")
	c3.Request = req3
	h(c3)
	if !c3.IsAborted() || w3.Code != http.StatusNoContent {
		t.Fatalf("preflight: aborted=%v code=%d, want true/204", c3.IsAborted(), w3.Code)
	}

	// Wildcard allows any origin.
	wh := CORS([]string{"*"})
	w4 := httptest.NewRecorder()
	c4, _ := gin.CreateTestContext(w4)
	req4 := httptest.NewRequest(http.MethodGet, "/x", nil)
	req4.Header.Set("Origin", "https://anything.example")
	c4.Request = req4
	wh(c4)
	if v := w4.Header().Get("Access-Control-Allow-Origin"); v != "*" {
		t.Fatalf("wildcard: allow-origin=%q, want *", v)
	}
}
