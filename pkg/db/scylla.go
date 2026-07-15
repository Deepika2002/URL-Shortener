package db

import (
	"strconv"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"urlshortener/pkg/config"
)

// NewScyllaSession creates and returns a new gocql session connected to ScyllaDB.
func NewScyllaSession(cfg *config.Config) (*gocql.Session, error) {
	hosts := strings.Split(cfg.ScyllaHosts, ",")
	cluster := gocql.NewCluster(hosts...)
	
	if cfg.ScyllaPort != "" {
		if port, err := strconv.Atoi(cfg.ScyllaPort); err == nil {
			cluster.Port = port
		}
	}
	
	cluster.Keyspace = cfg.ScyllaKeyspace
	cluster.Consistency = gocql.LocalQuorum
	cluster.Timeout = 5 * time.Second
	cluster.ConnectTimeout = 5 * time.Second

	return cluster.CreateSession()
}
