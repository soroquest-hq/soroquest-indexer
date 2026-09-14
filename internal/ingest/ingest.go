package ingest

import (
	"context"
	"log"
	"time"

	"github.com/ayomustap/soroquest-indexer/internal/db"
	"github.com/ayomustap/soroquest-indexer/internal/stellar"
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
		return err
	}

	events, err := idx.rpc.GetEvents(ctx, lastLedger+1)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := idx.processEvent(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// processEvent dispatches a single contract event to the appropriate handler.
// Wrapped in a database transaction: cursor only advances on success.
func (idx *Indexer) processEvent(ctx context.Context, event stellar.ContractEvent) error {
	// TODO: begin db transaction
	// dispatch by event.Type → handlePosted / handleClaimed / handleCompleted / handleCancelled
	// insert into events table
	// update cursor.last_ledger to event.Ledger
	// commit transaction
	panic("not implemented")
}
