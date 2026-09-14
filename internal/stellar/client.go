package stellar

import (
	"context"
)

// ContractEvent represents a decoded Soroban contract event.
type ContractEvent struct {
	Type    string         // posted|claimed|completed|cancelled
	BountyID uint64
	Ledger  uint64
	TxHash  string
	Payload map[string]any // decoded event data
}

// Client wraps the Stellar RPC connection.
type Client struct {
	rpcURL            string
	contractID        string
	networkPassphrase string
}

// NewClient creates a new Stellar RPC client.
func NewClient(rpcURL, contractID, networkPassphrase string) *Client {
	return &Client{
		rpcURL:            rpcURL,
		contractID:        contractID,
		networkPassphrase: networkPassphrase,
	}
}

// GetEvents fetches contract events from the given ledger onwards.
// Returns decoded ContractEvent structs ready for the ingestion handlers.
func (c *Client) GetEvents(ctx context.Context, fromLedger uint64) ([]ContractEvent, error) {
	// TODO: call Stellar RPC getEvents with contract filter
	// decode XDR → ContractEvent structs
	// return events ordered by ledger ASC
	panic("not implemented")
}

// GetLatestLedger returns the current ledger sequence number.
func (c *Client) GetLatestLedger(ctx context.Context) (uint64, error) {
	// TODO: call Stellar RPC getLatestLedger
	panic("not implemented")
}
