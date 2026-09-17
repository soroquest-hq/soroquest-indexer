package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/soroquest-hq/soroquest-indexer/internal/db"
)

type BountiesHandler struct {
	db *db.DB
}

func NewBountiesHandler(db *db.DB) *BountiesHandler {
	return &BountiesHandler{db: db}
}

// List handles GET /api/bounties
func (h *BountiesHandler) List(w http.ResponseWriter, r *http.Request) {
	statusQuery := r.URL.Query().Get("status")
	var status *string
	if statusQuery != "" {
		status = &statusQuery
	}

	ownerQuery := r.URL.Query().Get("owner")
	var owner *string
	if ownerQuery != "" {
		owner = &ownerQuery
	}

	claimantQuery := r.URL.Query().Get("claimant")
	var claimant *string
	if claimantQuery != "" {
		claimant = &claimantQuery
	}

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	bounties, total, err := h.db.ListBounties(r.Context(), status, owner, claimant, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if bounties == nil {
		bounties = []db.Bounty{} // Return empty array instead of null
	}

	res := map[string]any{
		"data": bounties,
		"meta": map[string]any{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Get handles GET /api/bounties/{id}
func (h *BountiesHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	bounty, err := h.db.GetBounty(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bounty)
}

// Events handles GET /api/bounties/{id}/events
func (h *BountiesHandler) Events(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	events, err := h.db.ListEventsByBounty(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if events == nil {
		events = []db.Event{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// Stats handles GET /api/stats
func (h *BountiesHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.db.GetStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
