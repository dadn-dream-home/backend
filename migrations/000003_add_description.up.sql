-- add description column to table feeds
ALTER TABLE feeds ADD COLUMN description TEXT NOT NULL DEFAULT '';
