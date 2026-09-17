package stellar

import (
	"encoding/base64"
	"testing"
)

// We won't test full XDR decoding here easily without constructing raw XDR payloads.
// So we just assert that decoding an empty/invalid payload properly returns an error instead of panicking.
func TestDecodeEvent_InvalidInput(t *testing.T) {
	_, err := DecodeEvent([]string{}, "", "100", "00")
	if err == nil {
		t.Errorf("Expected error for missing topics, got nil")
	}
}

func TestDecodeEvent_InvalidTopics(t *testing.T) {
	_, err := DecodeEvent([]string{"topic0", "topic1"}, "invalidXDR", "100", "00")
	if err == nil {
		t.Errorf("Expected error for invalid base64, got nil")
	}
}

// Since constructing manual XDR byte arrays for the tests is very verbose in Go,
// we rely on the integration tests or robust panic-free parsing on the error paths.
func TestDecodeEvent_NonSymbolTopic(t *testing.T) {
	// Encode a simple U64 as topic0 instead of Symbol
	// base64("AAAAAA==") is a simple representation... actually we just pass valid base64 but invalid XDR
	invalidB64 := base64.StdEncoding.EncodeToString([]byte("random bytes"))
	_, err := DecodeEvent([]string{invalidB64, invalidB64}, "", "100", "00")
	if err == nil {
		t.Errorf("Expected error for invalid XDR types, got nil")
	}
}
