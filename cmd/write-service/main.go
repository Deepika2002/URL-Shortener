package main

import (
	"log"
	"net/http"
	"time"

	"urlshortener/internal/write"
	"urlshortener/pkg/cache"
	"urlshortener/pkg/config"
	"urlshortener/pkg/db"
	"urlshortener/pkg/kafka"
	"urlshortener/pkg/kgsclient"
)

func main() {
	log.Println("Starting Write Service (Shortening)...")

	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Redis (for Idempotency Cache and Bloom Filter)
	redisClient := cache.NewRedisClient(cfg)
	defer redisClient.Close()
	log.Println("Connected to Redis.")

	// 3. Initialize ScyllaDB (for persistent URL mappings)
	scyllaSession, err := db.NewScyllaSession(cfg)
	if err != nil {
		log.Fatalf("Fatal error connecting to ScyllaDB: %v", err)
	}
	defer scyllaSession.Close()
	log.Println("Connected to ScyllaDB.")

	// 4. Initialize Kafka Publisher (for analytics events)
	kafkaPublisher, err := kafka.NewPublisher(cfg)
	if err != nil {
		log.Fatalf("Fatal error connecting to Kafka: %v", err)
	}
	defer kafkaPublisher.Close()
	log.Println("Connected to Kafka.")

	// 5. Initialize Key Generation Service (KGS) HTTP Client
	// For local execution, the client derives the KGS URL using the configured KGS port.
	kgsClient := kgsclient.NewClient(cfg)
	log.Println("Initialized KGS Client")

	// 6. Initialize Write Service Dependencies
	repo := write.NewRepository(scyllaSession)
	idempCache := write.NewIdempotencyCache(redisClient)
	bloom := write.NewBloomFilter(redisClient)

	// 7. Initialize Core Service
	svc := write.NewService(repo, idempCache, bloom, kgsClient, kafkaPublisher)

	// 8. Initialize HTTP Handler
	handler := write.NewHandler(svc, cfg.BaseURL)

	// 9. Start HTTP Server
	serverAddr := ":" + cfg.WriteServicePort
	log.Printf("Starting Write Service HTTP server on %s\n", serverAddr)

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Write Service HTTP server crashed: %v", err)
	}
}
