# Database Schema Blueprints

This document defines the strict data contracts for the URL Shortener architecture. All database migrations, queries, and ORM/driver implementations must adhere to these schemas.

## 1. Primary Database: ScyllaDB

**Keyspace Strategy:** `NetworkTopologyStrategy` (Active-Active)
**Table:** `url_mapping`

*   **Partition Key:** `short_code` (Must be exact-match for <10ms lookups)
*   **TTL:** Managed natively by ScyllaDB using `default_time_to_live` or explicitly passed during INSERT.

**CQL Blueprint:**
```sql
CREATE KEYSPACE IF NOT EXISTS url_shortener 
WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': 3};

CREATE TABLE IF NOT EXISTS url_shortener.url_mapping (
    short_code text PRIMARY KEY,
    original_url text,
    user_id text,
    created_at timestamp,
    expiration_date timestamp
) WITH default_time_to_live = 31536000; -- Default TTL of 1 year, override on insert if needed
```

## 2. Analytics Database: ClickHouse

**Strategy:** The Read Service pushes to Redpanda, and Consumer Workers bulk-insert into ClickHouse. We use `ReplacingMergeTree` for idempotent inserts and deduplication.
**Table:** `clicks_raw`

*   **Engine:** `ReplacingMergeTree` (uses timestamp to deduplicate identical rows inserted during consumer retries)
*   **Order By:** `(short_code, timestamp)` (Optimized for filtering analytics by specific short codes)

**SQL Blueprint:**
```sql
CREATE TABLE IF NOT EXISTS clicks_raw (
    click_id UUID,
    short_code String,
    timestamp DateTime,
    ip_address String,
    user_agent String,
    referer String,
    country_code FixedString(2)
) ENGINE = ReplacingMergeTree(timestamp)
ORDER BY (short_code, timestamp);
```

## 3. Cache Strategy: Redis
*   **Data Structure:** Simple Key-Value strings.
*   **Key Format:** `url:{short_code}`
*   **Value Format:** The `original_url` string.
*   **TTL Rule:** The Read Service MUST refresh the TTL on a cache hit (Rolling expiration).

## 4. API Contracts

### Write Service: `POST /shorten`
**Headers:**
- `Idempotency-Key`: UUID (Required)

**Request Body:**
```json
{
  "original_url": "https://example.com/very/long/path?q=123",
  "custom_alias": "my-promo" 
}
```
*(Note: `custom_alias` is optional)*

**Success Response (201 Created):**
```json
{
  "short_code": "my-promo",
  "short_url": "http://localhost:8081/my-promo",
  "original_url": "https://example.com/very/long/path?q=123",
  "expires_at": "2027-07-15T00:00:00Z"
}
```
*(Note: `short_url` domain is dynamically generated using the `BASE_URL` environment variable)*

**Error Responses:**
- `400 Bad Request`: Invalid URL format or missing idempotency key.
- `409 Conflict`: Custom alias already exists OR the `Idempotency-Key` was reused with a *different* `original_url`.
- `429 Too Many Requests`: Rate limit exceeded.
- `500 Internal Server Error`: Standard server error.

### Read Service: `GET /{short_code}`
**Responses:**
- `302 Found`: Standard traffic, redirects to `original_url`.
- `301 Moved Permanently`: Viral spike traffic (with aggressive `Cache-Control` headers).
- `404 Not Found`: Short code does not exist or has expired.
