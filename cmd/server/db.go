package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type Database struct {
	olddb *sql.DB
}

func NewDatabase() Database {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fSrc, err := (&file.File{}).Open("./migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("file", fSrc, "sqlite3", instance)
	if err != nil {
		log.Fatal(err)
	}

	m.Up()

	return Database{db}
}

func (db *Database) Close() {
	db.olddb.Close()
}

func (db *Database) InsertFeed(name string, ty pb.FeedType) *pb.Feed {
	id := uuid.New().String()
	_, err := db.olddb.Exec("INSERT INTO feeds (id, name, type) VALUES (?, ?, ?)", id, name, ty)
	if err != nil {
		panic(err)
	}

	return &pb.Feed{
		Id:   id,
		Name: name,
		Type: ty,
	}
}

func (db *Database) ListFeeds() []*pb.Feed {
	rows, err := db.olddb.Query("SELECT id, name, type FROM feeds")
	if err != nil {
		panic(err)
	}

	var feeds []*pb.Feed
	for rows.Next() {
		var feed pb.Feed
		if err := rows.Scan(&feed.Id, &feed.Name, &feed.Type); err != nil {
			panic(err)
		}
		feeds = append(feeds, &feed)
	}

	return feeds
}

func (db *Database) DeleteFeed(id string) error {
	_, err := db.olddb.Exec("DELETE FROM feeds WHERE id = ?", id)
	return err
}

func (db *Database) InsertFeedValue(feedId string, value float32) error {
	_, err := db.olddb.Exec("INSERT INTO feed_values (feed_id, value) VALUES (?, ?)", feedId, value)
	return err
}

func (db *Database) SummariseFeedValues(feedId string) (avg float32) {
	rows, err := db.olddb.Query("SELECT AVG(value) FROM feed_values WHERE feed_id = ?", feedId)
	if err != nil {
		panic(err)
	}

	if !rows.Next() {
		panic(rows.Err())
	}

	if err := rows.Scan(&avg); err != nil {
		panic(err)
	}

	return avg
}
