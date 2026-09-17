package ingest

import (
	"context"
	"log"
	"time"

	"github.com/soroquest-hq/soroquest-indexer/internal/db"
	"github.com/soroquest-hq/soroquest-indexer/internal/stellar"
)

const pollInterval = 5 * time.Second

// Indexer drives the ingestion loop.
type Indexer struct {
	db  *db.DB
	rpc *stellar.Client
}

// New creates a new Indexer.
func New(database *db.DB, rpc *stellar.Client) *Indexer {
	return &Indexer{db: database, rpc: rpc}
}

// Run starts the ingestion loop. Blocks until ctx is cancelled.
func (idx *Indexer) Run(ctx context.Context) error {
	for {
		if err := idx.tick(ctx); err != nil {
			log.Printf("ingest tick error: %v — retrying in %v", err, pollInterval)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(pollInterval):
		}
	}
}

// tick runs one ingestion cycle: fetch new events, process each one.
func (idx *Indexer) tick(ctx context.Context) error {
	lastLedger, err := GetLastLedger(ctx, idx.db)
	if err != nil {
		// Start at 0 if no cursor exists
		lastLedger = 0
	}

	events, err := idx.rpc.GetEvents(ctx, lastLedger+1)
	if err != nil {
		return err
	}

	var latestLedger uint64 = lastLedger

	for _, event := range events {
		if err := processEvent(ctx, idx.db, event); err != nil {
			return err
		}
		if event.Ledger > latestLedger {
			latestLedger = event.Ledger
		}
	}

	if latestLedger > lastLedger {
		// Update the global cursor
		tx, err := idx.db.Pool().Begin(ctx)
		if err == nil {
			if err := UpdateLastLedgerTx(ctx, tx, latestLedger); err == nil {
				tx.Commit(ctx)
			} else {
				tx.Rollback(ctx)
			}
		}
	}

	// If there were no events, maybe we can fast forward cursor to the network's latest ledger
	if len(events) == 0 {
		netLatest, err := idx.rpc.GetLatestLedger(ctx)
		if err == nil && netLatest > latestLedger {
			tx, err := idx.db.Pool().Begin(ctx)
			if err == nil {
				if err := UpdateLastLedgerTx(ctx, tx, netLatest); err == nil {
					tx.Commit(ctx)
				} else {
					tx.Rollback(ctx)
				}
			}
		}
	}

	return nil
}
