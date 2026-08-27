package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cookmode/internal/api"
	"cookmode/internal/config"
	"cookmode/internal/db"
	"cookmode/internal/jobs"
	"cookmode/internal/services"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	searchIndexer := services.NewSearchIndexer(cfg.MeiliHost, cfg.MeiliAPIKey)

	riverClient, err := jobs.NewClient(pool, searchIndexer)
	if err != nil {
		log.Fatalf("river client: %v", err)
	}

	if err := jobs.Migrate(ctx, pool); err != nil {
		log.Fatalf("river migrate: %v", err)
	}

	if err := riverClient.Start(ctx); err != nil {
		log.Fatalf("river start: %v", err)
	}

	router := api.NewRouter(api.Deps{
		Pool:             pool,
		SearchIndexer:    searchIndexer,
		VideoService:     services.NewVideoService(cfg.CFStreamAccountID, cfg.CFStreamAPIToken),
		NutritionService: services.NewNutritionService(cfg.EdamamAppID, cfg.EdamamAppKey, cfg.USDAAPIKey),
		JWTSecret:        cfg.SupabaseJWTSecret,
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	if err := riverClient.Stop(shutdownCtx); err != nil {
		log.Printf("river stop: %v", err)
	}

	cancel()
	log.Println("stopped")
}
