-- create notifications table

CREATE TABLE notifications (
    feed_id TEXT NOT NULL,
    time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    message TEXT NOT NULL,

    FOREIGN KEY (feed_id) REFERENCES feeds(id)
);

CREATE INDEX notifications_feed_id__time_idx ON notifications(feed_id, time DESC);
