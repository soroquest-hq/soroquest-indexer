package ingest

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/soroquest-hq/soroquest-indexer/internal/db"
)

// GetLastLedger reads the last processed ledger from the cursor table.
func GetLastLedger(ctx context.Context, database *db.DB) (uint64, error) {
	var lastLedger uint64
	query := `SELECT last_ledger FROM cursor WHERE id = 1`
	err := database.Pool().QueryRow(ctx, query).Scan(&lastLedger)
	return lastLedger, err
}

// UpdateLastLedgerTx sets the last processed ledger in the cursor table within a transaction.
func UpdateLastLedgerTx(ctx context.Context, tx pgx.Tx, ledger uint64) error {
	query := `UPDATE cursor SET last_ledger = $1, updated_at = NOW() WHERE id = 1`
	_, err := tx.Exec(ctx, query, ledger)
	return err
}
