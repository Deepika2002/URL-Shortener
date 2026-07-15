package main

import (
	"log"
	"net/http"
	"time"

	"urlshortener/internal/read"
	"urlshortener/pkg/cache"
	"urlshortener/pkg/config"
	"urlshortener/pkg/db"
	"urlshortener/pkg/kafka"
)

func main() {
	log.Println("Starting Read Service (Redirection)...")

	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Redis (for Fast Read Cache and Distributed Mutex)
	redisClient := cache.NewRedisClient(cfg)
	defer redisClient.Close()
	log.Println("Connected to Redis.")

	// 3. Initialize ScyllaDB (for persistent mapping retrieval)
	scyllaSession, err := db.NewScyllaSession(cfg)
	if err != nil {
		log.Fatalf("Fatal error connecting to ScyllaDB: %v", err)
	}
	defer scyllaSession.Close()
	log.Println("Connected to ScyllaDB.")

	// 4. Initialize Kafka Publisher (for redirect analytics events)
	kafkaPublisher, err := kafka.NewPublisher(cfg)
	if err != nil {
		log.Fatalf("Fatal error connecting to Kafka: %v", err)
	}
	defer kafkaPublisher.Close()
	log.Println("Connected to Kafka.")

	// 5. Initialize Read Service Dependencies
	repo := read.NewRepository(scyllaSession)
	cacheWrapper := read.NewCache(redisClient)
	mutex := read.NewMutex(redisClient)

	// 6. Initialize Core Redirection Service
	svc := read.NewService(repo, cacheWrapper, mutex, kafkaPublisher)

	// 7. Initialize HTTP Handler
	handler := read.NewHandler(svc)

	// 8. Start HTTP Server
	serverAddr := ":" + cfg.ReadServicePort
	log.Printf("Starting Read Service HTTP server on %s\n", serverAddr)

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Read Service HTTP server crashed: %v", err)
	}
}
