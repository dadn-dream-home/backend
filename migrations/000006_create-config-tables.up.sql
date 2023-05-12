-- create table for sensor config

CREATE TABLE sensor_config (
    sensor_id TEXT PRIMARY KEY,
    has_notification INTEGER NOT NULL,
    lower_threshold REAL NOT NULL,
    upper_threshold REAL NOT NULL,
    FOREIGN KEY(sensor_id) REFERENCES feeds(id)
);

CREATE TABLE actuator_config (
    actuator_id TEXT PRIMARY KEY,
    automatic INTEGER NOT NULL,
    turn_on_cron_expr TEXT NOT NULL,
    turn_off_cron_expr TEXT NOT NULL,
    FOREIGN KEY(actuator_id) REFERENCES feeds(id)
);