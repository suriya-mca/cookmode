package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	// minJWTSecretLen is the minimum accepted JWT signing secret length.
	// Short secrets make HS256 tokens brute-forceable.
	minJWTSecretLen = 32
	// legacyPlaceholderJWTSecret is the old .env.example value from the
	// Supabase era. It is no longer read, but still rejected so a stale
	// env cannot start the server with a publicly-known signing key.
	legacyPlaceholderJWTSecret = "your-supabase-jwt-secret"
)

// Validate fails fast on config that would silently weaken auth. In
// particular an empty, short, or well-known JWT secret must never reach
// the token verifier.
func (c *Config) Validate() error {
	if c.JWTSecret == "" {
		return errors.New("JWT_SECRET is not set — generate one with: openssl rand -base64 48")
	}
	if c.JWTSecret == legacyPlaceholderJWTSecret {
		return errors.New("JWT_SECRET is still the old placeholder value — set a random secret")
	}
	if strings.HasPrefix(c.JWTSecret, "changeme") {
		return errors.New("JWT_SECRET is still the .env.example changeme value — set a random secret")
	}
	if len(c.JWTSecret) < minJWTSecretLen {
		return fmt.Errorf("JWT_SECRET must be at least %d characters (got %d)", minJWTSecretLen, len(c.JWTSecret))
	}
	return nil
}

type Config struct {
	Port        string
	DatabaseURL string

	// JWTSecret signs and verifies locally-issued HS256 auth tokens.
	JWTSecret string

	// CORSAllowOriginsRaw is a comma-separated list of web origins allowed
	// to call the API ("*" = anywhere; dev only).
	CORSAllowOriginsRaw string

	MeiliHost   string
	MeiliAPIKey string

	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2Bucket          string

	CFStreamAccountID string
	CFStreamAPIToken  string

	EdamamAppID  string
	EdamamAppKey string
	USDAAPIKey   string
}

// Load reads application settings from the environment, with local defaults.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on environment")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),

		JWTSecret: getEnv("JWT_SECRET", ""),

		CORSAllowOriginsRaw: getEnv("CORS_ALLOW_ORIGINS", "http://localhost:8081"),

		MeiliHost:   getEnv("MEILI_HOST", "http://localhost:7700"),
		MeiliAPIKey: getEnv("MEILI_API_KEY", ""),

		R2AccountID:       getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2Bucket:          getEnv("R2_BUCKET", ""),

		CFStreamAccountID: getEnv("CF_STREAM_ACCOUNT_ID", ""),
		CFStreamAPIToken:  getEnv("CF_STREAM_API_TOKEN", ""),

		EdamamAppID:  getEnv("EDAMAM_APP_ID", ""),
		EdamamAppKey: getEnv("EDAMAM_APP_KEY", ""),
		USDAAPIKey:   getEnv("USDA_API_KEY", ""),
	}
}

// CORSAllowOrigins parses CORSAllowOriginsRaw into a trimmed allowlist.
func (c *Config) CORSAllowOrigins() []string {
	var out []string
	for _, o := range strings.Split(c.CORSAllowOriginsRaw, ",") {
		if o = strings.TrimSpace(o); o != "" {
			out = append(out, o)
		}
	}
	return out
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
