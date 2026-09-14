package ingest

import (
	"context"
	"fmt"

	"github.com/ayomustap/soroquest-indexer/internal/db"
)

// GetLastLedger reads the last processed ledger from the cursor table.
func GetLastLedger(ctx context.Context, database *db.DB) (uint64, error) {
	// TODO: SELECT last_ledger FROM cursor WHERE id = 1
	panic("not implemented")
}

// UpdateLastLedger sets the last processed ledger in the cursor table.
func UpdateLastLedger(ctx context.Context, database *db.DB, ledger uint64) error {
	// TODO: UPDATE cursor SET last_ledger = $1, updated_at = NOW() WHERE id = 1
	_ = fmt.Sprintf("updating cursor to ledger %d", ledger)
	panic("not implemented")
}
