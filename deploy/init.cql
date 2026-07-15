CREATE KEYSPACE IF NOT EXISTS url_shortener 
WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': 1};

CREATE TABLE IF NOT EXISTS url_shortener.url_mapping (
    short_code text PRIMARY KEY,
    original_url text,
    user_id text,
    created_at timestamp,
    expiration_date timestamp
) WITH default_time_to_live = 31536000; -- Default TTL of 1 year, override on insert if needed
