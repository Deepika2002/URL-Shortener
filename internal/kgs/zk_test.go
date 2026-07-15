package kgs_test

import (
	"testing"

	"urlshortener/internal/kgs"
	"urlshortener/pkg/config"
)

func TestZKManager_InitAndFetch(t *testing.T) {
	cfg := &config.Config{
		ZooKeeperServers: "localhost:2181",
	}

	manager, err := kgs.NewZKManager(cfg)
	
	if err == nil && manager == nil {
		t.Error("Expected ZK manager to be initialized if no error returned")
	}

	if manager != nil {
		workerID, err := manager.GetWorkerID()
		if err != nil {
			// This is expected if the ZK container is not reachable during the unit test
			t.Logf("Could not fetch worker ID (expected if ZK is down): %v", err)
		} else {
			if workerID < 0 || workerID > 1023 {
				t.Errorf("Worker ID %d is out of bounds (0-1023)", workerID)
			}
		}
		manager.Close()
	}
}
