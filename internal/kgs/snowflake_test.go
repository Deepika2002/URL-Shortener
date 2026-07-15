package kgs_test

import (
	"testing"

	"urlshortener/internal/kgs"
)

func TestSnowflakeGenerator_UniqueIDs(t *testing.T) {
	// Initialize a generator with worker ID 1
	gen, err := kgs.NewSnowflakeGenerator(1)
	if err != nil {
		t.Fatalf("Failed to initialize Snowflake generator: %v", err)
	}

	id1, err := gen.NextID()
	if err != nil {
		t.Fatalf("Failed to generate first ID: %v", err)
	}

	id2, err := gen.NextID()
	if err != nil {
		t.Fatalf("Failed to generate second ID: %v", err)
	}

	if id1 == id2 {
		t.Fatalf("Snowflake generator produced duplicate IDs: %d and %d", id1, id2)
	}
	
	if id1 == 0 || id2 == 0 {
		t.Fatalf("Snowflake generator produced invalid 0 ID")
	}
}

func TestSnowflakeGenerator_InvalidWorkerID(t *testing.T) {
	// Snowflake typically limits worker ID between 0 and 1023
	_, err := kgs.NewSnowflakeGenerator(2000)
	if err == nil {
		t.Error("Expected error when initializing Snowflake generator with invalid worker ID, got nil")
	}
}
