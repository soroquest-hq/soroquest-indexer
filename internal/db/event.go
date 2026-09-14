package db

import (
	"context"
	"encoding/json"
	"time"
)

// Event mirrors the events table.
type Event struct {
	ID        uint64          `json:"id"`
	BountyID  uint64          `json:"bounty_id"`
	EventType string          `json:"event_type"`
	Ledger    uint64          `json:"ledger"`
	TxHash    string          `json:"tx_hash"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

// InsertEvent appends a new event record.
func (d *DB) InsertEvent(ctx context.Context, e *Event) error {
	// TODO: implement INSERT INTO events
	panic("not implemented")
}

// ListEventsByBounty returns all events for a given bounty ID, ordered by ledger.
func (d *DB) ListEventsByBounty(ctx context.Context, bountyID uint64) ([]*Event, error) {
	// TODO: implement SELECT ... WHERE bounty_id = $1 ORDER BY ledger ASC
	panic("not implemented")
}
