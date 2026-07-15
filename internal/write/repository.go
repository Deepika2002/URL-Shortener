package write

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

// Repository defines the data access interface for the write service.
type Repository interface {
	SaveMapping(ctx context.Context, shortID, longURL string) error
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

// SaveMapping stores the generated short URL mapping into ScyllaDB.
func (r *repository) SaveMapping(ctx context.Context, shortID, longURL string) error {
	// The keyspace "url_shortener" is implicitly bound to the session during initialization.
	query := `INSERT INTO url_mappings (short_id, long_url, created_at) VALUES (?, ?, ?)`
	
	// Execute the query using the context for timeouts/cancellation
	err := r.session.Query(query, shortID, longURL, time.Now().UTC()).WithContext(ctx).Exec()
	if err != nil {
		return fmt.Errorf("failed to save URL mapping to ScyllaDB: %w", err)
	}
	
	return nil
}
