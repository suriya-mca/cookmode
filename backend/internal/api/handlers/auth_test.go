package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cookmode/internal/api/middleware"
)

const handlerTestSecret = "handler-test-secret-at-least-32-chars!!"

// testRouter builds a minimal /api/v1 router with the auth and users
// endpoints backed by the real database from DATABASE_URL. It skips when
// DATABASE_URL is unset (e.g. plain `go test ./...` without a DB); CI
// exports it and runs migrations before the tests.
func testRouter(t *testing.T) (*gin.Engine, *pgxpool.Pool) {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set — skipping DB-backed auth tests")
	}
	gin.SetMode(gin.TestMode)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("connect db: %v", err)
	}
	t.Cleanup(pool.Close)
	if err := pool.Ping(context.Background()); err != nil {
		t.Skipf("database unreachable: %v", err)
	}
	auth := middleware.NewAuth(handlerTestSecret)
	router := gin.New()
	v1 := router.Group("/api/v1")
	RegisterAuth(v1, pool, auth)
	RegisterUsers(v1, pool, auth)
	return router, pool
}

// doJSON sends a JSON request through the test router with an optional bearer token.
func doJSON(t *testing.T, router *gin.Engine, method, path string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// decodeBody decodes a recorded JSON response and fails the test on invalid JSON.
func decodeBody(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("decode %q: %v", w.Body.String(), err)
	}
	return m
}

// uniqueUsername derives a policy-valid username from a UUID (hex only) and
// registers cleanup so repeated local runs never collide or leave rows behind.
func uniqueUsername(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	u := "tu" + strings.ReplaceAll(uuid.NewString()[:16], "-", "")
	t.Cleanup(func() { deleteUser(t, pool, u) })
	return u
}

// deleteUser removes a test account by username during cleanup.
func deleteUser(t *testing.T, pool *pgxpool.Pool, username string) {
	t.Helper()
	_, _ = pool.Exec(context.Background(), "DELETE FROM users WHERE username=$1", username)
}

// TestSignupLoginFlow covers signup, duplicate signup, login, and failed logins.
func TestSignupLoginFlow(t *testing.T) {
	router, pool := testRouter(t)
	username := uniqueUsername(t, pool)
	password := "correct-horse-9"

	// Signup → 201 with token + user.
	w := doJSON(t, router, http.MethodPost, "/api/v1/auth/signup", map[string]string{
		"username": username, "password": password,
	}, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("signup code %d (body %s), want 201", w.Code, w.Body.String())
	}
	res := decodeBody(t, w)
	token, _ := res["token"].(string)
	if token == "" {
		t.Fatalf("signup returned no token: %v", res)
	}
	user, _ := res["user"].(map[string]any)
	if user["username"] != strings.ToLower(username) {
		t.Fatalf("signup user %v, want username %q", user, strings.ToLower(username))
	}

	// Duplicate signup → 409.
	w = doJSON(t, router, http.MethodPost, "/api/v1/auth/signup", map[string]string{
		"username": username, "password": password,
	}, "")
	if w.Code != http.StatusConflict {
		t.Fatalf("duplicate signup code %d, want 409", w.Code)
	}

	// Login → 200 with a fresh token that works on /users/me.
	w = doJSON(t, router, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"username": username, "password": password,
	}, "")
	if w.Code != http.StatusOK {
		t.Fatalf("login code %d (body %s), want 200", w.Code, w.Body.String())
	}
	loginToken, _ := decodeBody(t, w)["token"].(string)
	w = doJSON(t, router, http.MethodGet, "/api/v1/users/me", nil, loginToken)
	if w.Code != http.StatusOK {
		t.Fatalf("users/me code %d (body %s), want 200", w.Code, w.Body.String())
	}

	// Wrong password → 401 with the flat "invalid credentials" body.
	wrongKnownStart := time.Now()
	w = doJSON(t, router, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"username": username, "password": "wrong-password-1",
	}, "")
	if w.Code != http.StatusUnauthorized || !strings.Contains(w.Body.String(), "invalid credentials") {
		t.Fatalf("wrong password code %d body %s, want 401 invalid credentials", w.Code, w.Body.String())
	}

	// Unknown user → same 401 body. The dummy bcrypt compare should make
	// this path cost about the same as the real wrong-password path above,
	// so assert relatively instead of against a machine-dependent floor.
	unknown := uniqueUsername(t, pool)
	timeWrongKnown := time.Since(wrongKnownStart)
	start := time.Now()
	w = doJSON(t, router, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"username": unknown, "password": "wrong-password-1",
	}, "")
	elapsedUnknown := time.Since(start)
	if w.Code != http.StatusUnauthorized || !strings.Contains(w.Body.String(), "invalid credentials") {
		t.Fatalf("unknown user code %d body %s, want 401 invalid credentials", w.Code, w.Body.String())
	}
	if elapsedUnknown < timeWrongKnown/2 {
		t.Fatalf("unknown-user login took %v vs wrong-password %v — dummy bcrypt compare may be missing",
			elapsedUnknown, timeWrongKnown)
	}
}

// TestSignupValidation checks username normalization and credential limits.
func TestSignupValidation(t *testing.T) {
	router, pool := testRouter(t)
	// Uppercase input; the row is stored lowercased, so clean up by that.
	upper := "ValidUser" + strings.ToLower(uniqueUsername(t, pool)[2:])
	t.Cleanup(func() { deleteUser(t, pool, strings.ToLower(upper)) })
	cases := []struct {
		name     string
		username string
		password string
		wantCode int
	}{
		{"too short", "ab", "valid-password-1", http.StatusBadRequest},
		{"uppercase normalized", upper, "valid-password-1", http.StatusCreated},
		{"bad chars", "not a user!", "valid-password-1", http.StatusBadRequest},
		{"empty", "", "valid-password-1", http.StatusBadRequest},
		{"short password", uniqueUsername(t, pool), "short", http.StatusBadRequest},
		{"overlong password", uniqueUsername(t, pool), strings.Repeat("p", 73), http.StatusBadRequest},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := doJSON(t, router, http.MethodPost, "/api/v1/auth/signup", map[string]string{
				"username": tc.username, "password": tc.password,
			}, "")
			if w.Code != tc.wantCode {
				t.Fatalf("code %d (body %s), want %d", w.Code, w.Body.String(), tc.wantCode)
			}
			// Uppercase input must be stored lowercase.
			if tc.wantCode == http.StatusCreated {
				user, _ := decodeBody(t, w)["user"].(map[string]any)
				if user["username"] != strings.ToLower(tc.username) {
					t.Fatalf("stored username %v, want lowercase %q", user["username"], strings.ToLower(tc.username))
				}
			}
		})
	}
}

// TestUpdateMeRejectsEmptyUsername checks profile username validation and updates.
func TestUpdateMeRejectsEmptyUsername(t *testing.T) {
	router, pool := testRouter(t)
	username := uniqueUsername(t, pool)
	w := doJSON(t, router, http.MethodPost, "/api/v1/auth/signup", map[string]string{
		"username": username, "password": "valid-password-1",
	}, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("signup code %d, want 201", w.Code)
	}
	token, _ := decodeBody(t, w)["token"].(string)

	// Empty username must be a 400, not a silent NULL that breaks login.
	w = doJSON(t, router, http.MethodPatch, "/api/v1/users/me", map[string]string{"username": "   "}, token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("empty username patch code %d (body %s), want 400", w.Code, w.Body.String())
	}

	// Valid patch still works and is normalized (mixed-case input is
	// stored lowercase). Note PATCH returns the user object directly,
	// not a {user} envelope.
	suffix := strings.ReplaceAll(uuid.NewString()[:12], "-", "")
	target := "newname" + suffix
	t.Cleanup(func() { deleteUser(t, pool, target) })
	w = doJSON(t, router, http.MethodPatch, "/api/v1/users/me", map[string]string{"username": "NewName" + suffix}, token)
	if w.Code != http.StatusOK {
		t.Fatalf("valid patch code %d (body %s), want 200", w.Code, w.Body.String())
	}
	if got := decodeBody(t, w)["username"]; got != target {
		t.Fatalf("patched username %v, want %s", got, target)
	}
}
