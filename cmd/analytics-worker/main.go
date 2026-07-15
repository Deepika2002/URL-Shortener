package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"urlshortener/internal/analytics"
	"urlshortener/pkg/config"
	"urlshortener/pkg/db"
)

func main() {
	log.Println("Starting Analytics Worker...")

	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize ClickHouse Connection (The Analytical Data Warehouse)
	clickHouseConn, err := db.NewClickHouseConn(cfg)
	if err != nil {
		log.Fatalf("Fatal error connecting to ClickHouse: %v", err)
	}
	defer clickHouseConn.Close()
	log.Println("Connected to ClickHouse.")

	// 3. Initialize ClickHouse Repository
	repo := analytics.NewClickHouseRepository(clickHouseConn)

	// 4. Initialize Kafka/Redpanda Consumer Group
	consumer, err := analytics.NewConsumer(cfg, repo)
	if err != nil {
		log.Fatalf("Fatal error creating Kafka consumer: %v", err)
	}
	// Ensure we politely leave the consumer group on exit
	defer consumer.Close()
	log.Println("Connected to Kafka (Redpanda).")

	// 5. Setup Graceful Shutdown Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for OS interrupt signals (Ctrl+C, Docker stop)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Printf("Received shutdown signal (%v). Stopping consumer gracefully...\n", sig)
		cancel()
	}()

	// 6. Start the blocking consumer loop
	log.Println("Analytics Worker is now actively listening for events...")
	if err := consumer.Start(ctx); err != nil {
		log.Printf("Consumer exited with error: %v\n", err)
	}

	log.Println("Analytics Worker shut down successfully.")
}
