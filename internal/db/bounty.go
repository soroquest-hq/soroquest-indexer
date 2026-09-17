package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Bounty mirrors the bounties table.
type Bounty struct {
	ID            uint64    `json:"id"`
	Owner         string    `json:"owner"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Amount        string    `json:"amount"` // raw i128 as decimal string
	Token         string    `json:"token"`
	Status        string    `json:"status"` // open|claimed|completed|cancelled
	Claimant      *string   `json:"claimant"`
	CreatedAt     uint64    `json:"created_at"` // ledger sequence
	ClaimDeadline uint64    `json:"claim_deadline"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ListFilter holds query parameters for listing bounties.
type ListFilter struct {
	Status *string
	Limit  int
	Offset int
}

// InsertBounty inserts a new bounty record.
func (d *DB) InsertBounty(ctx context.Context, tx pgx.Tx, b Bounty) error {
	query := `
		INSERT INTO bounties (
			id, owner, title, description, amount, token, status, claimant, created_at, claim_deadline
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`
	_, err := tx.Exec(ctx, query,
		b.ID, b.Owner, b.Title, b.Description, b.Amount, b.Token, b.Status, b.Claimant, b.CreatedAt, b.ClaimDeadline,
	)
	return err
}

// UpdateBounty updates an existing bounty record.
func (d *DB) UpdateBounty(ctx context.Context, tx pgx.Tx, b Bounty) error {
	query := `
		UPDATE bounties 
		SET status = $1, claimant = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := tx.Exec(ctx, query, b.Status, b.Claimant, b.ID)
	return err
}

// GetBounty fetches a single bounty by ID.
func (d *DB) GetBounty(ctx context.Context, id uint64) (Bounty, error) {
	query := `
		SELECT id, owner, title, description, amount, token, status, claimant, created_at, claim_deadline, updated_at
		FROM bounties
		WHERE id = $1
	`
	var b Bounty
	err := d.pool.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.Owner, &b.Title, &b.Description, &b.Amount, &b.Token, &b.Status, &b.Claimant, &b.CreatedAt, &b.ClaimDeadline, &b.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Bounty{}, fmt.Errorf("bounty %d not found", id)
	}
	return b, err
}

// ListBounties fetches bounties with optional filters and pagination.
func (d *DB) ListBounties(ctx context.Context, statusFilter *string, ownerFilter *string, claimantFilter *string, limit, offset int) ([]Bounty, int, error) {
	var bounties []Bounty
	var total int

	queryConditions := []string{}
	args := []interface{}{}
	argId := 1

	if statusFilter != nil {
		queryConditions = append(queryConditions, fmt.Sprintf("status = $%d", argId))
		args = append(args, *statusFilter)
		argId++
	}
	if ownerFilter != nil {
		queryConditions = append(queryConditions, fmt.Sprintf("owner = $%d", argId))
		args = append(args, *ownerFilter)
		argId++
	}
	if claimantFilter != nil {
		queryConditions = append(queryConditions, fmt.Sprintf("claimant = $%d", argId))
		args = append(args, *claimantFilter)
		argId++
	}

	whereClause := ""
	if len(queryConditions) > 0 {
		whereClause = " WHERE " + strings.Join(queryConditions, " AND ")
	}

	// 1. Get total count
	countQuery := `SELECT COUNT(*) FROM bounties` + whereClause
	err := d.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 2. Fetch data
	dataQuery := `
		SELECT id, owner, title, description, amount, token, status, claimant, created_at, claim_deadline, updated_at
		FROM bounties
	` + whereClause + `
		ORDER BY id DESC
		LIMIT $` + fmt.Sprint(argId) + ` OFFSET $` + fmt.Sprint(argId+1)

	args = append(args, limit, offset)

	rows, err := d.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var b Bounty
		if err := rows.Scan(
			&b.ID, &b.Owner, &b.Title, &b.Description, &b.Amount, &b.Token, &b.Status, &b.Claimant, &b.CreatedAt, &b.ClaimDeadline, &b.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		bounties = append(bounties, b)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return bounties, total, nil
}

// GetStats returns platform-wide aggregate statistics.
func (d *DB) GetStats(ctx context.Context) (Stats, error) {
	query := `
		SELECT 
			COUNT(*),
			COUNT(*) FILTER (WHERE status = 'open'),
			COUNT(*) FILTER (WHERE status = 'claimed'),
			COUNT(*) FILTER (WHERE status = 'completed'),
			COUNT(*) FILTER (WHERE status = 'cancelled'),
			COALESCE(SUM(amount) FILTER (WHERE status = 'completed'), 0)
		FROM bounties
	`
	var s Stats
	err := d.pool.QueryRow(ctx, query).Scan(
		&s.Total, &s.Open, &s.Claimed, &s.Completed, &s.Cancelled, &s.TotalUSDC,
	)
	return s, err
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
