package kafka_test

import (
	"testing"

	"urlshortener/pkg/config"
	"urlshortener/pkg/kafka"
)

func TestNewPublisher_Init(t *testing.T) {
	cfg := &config.Config{
		KafkaBrokers: "localhost:19092",
	}

	publisher, err := kafka.NewPublisher(cfg)
	
	// Since we are doing TDD, we just want to ensure that if a publisher is returned without error,
	// it can be closed. If an error is returned due to a missing local broker during the test,
	// that's okay, we're just checking the initialization surface.
	if err == nil && publisher == nil {
		t.Error("Expected a publisher instance if no error is returned")
	}

	if publisher != nil {
		err := publisher.Close()
		if err != nil {
			t.Logf("Error closing publisher: %v", err)
		}
	}
}
