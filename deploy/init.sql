CREATE TABLE IF NOT EXISTS analytics_events (
    short_id String,
    long_url String,
    event_type String,
    timestamp DateTime
) ENGINE = MergeTree()
ORDER BY (short_id, timestamp);
