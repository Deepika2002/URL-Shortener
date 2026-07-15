package kgs

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
	"urlshortener/pkg/config"
)

// ZKManager defines the interface for interacting with ZooKeeper to coordinate distributed workers.
type ZKManager interface {
	GetWorkerID() (int64, error)
	Close()
}

type zkManager struct {
	conn *zk.Conn
}

// NewZKManager creates a new ZooKeeper manager.
func NewZKManager(cfg *config.Config) (ZKManager, error) {
	servers := strings.Split(cfg.ZooKeeperServers, ",")
	conn, _, err := zk.Connect(servers, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ZooKeeper: %w", err)
	}

	return &zkManager{
		conn: conn,
	}, nil
}

// GetWorkerID creates an ephemeral sequential node in ZooKeeper and derives a unique worker ID (0-1023) from it.
func (z *zkManager) GetWorkerID() (int64, error) {
	// Ensure the parent nodes exist
	exists, _, err := z.conn.Exists("/kgs")
	if err == nil && !exists {
		_, _ = z.conn.Create("/kgs", []byte{}, 0, zk.WorldACL(zk.PermAll))
	}
	workerPath := "/kgs/workers"
	exists, _, err = z.conn.Exists(workerPath)
	if err == nil && !exists {
		_, _ = z.conn.Create(workerPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
	}

	// Create an ephemeral sequential znode
	nodePath := workerPath + "/worker-"
	path, err := z.conn.Create(nodePath, []byte{}, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		return 0, fmt.Errorf("failed to create sequential znode: %w", err)
	}

	// path looks like /kgs/workers/worker-0000000001
	parts := strings.Split(path, "-")
	seqStr := parts[len(parts)-1]
	seq, err := strconv.ParseInt(seqStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse sequence number from path %s: %w", path, err)
	}

	// Snowflake requires the worker ID to be between 0 and 1023
	workerID := seq % 1024
	return workerID, nil
}

// Close closes the ZooKeeper connection.
func (z *zkManager) Close() {
	if z.conn != nil {
		z.conn.Close()
	}
}
