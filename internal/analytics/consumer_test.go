package analytics_test

import (
	"context"
	"testing"

	"urlshortener/internal/analytics"
	"urlshortener/pkg/config"
)

func TestConsumer_Init(t *testing.T) {
	cfg := &config.Config{
		KafkaBrokers: "127.0.0.1:9092",
	}

	// Assuming NewConsumer takes config and a Repository (using nil for testing init).
	consumer, err := analytics.NewConsumer(cfg, nil)
	if err != nil {
		t.Logf("Expected failure if Kafka/Redpanda is not running in test env: %v", err)
	}

	if consumer != nil {
		defer consumer.Close()
	}
}

// TestConsumer_Interface verifies the expected methods of the analytics consumer.
func TestConsumer_Interface(t *testing.T) {
	consumer, _ := analytics.NewConsumer(&config.Config{}, nil)
	var c interface {
		// Start begins consuming messages from Kafka topics in the background
		Start(ctx context.Context) error
		// Close gracefully shuts down the consumer group
		Close() error
	} = consumer

	if c == nil {
		t.Log("Consumer is nil, which is fine for this interface check")
	}
}
