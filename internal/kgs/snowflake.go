package kgs

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

// SnowflakeGenerator generates distributed unique IDs.
type SnowflakeGenerator struct {
	node *snowflake.Node
}

// NewSnowflakeGenerator creates a new generator for a specific worker/node ID.
func NewSnowflakeGenerator(workerID int64) (*SnowflakeGenerator, error) {
	// Snowflake limits node ID bits. bwmarrin/snowflake defaults to 10 bits (0-1023).
	if workerID < 0 || workerID > 1023 {
		return nil, fmt.Errorf("worker ID %d is out of range (0-1023)", workerID)
	}

	node, err := snowflake.NewNode(workerID)
	if err != nil {
		return nil, fmt.Errorf("failed to create snowflake node: %w", err)
	}

	return &SnowflakeGenerator{
		node: node,
	}, nil
}

// NextID generates and returns a new unique ID in thread-safe manner.
func (s *SnowflakeGenerator) NextID() (int64, error) {
	id := s.node.Generate()
	return id.Int64(), nil
}
