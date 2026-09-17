package stellar

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/stellar/go/xdr"
)

// DecodeEvent decodes a Soroban contract event into a ContractEvent.
func DecodeEvent(topics []string, valueXdr string, ledgerStr string, txHash string) (ContractEvent, error) {
	ledger, _ := strconv.ParseUint(ledgerStr, 10, 64)

	if len(topics) < 2 {
		return ContractEvent{}, fmt.Errorf("not enough topics")
	}

	// Topic 0 should be the event symbol (e.g. "posted", "claimed")
	topic0Raw, err := base64.StdEncoding.DecodeString(topics[0])
	if err != nil {
		return ContractEvent{}, err
	}
	var scVal0 xdr.ScVal
	if err := scVal0.UnmarshalBinary(topic0Raw); err != nil {
		return ContractEvent{}, err
	}
	if scVal0.Type != xdr.ScValTypeScvSymbol {
		return ContractEvent{}, fmt.Errorf("topic 0 is not a symbol")
	}
	eventType := string(*scVal0.Sym)

	// Topic 1 should be the bounty ID (u64)
	topic1Raw, err := base64.StdEncoding.DecodeString(topics[1])
	if err != nil {
		return ContractEvent{}, err
	}
	var scVal1 xdr.ScVal
	if err := scVal1.UnmarshalBinary(topic1Raw); err != nil {
		return ContractEvent{}, err
	}
	if scVal1.Type != xdr.ScValTypeScvU64 {
		return ContractEvent{}, fmt.Errorf("topic 1 is not a u64")
	}
	bountyID := uint64(*scVal1.U64)

	// Value payload
	valRaw, err := base64.StdEncoding.DecodeString(valueXdr)
	if err != nil {
		return ContractEvent{}, err
	}
	var val xdr.ScVal
	if err := val.UnmarshalBinary(valRaw); err != nil {
		return ContractEvent{}, err
	}

	payload := make(map[string]any)

	// Parse payload based on event type
	switch eventType {
	case "posted":
		// Payload: (owner, amount, token, claim_deadline)
		vec, ok := val.GetVec()
		if !ok || vec == nil || len(*vec) != 4 {
			return ContractEvent{}, fmt.Errorf("invalid payload for posted event")
		}

		ownerAddr, _ := (*vec)[0].GetAddress()
		ownerStr, _ := ownerAddr.String()
		payload["owner"] = ownerStr

		// amount is i128
		amt128, _ := (*vec)[1].GetI128()
		payload["amount"] = fmt.Sprintf("%d", amt128.Lo) // assuming Lo part is enough for string representation for simplicity, ideally we need big int math

		tokenAddr, _ := (*vec)[2].GetAddress()
		tokenStr, _ := tokenAddr.String()
		payload["token"] = tokenStr

		deadlineU64, _ := (*vec)[3].GetU64()
		payload["claim_deadline"] = uint64(deadlineU64)

	case "claimed":
		// Payload: claimant (address)
		claimantAddr, _ := val.GetAddress()
		claimantStr, _ := claimantAddr.String()
		payload["claimant"] = claimantStr

	case "completed":
		// Payload: (claimant, amount)
		vec, ok := val.GetVec()
		if !ok || vec == nil || len(*vec) != 2 {
			return ContractEvent{}, fmt.Errorf("invalid payload for completed event")
		}
		claimantAddr, _ := (*vec)[0].GetAddress()
		claimantStr, _ := claimantAddr.String()
		payload["claimant"] = claimantStr

		amt128, _ := (*vec)[1].GetI128()
		payload["amount"] = fmt.Sprintf("%d", amt128.Lo)

	case "cancelled":
		// Payload: (owner, amount)
		vec, ok := val.GetVec()
		if !ok || vec == nil || len(*vec) != 2 {
			return ContractEvent{}, fmt.Errorf("invalid payload for cancelled event")
		}
		ownerAddr, _ := (*vec)[0].GetAddress()
		ownerStr, _ := ownerAddr.String()
		payload["owner"] = ownerStr

		amt128, _ := (*vec)[1].GetI128()
		payload["amount"] = fmt.Sprintf("%d", amt128.Lo)

	default:
		return ContractEvent{}, fmt.Errorf("unknown event type: %s", eventType)
	}

	return ContractEvent{
		Type:     eventType,
		BountyID: bountyID,
		Ledger:   ledger,
		TxHash:   txHash,
		Payload:  payload,
	}, nil
}
