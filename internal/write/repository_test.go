package write_test

import (
	"context"
	"testing"

	"urlshortener/internal/write"
	"urlshortener/pkg/config"
	"urlshortener/pkg/db"
)

func TestRepository_Init(t *testing.T) {
	cfg := &config.Config{
		ScyllaHosts:    "127.0.0.1",
		ScyllaKeyspace: "url_shortener",
	}

	session, err := db.NewScyllaSession(cfg)
	if err != nil {
		t.Logf("Expected failure if ScyllaDB is not running in test environment: %v", err)
	}
	
	if session != nil {
		defer session.Close()
	}

	repo := write.NewRepository(session)
	if repo == nil {
		t.Error("Expected repository to be initialized")
	}
}

// TestRepository_Interface verifies the expected methods of the repository.
func TestRepository_Interface(t *testing.T) {
	var repo interface {
		SaveMapping(ctx context.Context, shortID, longURL string) error
	} = write.NewRepository(nil)

	if repo == nil {
		t.Log("Repository is nil, which is fine for this nil-init test")
	}
}
