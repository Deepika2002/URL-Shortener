package main

import (
	"log"
	"net/http"
	"time"

	"urlshortener/internal/kgs"
	"urlshortener/pkg/config"
)

func main() {
	log.Println("Starting Key Generation Service (KGS)...")

	// Load Configuration
	cfg := config.LoadConfig()

	// Connect to ZooKeeper
	log.Println("Connecting to ZooKeeper...")
	zkManager, err := kgs.NewZKManager(cfg)
	if err != nil {
		log.Fatalf("Fatal error connecting to ZooKeeper: %v", err)
	}
	defer zkManager.Close()

	// Fetch Worker ID from ZooKeeper
	// The underlying library will block/retry during the session timeout if connecting
	workerID, err := zkManager.GetWorkerID()
	if err != nil {
		log.Fatalf("Fatal error fetching Worker ID from ZooKeeper: %v", err)
	}
	log.Printf("Successfully acquired Worker ID: %d\n", workerID)

	// Initialize Snowflake Generator
	sfGen, err := kgs.NewSnowflakeGenerator(workerID)
	if err != nil {
		log.Fatalf("Fatal error initializing Snowflake generator: %v", err)
	}
	log.Println("Snowflake generator initialized successfully.")

	// Initialize HTTP Handler
	handler := kgs.NewHandler(sfGen)

	// Start HTTP Server
	serverAddr := ":" + cfg.KGSPort
	log.Printf("Starting KGS HTTP server on %s\n", serverAddr)
	
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("KGS HTTP server crashed: %v", err)
	}
}
