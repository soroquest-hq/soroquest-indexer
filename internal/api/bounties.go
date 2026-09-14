package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleListBounties(w http.ResponseWriter, r *http.Request) {
	// TODO: parse query params: status, limit (default 20), offset (default 0)
	// call db.ListBounties
	// return JSON: { bounties: [], total: 0, limit: 20, offset: 0 }
	writeJSON(w, http.StatusOK, map[string]any{
		"bounties": []any{},
		"total":    0,
		"limit":    20,
		"offset":   0,
	})
}

func (s *Server) handleGetBounty(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid bounty id")
		return
	}
	// TODO: call db.GetBounty(id)
	// return 404 if not found
	_ = id
	writeError(w, http.StatusNotImplemented, "not implemented")
}

func (s *Server) handleGetBountyEvents(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid bounty id")
		return
	}
	// TODO: call db.ListEventsByBounty(id)
	_ = id
	writeError(w, http.StatusNotImplemented, "not implemented")
}

func (s *Server) handleGetStats(w http.ResponseWriter, r *http.Request) {
	// TODO: call db.GetStats()
	writeError(w, http.StatusNotImplemented, "not implemented")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
