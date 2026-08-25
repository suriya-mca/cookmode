package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string

	SupabaseURL       string
	SupabaseJWTSecret string

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

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on environment")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),

		SupabaseURL:       getEnv("SUPABASE_URL", ""),
		SupabaseJWTSecret: getEnv("SUPABASE_JWT_SECRET", ""),

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

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
