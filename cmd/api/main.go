package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/soroquest-hq/soroquest-indexer/internal/api"
	"github.com/soroquest-hq/soroquest-indexer/internal/db"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Received termination signal, shutting down API...")
		cancel()
	}()

	log.Println("Connecting to database...")
	database, err := db.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer database.Close()

	server := api.NewServer(database, port)

	log.Printf("Starting API server on :%s", port)
	if err := server.Run(ctx); err != nil {
		log.Printf("Server stopped: %v", err)
	}
}
