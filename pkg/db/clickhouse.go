package db

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"urlshortener/pkg/config"
)

// NewClickHouseConn initializes and pings a connection to ClickHouse.
func NewClickHouseConn(cfg *config.Config) (driver.Conn, error) {
	addr := cfg.ClickHouseAddr
	if addr == "" {
		addr = "127.0.0.1:19000"
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "password",
		},
		MaxOpenConns: 10,
		MaxIdleConns: 5,
		DialTimeout:  time.Second * 10,
	})
	if err != nil {
		return nil, err
	}

	// Verify the connection is truly alive
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}
