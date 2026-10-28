package db

import (
	"context"
	"testing"
	"time"
)

func TestConnectInvalidURL(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Should fail fast on invalid postgres URL format
	_, err := Connect(ctx, "invalid://database:url")
	if err == nil {
		t.Errorf("expected error connecting to invalid database URL, got nil")
	}
}
