# Phase 5 Microservices Implementation Plan

This plan details the step-by-step implementation for Phase 5 of the URL Shortener project. It follows strict Test-Driven Development (TDD), automated verifications, explicit manual user verification points, and Git checkpoints.

## Proposed Changes

### Phase 5.1: Infrastructure & Scaffolding
- [x] **Step 1:** Initialize Go module (`go mod init urlshortener`).
- [x] **Step 2:** Create the monorepo directory structure (`/cmd`, `/internal`, `/pkg`, `/deploy`, `/scripts`).
- [x] **Step 3:** Create ScyllaDB schema at `deploy/init.cql` (keyspace, url_mapping table).
- [x] **Step 4:** Create ClickHouse schema at `deploy/init.sql` (clicks_raw table).
- [x] **Step 5:** Create `docker-compose.yml` defining ScyllaDB, RedisStack (for RedisBloom), Redpanda, ClickHouse, and ZooKeeper.
- [x] **Step 6:** **AI Verification**: Run `docker-compose config` and `go mod tidy` to validate configs.
- [x] **Step 7:** **WAIT FOR USER**: Run `docker-compose up -d` and wait for containers to be fully healthy. Wait for user "go ahead".
- [x] **Step 8:** **AI Verification**: Execute DB schemas using `docker exec -i urlshortener-scylla cqlsh < deploy/init.cql` and `docker exec -i urlshortener-clickhouse clickhouse-client < deploy/init.sql`.
- [x] **Step 9:** **[USER ACTION] Git Commit**: "chore: initialize go module, monorepo structure, and docker-compose infrastructure"

### Phase 5.2: Shared Packages (`/pkg`)
- [x] **Step 10:** Implement Environment Config parser in `/pkg/config/env.go` to load DB/Cache hosts and ports.
- [x] **Step 11:** Write unit tests for ScyllaDB connection wrapper.
- [x] **Step 12:** Implement ScyllaDB driver setup in `/pkg/db/scylla.go`.
- [x] **Step 13:** Write unit tests for Redis connection wrapper.
- [x] **Step 14:** Implement Redis client setup in `/pkg/cache/redis.go`.
- [x] **Step 15:** Write unit tests for Redpanda publisher wrapper.
- [x] **Step 16:** Implement Redpanda publisher in `/pkg/kafka/publisher.go`.
- [x] **Step 17:** Implement KGS HTTP client in `/pkg/kgsclient/client.go` to fetch Snowflake IDs.
- [x] **Step 18:** **AI Verification**: Run `go test ./pkg/...` and `go vet ./pkg/...`.
- [ ] **Step 19:** **[USER ACTION] Git Commit**: "feat: implement shared config, db, cache, kafka, and KGS client packages"

### Phase 5.3: Key Generation Service (KGS)
- [ ] **Step 20:** Write unit tests for Snowflake ID logic in `/internal/kgs/snowflake_test.go`.
- [ ] **Step 21:** Implement Snowflake ID generation in `/internal/kgs/snowflake.go`.
- [ ] **Step 22:** Write unit tests for Base62 encoding in `/internal/kgs/base62_test.go`.
- [ ] **Step 23:** Implement Base62 encoding in `/internal/kgs/base62.go`.
- [ ] **Step 24:** Write unit tests for ZooKeeper Worker ID fetcher in `/internal/kgs/zk_test.go`.
- [ ] **Step 25:** Implement ZooKeeper Worker ID fetcher in `/internal/kgs/zk.go`.
- [ ] **Step 26:** Implement KGS HTTP handler and router in `/internal/kgs/handler.go`.
- [ ] **Step 27:** Create KGS entrypoint in `/cmd/kgs/main.go` wiring up config, zk, and HTTP server.
- [ ] **Step 28:** **AI Verification**: Run `go test ./internal/kgs/...`, `go build ./cmd/kgs`, and `go vet ./...`.
- [ ] **Step 29:** **WAIT FOR USER**: Run the KGS locally (`go run ./cmd/kgs/main.go`). Test via `curl http://localhost:8083/`. Wait for user "go ahead".
- [ ] **Step 30:** **[USER ACTION] Git Commit**: "feat: implement Key Generation Service (KGS) with ZooKeeper worker ID assignment"

### Phase 5.4: Write Service (Shortening)
- [ ] **Step 31:** Write unit tests for Bloom Filter logic in `/internal/write/bloom_test.go`.
- [ ] **Step 32:** Implement RedisBloom check and add logic in `/internal/write/bloom.go`.
- [ ] **Step 33:** Write unit tests for idempotency caching in `/internal/write/idempotency_test.go`.
- [ ] **Step 34:** Implement idempotency caching in `/internal/write/idempotency.go`.
- [ ] **Step 35:** Write unit tests for ScyllaDB repository (store mapping) in `/internal/write/repository_test.go`.
- [ ] **Step 36:** Implement ScyllaDB repository in `/internal/write/repository.go`.
- [ ] **Step 37:** Write unit tests for core shortening service logic in `/internal/write/service_test.go`.
- [ ] **Step 38:** Implement core shortening service logic (auto-gen vs custom alias) in `/internal/write/service.go`.
- [ ] **Step 39:** Implement HTTP handler for `POST /shorten` in `/internal/write/handler.go`.
- [ ] **Step 40:** Create Write Service entrypoint in `/cmd/write-service/main.go`.
- [ ] **Step 41:** **AI Verification**: Run `go test ./internal/write/...`, `go build ./cmd/write-service`, and `go vet ./...`.
- [ ] **Step 42:** **WAIT FOR USER**: Run Write Service locally. Send a `POST /shorten` request via cURL/Postman. Verify DB/Cache entries. Wait for user "go ahead".
- [ ] **Step 43:** **[USER ACTION] Git Commit**: "feat: implement Write Service with idempotency, Bloom filter, and ScyllaDB persistence"

### Phase 5.5: Read Service (Redirection)
- [ ] **Step 44:** Write unit tests for Redis cache read and TTL refresh in `/internal/read/cache_test.go`.
- [ ] **Step 45:** Implement Redis cache read and TTL refresh in `/internal/read/cache.go`.
- [ ] **Step 46:** Write unit tests for Mutex lock cache stampede prevention in `/internal/read/mutex_test.go`.
- [ ] **Step 47:** Implement Mutex lock logic in `/internal/read/mutex.go`.
- [ ] **Step 48:** Write unit tests for ScyllaDB repository (get mapping) in `/internal/read/repository_test.go`.
- [ ] **Step 49:** Implement ScyllaDB repository in `/internal/read/repository.go`.
- [ ] **Step 50:** Write unit tests for core redirection logic in `/internal/read/service_test.go`.
- [ ] **Step 51:** Implement core redirection service logic (Cache Hit vs Miss) in `/internal/read/service.go`.
- [ ] **Step 52:** Implement HTTP handler for `GET /{short_code}` in `/internal/read/handler.go` (including 302/301 logic).
- [ ] **Step 53:** Create Read Service entrypoint in `/cmd/read-service/main.go`.
- [ ] **Step 54:** **AI Verification**: Run `go test ./internal/read/...`, `go build ./cmd/read-service`, and `go vet ./...`.
- [ ] **Step 55:** **WAIT FOR USER**: Run Read Service locally. Access the short URL via browser/cURL. Verify redirection and Redis TTL refresh. Wait for user "go ahead".
- [ ] **Step 56:** **[USER ACTION] Git Commit**: "feat: implement Read Service with dynamic caching and stampede prevention"

### Phase 5.6: Analytics Workers
- [ ] **Step 57:** Write unit tests for Redpanda consumer logic in `/internal/analytics/consumer_test.go`.
- [ ] **Step 58:** Implement Redpanda consumer worker in `/internal/analytics/consumer.go`.
- [ ] **Step 59:** Write unit tests for ClickHouse bulk insertion logic in `/internal/analytics/clickhouse_test.go`.
- [ ] **Step 60:** Implement ClickHouse bulk inserter in `/internal/analytics/clickhouse.go`.
- [ ] **Step 61:** Create Analytics Worker entrypoint in `/cmd/analytics-worker/main.go`.
- [ ] **Step 62:** **AI Verification**: Run `go test ./internal/analytics/...`, `go build ./cmd/analytics-worker`, and `go vet ./...`.
- [ ] **Step 63:** **WAIT FOR USER**: Run the Analytics Worker locally. Click short links, verify deduplicated bulk inserts in ClickHouse. Wait for user "go ahead".
- [ ] **Step 64:** **[USER ACTION] Git Commit**: "feat: implement Analytics Workers for Redpanda consumption and ClickHouse insertion"

## Verification Plan

### Automated Tests
- `go test ./...` and `go vet ./...` will be executed frequently by the AI after each microservice feature.
- `go build` will be used to ensure compiling states before manual tests.

### Manual Verification
- You will be explicitly prompted to perform manual tests via `curl`, Postman, or browser at the end of each major checkpoint.
- You should inspect the database state (ScyllaDB, Redis, ClickHouse, Redpanda) during these pauses.
