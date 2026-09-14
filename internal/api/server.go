package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ayomustap/soroquest-indexer/internal/db"
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

	// Routes
	r.Get("/health", s.handleHealth)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/bounties", s.handleListBounties)
		r.Get("/bounties/{id}", s.handleGetBounty)
		r.Get("/bounties/{id}/events", s.handleGetBountyEvents)
		r.Get("/stats", s.handleGetStats)
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
		w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: restrict to Vercel domain
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
