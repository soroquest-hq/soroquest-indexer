package ingest

import (
	"context"
	"fmt"
	"log"

	"github.com/soroquest-hq/soroquest-indexer/internal/db"
	"github.com/soroquest-hq/soroquest-indexer/internal/stellar"
)

func processEvent(ctx context.Context, database *db.DB, evt stellar.ContractEvent) error {
	log.Printf("Processing event %s for bounty %d at ledger %d", evt.Type, evt.BountyID, evt.Ledger)

	// We execute event insertion and state update in a single transaction
	tx, err := database.Pool().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Record the event itself
	dbEvt := db.Event{
		BountyID:  evt.BountyID,
		EventType: evt.Type,
		Ledger:    evt.Ledger,
		TxHash:    evt.TxHash,
		Payload:   evt.Payload,
	}
	if err := database.InsertEvent(ctx, tx, dbEvt); err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	// 2. Update the bounty state based on the event
	switch evt.Type {
	case "posted":
		owner, _ := evt.Payload["owner"].(string)
		amount, _ := evt.Payload["amount"].(string)
		token, _ := evt.Payload["token"].(string)
		deadlineF, _ := evt.Payload["claim_deadline"].(float64)
		deadlineU64 := uint64(deadlineF)
		// Usually a "posted" event lacks the title/description because those might be stored elsewhere or we need to fetch them
		// But for our simple model, we assume the initial insert was created by our API, OR we create a stub.
		// Since we don't have title/desc in the event (too large for events), the frontend should have posted this to our API first.
		// So we will just Upsert or Update. Let's try to update, if it doesn't exist, we insert a stub.

		_, err := database.GetBounty(ctx, evt.BountyID)
		if err != nil {
			// Doesn't exist, insert stub
			b := db.Bounty{
				ID:            evt.BountyID,
				Owner:         owner,
				Title:         "Unknown (Posted via Contract)", // Fallback if API didn't see it first
				Description:   "",
				Amount:        amount,
				Token:         token,
				Status:        "open",
				CreatedAt:     evt.Ledger,
				ClaimDeadline: deadlineU64,
			}
			if err := database.InsertBounty(ctx, tx, b); err != nil {
				return err
			}
		} else {
			// It exists (created by our API POST), just update status to open
			b := db.Bounty{
				ID:     evt.BountyID,
				Status: "open",
			}
			if err := database.UpdateBounty(ctx, tx, b); err != nil {
				return err
			}
		}

	case "claimed":
		claimant, _ := evt.Payload["claimant"].(string)
		b := db.Bounty{
			ID:       evt.BountyID,
			Status:   "claimed",
			Claimant: &claimant,
		}
		if err := database.UpdateBounty(ctx, tx, b); err != nil {
			return err
		}

	case "completed":
		b := db.Bounty{
			ID:     evt.BountyID,
			Status: "completed",
		}
		// Notice we don't change the claimant (it's already set)
		// But since UpdateBounty takes the whole struct, we should fetch first to preserve claimant
		existing, err := database.GetBounty(ctx, evt.BountyID)
		if err == nil {
			b.Claimant = existing.Claimant
		}
		if err := database.UpdateBounty(ctx, tx, b); err != nil {
			return err
		}

	case "cancelled":
		b := db.Bounty{
			ID:     evt.BountyID,
			Status: "cancelled",
		}
		existing, err := database.GetBounty(ctx, evt.BountyID)
		if err == nil {
			b.Claimant = existing.Claimant
		}
		if err := database.UpdateBounty(ctx, tx, b); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
