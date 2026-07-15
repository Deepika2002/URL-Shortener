package read

import (
	"context"
	"errors"
	"fmt"

	"github.com/gocql/gocql"
)

// ErrNotFound is returned when a short ID does not exist in the database.
var ErrNotFound = errors.New("short URL mapping not found")

// Repository defines the data access interface for the Read service.
type Repository interface {
	GetMapping(ctx context.Context, shortID string) (string, error)
}

type repository struct {
	session *gocql.Session
}

// NewRepository creates a new ScyllaDB-backed repository.
func NewRepository(session *gocql.Session) Repository {
	return &repository{
		session: session,
	}
}

// GetMapping retrieves the original long URL for a given short ID.
func (r *repository) GetMapping(ctx context.Context, shortID string) (string, error) {
	var longURL string
	
	// Query ScyllaDB for the mapping
	query := `SELECT long_url FROM url_mappings WHERE short_id = ? LIMIT 1`
	err := r.session.Query(query, shortID).WithContext(ctx).Scan(&longURL)
	if err != nil {
		if err == gocql.ErrNotFound {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("failed to fetch mapping from scylladb: %w", err)
	}

	return longURL, nil
}
