package read_test

import (
	"context"
	"testing"
	"time"

	"urlshortener/internal/read"
)

// Mocks to verify the orchestration logic of the core Read Service

type mockReadRepo struct{}

func (m *mockReadRepo) GetMapping(ctx context.Context, shortID string) (string, error) {
	if shortID == "exists" {
		return "https://example.com/long", nil
	}
	return "", read.ErrNotFound
}

type mockCache struct{}

func (m *mockCache) Get(ctx context.Context, shortID string) (string, error) { return "", nil }
func (m *mockCache) Set(ctx context.Context, shortID string, longURL string, expiration time.Duration) error {
	return nil
}

type mockMutex struct{}

func (m *mockMutex) AcquireLock(ctx context.Context, shortID string, expiration time.Duration) (bool, error) {
	return true, nil
}
func (m *mockMutex) ReleaseLock(ctx context.Context, shortID string) error { return nil }

type mockPublisher struct{}

func (m *mockPublisher) Publish(topic string, key string, message []byte) error { return nil }
func (m *mockPublisher) Close() error { return nil }

func TestService_Init(t *testing.T) {
	repo := &mockReadRepo{}
	cache := &mockCache{}
	mutex := &mockMutex{}
	publisher := &mockPublisher{}

	svc := read.NewService(repo, cache, mutex, publisher)
	if svc == nil {
		t.Error("Expected Read Service to be initialized")
	}
}

// TestService_Interface verifies the expected methods of the read service.
func TestService_Interface(t *testing.T) {
	var svc interface {
		GetLongURL(ctx context.Context, shortID string) (string, error)
	} = read.NewService(nil, nil, nil, nil)

	if svc == nil {
		t.Log("Service is nil, which is fine for this nil-init test")
	}
}
