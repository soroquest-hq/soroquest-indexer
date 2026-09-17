package api

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/soroquest-hq/soroquest-indexer/internal/db"
)

// Server holds the HTTP server and its dependencies.
type Server struct {
	db   *db.DB
	port string
}

// NewServer creates a new HTTP server.
func NewServer(database *db.DB, port string) *Server {
	return &Server{db: database, port: port}
}

// Run starts the HTTP server and blocks until ctx is cancelled.
func (s *Server) Run(ctx context.Context) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	bountiesHandler := NewBountiesHandler(s.db)
	healthHandler := NewHealthHandler(s.db)

	// Routes
	r.Get("/health", healthHandler.Check)
	r.Route("/api", func(r chi.Router) {
		r.Get("/bounties", bountiesHandler.List)
		r.Get("/bounties/{id}", bountiesHandler.Get)
		r.Get("/bounties/{id}/events", bountiesHandler.Events)
		r.Get("/stats", bountiesHandler.Stats)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.port),
		Handler: r,
	}

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	select {
	case <-ctx.Done():
		return srv.Shutdown(context.Background())
	case err := <-errCh:
		return err
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Vercel domain restriction (using environment variable, default to * for dev)
		origin := os.Getenv("ALLOWED_ORIGIN")
		if origin == "" {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
