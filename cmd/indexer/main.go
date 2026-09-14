package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/ayomustap/soroquest-indexer/internal/api"
	"github.com/ayomustap/soroquest-indexer/internal/db"
	"github.com/ayomustap/soroquest-indexer/internal/ingest"
	"github.com/ayomustap/soroquest-indexer/internal/stellar"
)

func main() {
	_ = godotenv.Load() // load .env if present, ignore error in prod

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Database
	database, err := db.New(ctx, mustEnv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer database.Close()

	// Stellar RPC client
	stellarClient := stellar.NewClient(
		mustEnv("STELLAR_RPC_URL"),
		mustEnv("CONTRACT_ID"),
		mustEnv("NETWORK_PASSPHRASE"),
	)

	// Ingestion loop (runs in background)
	indexer := ingest.New(database, stellarClient)
	go func() {
		if err := indexer.Run(ctx); err != nil {
			log.Printf("ingest: %v", err)
		}
	}()

	// HTTP API
	port := getEnv("PORT", "8080")
	server := api.NewServer(database, port)
	if err := server.Run(ctx); err != nil {
		log.Fatalf("api: %v", err)
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %q is not set", key)
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
