package ingest

import (
	"context"

	"github.com/ayomustap/soroquest-indexer/internal/db"
	"github.com/ayomustap/soroquest-indexer/internal/stellar"
)

// handlePosted processes a bounty_posted event.
// Creates a new Bounty row with status=open.
func handlePosted(ctx context.Context, database *db.DB, event stellar.ContractEvent) error {
	// TODO: decode event.Payload → Bounty fields, call db.UpsertBounty
	panic("not implemented")
}

// handleClaimed processes a bounty_claimed event.
// Updates bounty: status=claimed, claimant=<address>
func handleClaimed(ctx context.Context, database *db.DB, event stellar.ContractEvent) error {
	// TODO: decode event.Payload → claimant, call db.UpsertBounty with updated fields
	panic("not implemented")
}

// handleCompleted processes a bounty_completed event.
// Updates bounty: status=completed
func handleCompleted(ctx context.Context, database *db.DB, event stellar.ContractEvent) error {
	// TODO: decode event.Payload, call db.UpsertBounty with status=completed
	panic("not implemented")
}

// handleCancelled processes a bounty_cancelled event.
// Updates bounty: status=cancelled, claimant=null
func handleCancelled(ctx context.Context, database *db.DB, event stellar.ContractEvent) error {
	// TODO: decode event.Payload, call db.UpsertBounty with status=cancelled
	panic("not implemented")
}
