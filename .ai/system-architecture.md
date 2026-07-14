# Production-Ready URL Shortener Architecture v2.0 (High-Scale PoC Edition)

This updated architecture resolves extreme-scale bottlenecks by implementing stateless key generation, dynamic edge caching for viral spikes, and optimized partition routing. It is designed to handle billions of URLs and massive read-heavy traffic with strict consistency and fault tolerance.

## 1. High-Level Component Stack

| Layer | Recommended Tech | Purpose |
| :--- | :--- | :--- |
| **Edge / DNS** | Cloudflare (Free Tier) | Global DNS routing, basic DDoS protection, and dynamic edge caching. |
| **Load Balancer** | NGINX (Open Source) | SSL termination, reverse proxy, and local multi-service routing. |
| **API Gateway** | NGINX / Traefik | Local rate limiting, path routing, and request validation. |
| **Compute** | Docker Compose / Local K8s | Containerized microservices orchestration for local PoC testing. |
| **Cache** | Redis (Open Source) | Ultra-fast key-value store with distributed locks and TTL management. |
| **Primary DB** | ScyllaDB / PostgreSQL | Horizontally scalable datastore. |
| **Coordination** | ZooKeeper / etcd | Distributed coordination for assigning Worker IDs to the KGS. |
| **Message Broker** | Redpanda / Apache Kafka | Asynchronous queue using random partitioning for analytics ingestion. |
| **Analytics DB** | ClickHouse (Community Edition) | Columnar database optimized for heavy aggregation queries and deduplication. |
| **Observability** | Prometheus, Grafana, OpenTelemetry | Distributed tracing, metrics collection, and dashboard visualization. |

---

## 2. The Core Microservices

### The Read Service (Redirection & Dynamic Caching)
*   **Role:** Handles `GET /{code}` requests with latency <10ms.
*   **Logic:**
    1.  Checks Redis cache for the `{code}`.
    2.  **Cache Hit:** If found, refreshes the Redis TTL to ensure highly trafficked links remain in memory.
    3.  **Cache Miss & Stampede Prevention:** Uses a Distributed Lock (Mutex) so only one thread queries the Primary DB.
    4.  **Analytics Ingestion:** Drops a click event payload into Redpanda/Kafka using a `RoundRobinPartitioner`. 
    5.  **Dynamic Response Strategy (Viral Spike Protection):**
        *   *Standard Traffic:* Returns a `302 Found` to ensure the click hits the backend for analytics.
        *   *Viral Spike Detected:* If the read rate for a specific code crosses a critical threshold (e.g., >10k req/sec), the service automatically switches to returning a `301 Permanent Redirect` or sets aggressive `Cache-Control` headers, offloading the traffic to Cloudflare and protecting the backend.

### The Write Service (Shortening)
*   **Role:** Handles `POST /shorten`.
*   **Logic:**
    1.  Validates the URL and checks an `Idempotency-Key`. The mapping of `{Idempotency-Key} -> {short_code}` is temporarily cached. If a client crashes and retries with the same key, it receives the exact same short code without duplicate processing.
    2.  **Duplicate URL Check:** Queries a Redis-backed Bloom Filter to check if the long URL already exists. If positive, it queries the DB to return the existing code, saving massive namespace exhaustion from bot abuse.
    3.  **Branch A (Auto-Generated):** Requests a pre-generated code from the KGS and writes the `{code} -> {long_URL}` mapping.
    4.  **Branch B (Custom Alias):** Routes all custom alias write requests to a designated "Leader" region to avoid Paxos Lightweight Transaction (LWT) cross-region consensus penalties.

### Key Generation Service (KGS) - Stateless
*   **Role:** Generates perfectly collision-free strings without a database bottleneck.
*   **Logic:** Replaces the expensive "Unused Codes" table with a stateless generator.
*   1.  Upon startup, each KGS node requests a unique `Worker ID` from ZooKeeper/etcd.
    2.  The node combines its `Worker ID`, a timestamp, and a local atomic counter to generate a unique integer (similar to Twitter Snowflake).
    3.  This integer is encoded into a Base62 string. This completely eliminates the need to store trillions of unused strings or rely on a centralized sequence database.

---

## 3. Data Schema & Strategy

### URL Mapping Table (ScyllaDB / PostgreSQL)
*   **Partition Key / Primary Key:** `short_code` (Ensures fast, exact-match lookups).
*   **Columns:** `original_url`, `created_at`, `expiration_date`, `user_id`.
*   **Routing:** Active-Active for reads; custom alias writes routed to a Leader region.
*   **TTL (Time to Live):** Uses native NoSQL TTL features in ScyllaDB or background workers in PostgreSQL to purge expired links.

### Cache Strategy (Redis)
*   **Eviction Policy:** `volatile-lru`.
*   **Rule:** Assign explicit TTLs to all cached mappings. Highly trafficked links have their TTL refreshed upon access (rolling expiration), ensuring vital links are never evicted to make room for one-off new links.

---

## 4. The Analytics Pipeline (Async)

1.  The Read Service pushes raw JSON click events to **Redpanda/Kafka**.
2.  **Hot Partition Mitigation:** Producers use a `RoundRobinPartitioner` to perfectly distribute the load across all brokers, ignoring IP or short code to prevent CGNAT hotspots.
3.  A fleet of Consumer Workers bulk-inserts the raw data into **ClickHouse**.
4.  ClickHouse utilizes Materialized Views and `ReplacingMergeTree` engines to handle data deduplication and real-time aggregations downstream.

---

## 5. Security & Abuse Protections

*   **Rate Limiting & WAF:** Enforced at the API Gateway and Edge layers (Cloudflare + NGINX rate limiting).
*   **Real-Time Malicious Sweeps:** Instead of delayed cron jobs, an asynchronous stream processor consumes a randomized sample of newly created links and active click streams. It immediately evaluates destination URLs against the Google Safe Browsing API, instantly flagging and disabling payloads that rotate to malware post-creation.

---

## 6. Observability & SRE

*   **Distributed Tracing:** OpenTelemetry traces requests across the Gateway, Read Service, Cache, and DB.
*   **Metrics:** Prometheus scrapes cache hit/miss ratios, Kafka consumer lag, and API latencies.
*   **Dashboards & Alerts:** Grafana visualizes metrics with configured webhook alerts (Discord/Slack/Email) for anomalous thresholds like sudden traffic spikes or KGS node failures.

---

## 7. Component Interaction Boundaries

### 7.1 The Read Service

#### Cache Hit Flow
```
Client -> Cloudflare -> NGINX -> Read Service -> Redis (Cache Hit & Refresh TTL) -> Read Service -> Client (302 Found / 301 Redirect)
```
*Asynchronous event logged concurrently:*
```
Read Service -> Redpanda/Kafka (RoundRobinPartitioner)
```

#### Cache Miss Flow
```
Client -> Cloudflare -> NGINX -> Read Service -> Redis (Cache Miss) -> Redis (Acquire Mutex Lock) -> Primary DB (ScyllaDB/PostgreSQL lookup) -> Redis (Cache Mapping & Release Lock) -> Read Service -> Client (302 Found / 301 Redirect)
```
*Asynchronous event logged concurrently:*
```
Read Service -> Redpanda/Kafka (RoundRobinPartitioner)
```

---

### 7.2 The Write Service

#### Auto-Generated Flow
```
Client -> Cloudflare -> NGINX -> Write Service -> Redis (Verify Idempotency & Bloom Filter Check) -> KGS -> Write Service -> Primary DB (Store Mapping) -> Write Service -> Client (Short Code Response)
```
*Background KGS Startup worker coordination:*
```
KGS Node -> ZooKeeper/etcd (Fetch Worker ID registration)
```

#### Custom Alias Flow
```
Client -> Cloudflare -> NGINX -> Write Service -> Redis (Verify Idempotency & Bloom Filter Check) -> Write Service (Route to Leader Region) -> Primary DB (LWT/Consensus & Store Mapping) -> Write Service -> Client (Short Code Response)
```

---

### 7.3 The Analytics Pipeline
```
Read Service -> Redpanda/Kafka (RoundRobinPartitioner) -> Consumer Workers -> ClickHouse (Materialized Views & ReplacingMergeTree deduplication/aggregation)
```
