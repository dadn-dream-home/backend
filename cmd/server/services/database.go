package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
)

func NewDatabase() (db *sql.DB, err error) {

	// open db with database/sql

	dbPath := "./db.sqlite3"
	migrationsPath := "./migrations"

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s with database/sql: %w", dbPath, err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, db.Close())
		}
	}()

	// open db with driver for migrators

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open %s with driver for migrators: %w", dbPath, err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, instance.Close())
		}
	}()

	// open migrations

	fSrc, err := (&file.File{}).Open(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open migrations at %s: %w", migrationsPath, err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, fSrc.Close())
		}
	}()

	// init migrator

	m, err := migrate.NewWithInstance("file", fSrc, "sqlite3", instance)
	if err != nil {
		return nil, fmt.Errorf("failed to init migrator: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return db, nil
}
