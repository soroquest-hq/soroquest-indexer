package api

import (
	"net/http"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	// TODO: get last_ledger from cursor table via db
	// get current ledger from Stellar RPC
	// compute lag = current - last
	writeJSON(w, http.StatusOK, map[string]any{
		"status":      "ok",
		"last_ledger": 0,
		"lag":         0,
	})
}
