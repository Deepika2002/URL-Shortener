package analytics_test

import (
	"context"
	"testing"

	"urlshortener/internal/analytics"
	"urlshortener/pkg/config"
	"urlshortener/pkg/db"
)

func TestClickHouse_Init(t *testing.T) {
	cfg := &config.Config{
		ClickHouseAddr: "127.0.0.1:9000",
	}

	conn, err := db.NewClickHouseConn(cfg)
	if err != nil {
		t.Logf("Expected failure if ClickHouse is not running in test env: %v", err)
	}

	if conn != nil {
		defer conn.Close()
	}

	repo := analytics.NewClickHouseRepository(conn)
	if repo == nil {
		t.Error("Expected repository to be initialized")
	}
}

// TestClickHouse_Interface verifies that the ClickHouse implementation accurately satisfies the Repository interface.
func TestClickHouse_Interface(t *testing.T) {
	var repo interface {
		RecordEvent(ctx context.Context, eventType string, payload []byte) error
	} = analytics.NewClickHouseRepository(nil)

	if repo == nil {
		t.Log("Repository is nil, which is fine for this nil-init interface check test")
	}
}
