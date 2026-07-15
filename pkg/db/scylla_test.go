package db_test

import (
	"testing"

	"urlshortener/pkg/config"
	"urlshortener/pkg/db"
)

func TestNewScyllaSession_InvalidHost(t *testing.T) {
	cfg := &config.Config{
		ScyllaHosts:    "invalid_host:9999",
		ScyllaKeyspace: "url_shortener",
	}

	session, err := db.NewScyllaSession(cfg)
	if err == nil {
		t.Error("Expected error when connecting to invalid ScyllaDB host, got nil")
	}
	if session != nil {
		session.Close()
	}
}
