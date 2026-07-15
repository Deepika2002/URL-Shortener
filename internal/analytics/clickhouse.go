package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type clickHouseRepo struct {
	conn driver.Conn
}

// NewClickHouseRepository creates a new ClickHouse-backed analytics repository.
func NewClickHouseRepository(conn driver.Conn) Repository {
	return &clickHouseRepo{
		conn: conn,
	}
}

// RecordEvent parses the Kafka JSON payload and inserts the event into ClickHouse.
func (r *clickHouseRepo) RecordEvent(ctx context.Context, eventType string, payload []byte) error {
	var data map[string]string
	if err := json.Unmarshal(payload, &data); err != nil {
		return fmt.Errorf("failed to parse event payload: %w", err)
	}

	shortID := data["short_id"]
	longURL := data["long_url"]
	
	// Dynamically determine the timestamp based on the event type payload
	var eventTime time.Time
	var err error
	if tStr, ok := data["created_at"]; ok {
		eventTime, err = time.Parse(time.RFC3339Nano, tStr)
	} else if tStr, ok := data["accessed_at"]; ok {
		eventTime, err = time.Parse(time.RFC3339Nano, tStr)
	}
	
	// Fallback to the exact current time if parsing fails or data is malformed
	if err != nil || eventTime.IsZero() {
		eventTime = time.Now().UTC()
	}

	// In a true high-throughput production environment, we would queue these events in memory
	// and flush them to ClickHouse in massive batches (e.g., 100,000 rows at once).
	// For this worker architecture, we use the driver's batch mechanism (batch size of 1 for now).
	batch, err := r.conn.PrepareBatch(ctx, "INSERT INTO analytics_events (short_id, long_url, event_type, timestamp)")
	if err != nil {
		return fmt.Errorf("failed to prepare clickhouse batch: %w", err)
	}

	if err := batch.Append(shortID, longURL, eventType, eventTime); err != nil {
		return fmt.Errorf("failed to append row to clickhouse batch: %w", err)
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send clickhouse batch: %w", err)
	}

	log.Printf("[Analytics] Recorded '%s' event for ID: %s", eventType, shortID)
	return nil
}
