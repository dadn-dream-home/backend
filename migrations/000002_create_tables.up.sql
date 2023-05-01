CREATE TABLE feeds (
    id TEXT PRIMARY KEY,
    type INTEGER NOT NULL
);

CREATE TABLE feed_values (
    time TIMESTAMPT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    feed_id TEXT NOT NULL,
    value REAL NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES sensor_feeds(id)
);

CREATE INDEX feed_values_feed_id_time_idx ON feed_values (feed_id, time DESC);