# Tech Stack & Dependency Lock

This document strictly defines the technologies, versions, and programming languages required for the URL Shortener architecture. All code and configurations generated MUST comply with these constraints.

## 1. Infrastructure & Coordination
*   **Container Orchestration:** Docker Compose (for local development and AI validation)
*   **API Gateway & Load Balancer:** NGINX (latest stable)
*   **Coordination Service:** ZooKeeper / etcd (used exclusively for KGS Worker ID assignment)

## 2. Databases & Message Brokers
*   **Primary Database:** ScyllaDB 2026.1.x (configured for NoSQL TTL features)
*   **Cache:** Redis (latest stable) (configured with `volatile-lru` eviction policy)
*   **Message Broker:** Redpanda v26.1.x (configured for asynchronous analytics ingestion)
*   **Analytics Database:** ClickHouse 26.3 LTS (configured with `ReplacingMergeTree` engines for deduplication)

## 3. Microservice Languages & Frameworks
*   **The Read Service (Redirection):**
    *   **Language:** Go (Golang)
    *   **Reasoning:** Essential for achieving <10ms latency requirements, high concurrency handling, and efficient connection pooling with Redis and ScyllaDB.
*   **The Write Service (Shortening):**
    *   **Language:** Go (Golang)
    *   **Reasoning:** Fast API response times and robust Bloom filter implementation.
*   **Key Generation Service (KGS):**
    *   **Language:** Go (Golang)
    *   **Reasoning:** Must handle extremely high throughput and atomic counter operations flawlessly. A unified Go ecosystem simplifies CI/CD and library sharing.
*   **Analytics Consumer Workers:**
    *   **Language:** Go (Golang)
    *   **Reasoning:** Go is preferred for high-throughput batching into ClickHouse and maintaining a unified ecosystem across all microservices.

## 4. Environment Variables & Infrastructure Connections
For local AI development and Docker Compose connectivity, all microservices MUST use the following standard environment variables:

*   **ScyllaDB:** `SCYLLA_HOSTS=127.0.0.1`, `SCYLLA_PORT=9042`, `SCYLLA_KEYSPACE=url_shortener`
*   **Redis:** `REDIS_ADDR=127.0.0.1:6379`
*   **Redpanda (Kafka):** `KAFKA_BROKERS=127.0.0.1:19092`
*   **ClickHouse:** `CLICKHOUSE_ADDR=127.0.0.1:9000` (TCP) or `8123` (HTTP)
*   **ZooKeeper:** `ZOOKEEPER_SERVERS=127.0.0.1:2181`
*   **Service Ports:** `READ_SERVICE_PORT=8081`, `WRITE_SERVICE_PORT=8082`, `KGS_PORT=8083`
*   **Application Config:** `BASE_URL=http://localhost:8081` (Dynamic base for short URLs)
