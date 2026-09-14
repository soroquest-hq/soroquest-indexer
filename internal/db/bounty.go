package db

import (
	"context"
	"time"
)

// Bounty mirrors the bounties table.
type Bounty struct {
	ID            uint64    `json:"id"`
	Owner         string    `json:"owner"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Amount        string    `json:"amount"`        // raw i128 as decimal string
	Token         string    `json:"token"`
	Status        string    `json:"status"`        // open|claimed|completed|cancelled
	Claimant      *string   `json:"claimant"`
	CreatedAt     uint64    `json:"created_at"`    // ledger sequence
	ClaimDeadline uint64    `json:"claim_deadline"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ListFilter holds query parameters for listing bounties.
type ListFilter struct {
	Status *string
	Limit  int
	Offset int
}

// UpsertBounty inserts or updates a bounty record.
func (d *DB) UpsertBounty(ctx context.Context, b *Bounty) error {
	// TODO: implement INSERT ... ON CONFLICT (id) DO UPDATE
	panic("not implemented")
}

// GetBounty fetches a single bounty by ID.
func (d *DB) GetBounty(ctx context.Context, id uint64) (*Bounty, error) {
	// TODO: implement SELECT ... WHERE id = $1
	panic("not implemented")
}

// ListBounties fetches bounties with optional status filter and pagination.
func (d *DB) ListBounties(ctx context.Context, f ListFilter) ([]*Bounty, int, error) {
	// TODO: implement SELECT with WHERE, LIMIT, OFFSET
	panic("not implemented")
}

// GetStats returns platform-wide aggregate statistics.
func (d *DB) GetStats(ctx context.Context) (*Stats, error) {
	// TODO: implement COUNT(*) GROUP BY status + SUM(amount) WHERE status=completed
	panic("not implemented")
}

// Stats holds aggregate bounty statistics.
type Stats struct {
	Total     int    `json:"total"`
	Open      int    `json:"open"`
	Claimed   int    `json:"claimed"`
	Completed int    `json:"completed"`
	Cancelled int    `json:"cancelled"`
	TotalUSDC string `json:"total_usdc"` // sum of completed bounty amounts
}
