-- add threshold columns
ALTER TABLE sensor_configs
ADD COLUMN lower_has_trigger INTEGER NOT NULL DEFAULT 0;

ALTER TABLE sensor_configs
ADD COLUMN lower_feed_id TEXT
    CONSTRAINT fk_sensor_configs_lower_feed_id
    REFERENCES feeds (id);

ALTER TABLE sensor_configs
ADD COLUMN lower_state INTEGER NOT NULL DEFAULT 0;

ALTER TABLE sensor_configs
ADD COLUMN upper_has_trigger INTEGER NOT NULL DEFAULT 0;

ALTER TABLE sensor_configs
ADD COLUMN upper_feed_id TEXT
    CONSTRAINT fk_sensor_configs_upper_feed_id
    REFERENCES feeds (id);

ALTER TABLE sensor_configs
ADD COLUMN upper_state INTEGER NOT NULL DEFAULT 0;
