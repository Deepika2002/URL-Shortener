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
