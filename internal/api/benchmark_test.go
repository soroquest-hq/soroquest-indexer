package api

import (
	"encoding/json"
	"testing"
)

type sampleBounty struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
	Token       string `json:"token"`
	Status      string `json:"status"`
}

func BenchmarkBountyJSONSerialization(b *testing.B) {
	bounty := sampleBounty{
		ID:          1,
		Title:       "Implement Telegram Notifications for Bounty Claims",
		Description: "Build a lightweight webhook listener or bot that sends a real-time message to the community Telegram channel.",
		Amount:      "50000000",
		Token:       "CBIELTK6YBZJU5UP2WWQEUCYKLPU6AUNZ2BQ4WWFEIE3USCIHMXQDAMA",
		Status:      "open",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(bounty)
		if err != nil {
			b.Fatal(err)
		}
	}
}
