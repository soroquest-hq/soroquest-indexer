package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/soroquest-hq/soroquest-indexer/internal/db"
	"github.com/soroquest-hq/soroquest-indexer/internal/ingest"
)

type HealthHandler struct {
	db        *db.DB
	startTime time.Time
}

func NewHealthHandler(db *db.DB) *HealthHandler {
	return &HealthHandler{
		db:        db,
		startTime: time.Now(),
	}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	// Ping DB
	dbErr := h.db.Pool().Ping(r.Context())
	dbStatus := "ok"
	if dbErr != nil {
		dbStatus = "error: " + dbErr.Error()
	}

	// Get latest indexed ledger
	lastLedger, err := ingest.GetLastLedger(r.Context(), h.db)
	if err != nil {
		lastLedger = 0
	}

	res := map[string]any{
		"status":      "ok",
		"uptime":      time.Since(h.startTime).String(),
		"db":          dbStatus,
		"last_ledger": lastLedger,
	}

	if dbErr != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
