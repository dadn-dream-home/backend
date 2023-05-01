package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	driver "github.com/mattn/go-sqlite3"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type Database struct {
	olddb *sql.DB
}

var (
	ErrorFeedNotFound = errors.New("feed not found")
	ErrorFeedExists   = errors.New("feed already exists")
)

func NewDatabase() (db Database, err error) {

	// open db with database/sql

	dbPath := "./db.sqlite3"
	migrationsPath := "./migrations"

	olddb, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return Database{}, fmt.Errorf("failed to open %s with database/sql: %w", dbPath, err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, olddb.Close())
		}
	}()

	// open db with driver for migrators

	instance, err := sqlite3.WithInstance(olddb, &sqlite3.Config{})
	if err != nil {
		return Database{}, fmt.Errorf("failed to open %s with driver for migrators: %w", dbPath, err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, instance.Close())
		}
	}()

	// open migrations

	fSrc, err := (&file.File{}).Open(migrationsPath)
	if err != nil {
		return Database{}, fmt.Errorf("failed to open migrations at %s: %w", migrationsPath, err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, fSrc.Close())
		}
	}()

	// init migrator

	m, err := migrate.NewWithInstance("file", fSrc, "sqlite3", instance)
	if err != nil {
		return Database{}, fmt.Errorf("failed to init migrator: %w", err)
	}
	// defer func() {
	// 	serr, dberr := m.Close()
	// 	err = errors.Join(err, serr, dberr)
	// }()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return Database{}, fmt.Errorf("failed to migrate: %w", err)
	}

	return Database{olddb}, nil
}

func (db *Database) Close() error {
	return db.olddb.Close()
}

func (db *Database) InsertFeed(ctx context.Context, id string, ty pb.FeedType) (feed *pb.Feed, err error) {
	log := LoggerFromContext(ctx)
	log = log.WithField("id", id).WithField("type", ty)
	log.Infof("inserting feed\n")
	defer func() {
		if err != nil {
			log.WithError(err).Errorf("error inserting feeds\n")
		} else {
			log.Infof("inserted feed\n")
		}
	}()

	_, err = db.olddb.Exec("INSERT INTO feeds (id, type) VALUES (?, ?)", id, ty)
	if err != nil {
		if errors.Is(err, driver.ErrConstraintUnique) {
			return nil, errors.Join(ErrorFeedExists, err)
		}
		return nil, fmt.Errorf("error inserting feed: %w", err)
	}

	return &pb.Feed{
		Id:   id,
		Type: ty,
	}, nil
}

func (db *Database) ListFeeds(ctx context.Context) (feeds []*pb.Feed, err error) {
	log := LoggerFromContext(ctx)

	log.Infof("listing feeds\n")
	defer func() {
		if err != nil {
			log.WithError(err).Errorf("error listing feeds\n")
		}
	}()

	rows, err := db.olddb.Query("SELECT id, type FROM feeds")
	if err != nil {
		return nil, fmt.Errorf("error querying feeds: %w", err)
	}

	for rows.Next() {
		var feed pb.Feed
		if err := rows.Scan(&feed.Id, &feed.Type); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		feeds = append(feeds, &feed)
	}

	return feeds, nil
}

func (db *Database) DeleteFeed(ctx context.Context, id string) (err error) {
	log := LoggerFromContext(ctx)
	log = log.WithField("feed", id)
	log.Infof("deleting feeds\n")
	defer func() {
		if err != nil {
			log.WithError(err).Errorf("error listing feeds\n")
		} else {
			log.Infof("listed feeds\n")
		}
	}()

	r, err := db.olddb.Exec("DELETE FROM feeds WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting feed: %w", err)
	}

	n, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if n == 0 {
		return ErrorFeedNotFound
	}

	return nil
}

func (db *Database) InsertFeedValue(ctx context.Context, feedId string, value float32) (err error) {
	log := LoggerFromContext(ctx).WithField("feed", feedId).WithField("value", value)
	log.Infof("inserting feed value\n")
	defer func() {
		if err != nil {
			log.WithError(err).Errorf("error inserting feed value\n")
		} else {
			log.Infof("inserted feed value\n")
		}
	}()

	_, err = db.olddb.Exec("INSERT INTO feed_values (feed_id, value) VALUES (?, ?)", feedId, value)
	if err != nil {
		return fmt.Errorf("error inserting feed value: %w", err)
	}
	return nil
}

func (db *Database) SummariseFeedValues(feedId string) (avg float32, err error) {
	rows, err := db.olddb.Query("SELECT AVG(value) FROM feed_values WHERE feed_id = ?", feedId)
	if err != nil {
		return 0, fmt.Errorf("error querying: %w", err)
	}

	if !rows.Next() {
		return 0, fmt.Errorf("no row: %w", rows.Err())
	}

	if err := rows.Scan(&avg); err != nil {
		return 0, fmt.Errorf("error scanning: %w", err)
	}

	return avg, nil
}
