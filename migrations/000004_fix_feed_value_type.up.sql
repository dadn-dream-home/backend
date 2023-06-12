-- rename old table to value_old
ALTER TABLE feed_values RENAME TO feed_values_old;

-- create new table
CREATE TABLE feed_values (
    time TIMESTAMPT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    feed_id TEXT NOT NULL,
    value BLOB NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feeds(id)
);

-- create index
DROP INDEX feed_values_feed_id_time_idx;
CREATE INDEX feed_values_feed_id_time_idx ON feed_values (feed_id, time DESC);

-- copy data from old table to new table
INSERT INTO feed_values (feed_id, value, time)
SELECT feed_id, CAST(value as BLOB), time FROM feed_values_old;

-- drop old table
DROP TABLE feed_values_old;