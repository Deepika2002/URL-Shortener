package write_test

import (
	"context"
	"testing"

	"urlshortener/internal/write"
	"urlshortener/pkg/kafka"
	"urlshortener/pkg/kgsclient"
)

func TestService_Init(t *testing.T) {
	// For TDD, we verify the service can be instantiated with its dependencies.
	// In a complete test suite, these would be mock implementations.
	var repo write.Repository = nil
	var cache write.IdempotencyCache = nil
	var bloom write.BloomFilter = nil
	var kgs kgsclient.Client = nil
	var publisher kafka.Publisher = nil

	svc := write.NewService(repo, cache, bloom, kgs, publisher)
	if svc == nil {
		t.Error("Expected write service to be initialized")
	}
}

// TestService_Interface verifies the expected methods of the core service.
func TestService_Interface(t *testing.T) {
	var svc interface {
		ShortenURL(ctx context.Context, longURL string, customAlias string) (string, error)
	} = write.NewService(nil, nil, nil, nil, nil)

	if svc == nil {
		t.Log("Service is nil, which is fine for this nil-init test")
	}
}
