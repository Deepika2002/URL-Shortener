package config

import (
	"os"
)

type Config struct {
	// ScyllaDB
	ScyllaHosts    string
	ScyllaPort     string
	ScyllaKeyspace string

	// Redis
	RedisAddr string

	// Redpanda (Kafka)
	KafkaBrokers string

	// ClickHouse
	ClickHouseAddr string

	// ZooKeeper
	ZooKeeperServers string

	// Service Ports
	ReadServicePort  string
	WriteServicePort string
	KGSPort          string

	// Application Config
	BaseURL string
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func LoadConfig() *Config {
	return &Config{
		ScyllaHosts:      getEnv("SCYLLA_HOSTS", "127.0.0.1"),
		ScyllaPort:       getEnv("SCYLLA_PORT", "9042"),
		ScyllaKeyspace:   getEnv("SCYLLA_KEYSPACE", "url_shortener"),
		RedisAddr:        getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		KafkaBrokers:     getEnv("KAFKA_BROKERS", "127.0.0.1:19092"),
		ClickHouseAddr:   getEnv("CLICKHOUSE_ADDR", "127.0.0.1:19000"),
		ZooKeeperServers: getEnv("ZOOKEEPER_SERVERS", "127.0.0.1:2181"),
		ReadServicePort:  getEnv("READ_SERVICE_PORT", "8081"),
		WriteServicePort: getEnv("WRITE_SERVICE_PORT", "8082"),
		KGSPort:          getEnv("KGS_PORT", "8083"),
		BaseURL:          getEnv("BASE_URL", "http://localhost:8081"),
	}
}
