package config

import (
	"strings"
	"testing"
)

// TestValidate checks rejection of weak JWT secrets and acceptance of valid ones.
func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{"empty", "", true},
		{"short", "too-short-secret", true},
		{"31 chars", strings.Repeat("a", 31), true},
		{"legacy supabase placeholder", legacyPlaceholderJWTSecret, true},
		{"changeme placeholder", "changeme-change-this-to-a-random-secret-min-32-chars", true},
		{"32 chars", strings.Repeat("a", 32), false},
		{"long random", "NzmJmiDrC2Xnm-AUUIv8MwHRK10cDMpY4i-ZkKtOFAXvFJXyNcXQ0zuvUh0G9wbZ", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Config{JWTSecret: tc.secret}
			if err := c.Validate(); (err != nil) != tc.wantErr {
				t.Fatalf("Validate() err=%v, wantErr=%v", err, tc.wantErr)
			}
		})
	}
}

// TestValidateEmptyNamesNewVar checks that a missing secret names JWT_SECRET.
func TestValidateEmptyNamesNewVar(t *testing.T) {
	err := (&Config{}).Validate()
	if err == nil || !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Fatalf("empty secret error = %v, want a JWT_SECRET message", err)
	}
}

// TestCORSAllowOrigins checks splitting and trimming of the configured origins.
func TestCORSAllowOrigins(t *testing.T) {
	c := &Config{CORSAllowOriginsRaw: "http://localhost:8081, https://app.cookmode.dev "}
	got := c.CORSAllowOrigins()
	if len(got) != 2 || got[0] != "http://localhost:8081" || got[1] != "https://app.cookmode.dev" {
		t.Fatalf("unexpected allowlist: %q", got)
	}
	if got := (&Config{}).CORSAllowOrigins(); len(got) != 0 {
		t.Fatalf("empty raw should give empty allowlist, got %q", got)
	}
}
