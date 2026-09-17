package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// Event mirrors the events table.
type Event struct {
	ID        uint64         `json:"id"`
	BountyID  uint64         `json:"bounty_id"`
	EventType string         `json:"event_type"`
	Ledger    uint64         `json:"ledger"`
	TxHash    string         `json:"tx_hash"`
	Payload   map[string]any `json:"payload"`
	CreatedAt time.Time      `json:"created_at"`
}

// InsertEvent inserts a new event record.
func (d *DB) InsertEvent(ctx context.Context, tx pgx.Tx, e Event) error {
	query := `
		INSERT INTO events (
			bounty_id, event_type, ledger, tx_hash, payload
		) VALUES (
			$1, $2, $3, $4, $5
		)
	`
	_, err := tx.Exec(ctx, query,
		e.BountyID, e.EventType, e.Ledger, e.TxHash, e.Payload,
	)
	return err
}

// ListEventsByBounty fetches all events for a given bounty.
func (d *DB) ListEventsByBounty(ctx context.Context, bountyID uint64) ([]Event, error) {
	query := `
		SELECT id, bounty_id, event_type, ledger, tx_hash, payload, created_at
		FROM events
		WHERE bounty_id = $1
		ORDER BY ledger ASC
	`
	rows, err := d.pool.Query(ctx, query, bountyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(
			&e.ID, &e.BountyID, &e.EventType, &e.Ledger, &e.TxHash, &e.Payload, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
