package stellar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ContractEvent represents a decoded Soroban contract event.
type ContractEvent struct {
	Type     string // posted|claimed|completed|cancelled
	BountyID uint64
	Ledger   uint64
	TxHash   string
	Payload  map[string]any // decoded event data
}

// Client wraps the Stellar RPC connection.
type Client struct {
	rpcURL            string
	contractID        string
	networkPassphrase string
	httpClient        *http.Client
}

// NewClient creates a new Stellar RPC client.
func NewClient(rpcURL, contractID, networkPassphrase string) *Client {
	return &Client{
		rpcURL:            rpcURL,
		contractID:        contractID,
		networkPassphrase: networkPassphrase,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// rpcRequest is a generic JSON-RPC request.
type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

type eventFilter struct {
	Type        string   `json:"type"`
	ContractIDs []string `json:"contractIds"`
}

type getEventsParams struct {
	StartLedger uint64         `json:"startLedger"`
	Filters     []eventFilter  `json:"filters"`
	Pagination  map[string]any `json:"pagination,omitempty"`
}

// GetEvents fetches contract events from the given ledger onwards.
func (c *Client) GetEvents(ctx context.Context, fromLedger uint64) ([]ContractEvent, error) {
	reqBody := rpcRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "getEvents",
		Params: getEventsParams{
			StartLedger: fromLedger,
			Filters: []eventFilter{
				{
					Type:        "contract",
					ContractIDs: []string{c.contractID},
				},
			},
			Pagination: map[string]any{
				"limit": 1000,
			},
		},
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.rpcURL, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rpc request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcRes struct {
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
		Result struct {
			Events []struct {
				Type   string   `json:"type"`
				Ledger string   `json:"ledger"`
				TxHash string   `json:"txHash"`
				Topic  []string `json:"topic"`
				Value  struct {
					Xdr string `json:"xdr"`
				} `json:"value"`
			} `json:"events"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &rpcRes); err != nil {
		return nil, fmt.Errorf("failed to parse rpc response: %w. Body: %s", err, string(body))
	}

	if rpcRes.Error != nil {
		return nil, fmt.Errorf("rpc error: %s", rpcRes.Error.Message)
	}

	var events []ContractEvent
	for _, raw := range rpcRes.Result.Events {
		// Only parse our own contract events
		if raw.Type != "contract" {
			continue
		}

		evt, err := DecodeEvent(raw.Topic, raw.Value.Xdr, raw.Ledger, raw.TxHash)
		if err != nil {
			// Log and skip unknown/invalid events
			fmt.Printf("skipping event: %v\n", err)
			continue
		}
		events = append(events, evt)
	}

	return events, nil
}

// GetLatestLedger returns the current ledger sequence number.
func (c *Client) GetLatestLedger(ctx context.Context) (uint64, error) {
	reqBody := rpcRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "getLatestLedger",
	}

	b, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(ctx, "POST", c.rpcURL, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var rpcRes struct {
		Error  *struct{ Message string } `json:"error"`
		Result struct {
			Sequence uint64 `json:"sequence"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rpcRes); err != nil {
		return 0, err
	}
	if rpcRes.Error != nil {
		return 0, fmt.Errorf("rpc error: %s", rpcRes.Error.Message)
	}

	return rpcRes.Result.Sequence, nil
}
