-- create table for sensor config

CREATE TABLE sensor_configs (
    feed_id TEXT PRIMARY KEY,
    has_notification INTEGER NOT NULL DEFAULT 0,
    lower_threshold REAL NOT NULL DEFAULT 10,
    upper_threshold REAL NOT NULL DEFAULT 40,
    FOREIGN KEY(feed_id) REFERENCES feeds(id)
);

CREATE TABLE actuator_configs (
    feed_id TEXT PRIMARY KEY,
    automatic INTEGER NOT NULL DEFAULT 0,
    turn_on_cron_expr TEXT NOT NULL DEFAULT '0 0 0 0 0',
    turn_off_cron_expr TEXT NOT NULL DEFAULT '0 0 0 0 0',
    FOREIGN KEY(feed_id) REFERENCES feeds(id)
);

INSERT INTO sensor_configs (feed_id)
SELECT id FROM feeds
WHERE type = 'TEMPERATURE' OR type = 'HUMIDITY';

INSERT INTO actuator_configs (feed_id)
SELECT id FROM feeds
WHERE type = 'LIGHT';