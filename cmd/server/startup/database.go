package startup

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dadn-dream-home/x/server/telemetry"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func OpenDatabase(ctx context.Context, config DatabaseConfig) *sql.DB {
	log := telemetry.GetLogger(ctx)

	db, err := sql.Open("sqlite3", config.ConnectionString)
	if err != nil {
		log.Fatal("failed to open database", zap.Error(err))
	}

	log.Info("opened database")

	return db
}

func Migrate(ctx context.Context, db *sql.DB, config DatabaseConfig) {
	log := telemetry.GetLogger(ctx)

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal("failed to init migrator driver", zap.Error(err))
	}
	defer instance.Close()

	log.Info("initialized migrator driver")

	// open migrations

	migrations, err := (&file.File{}).Open(config.MigrationsPath)
	if err != nil {
		log.Fatal("failed to open migrations", zap.Error(err), zap.String("path", config.MigrationsPath))
	}
	defer migrations.Close()

	log.Info("opened migrations", zap.String("path", config.MigrationsPath))

	// init migrator

	m, err := migrate.NewWithInstance("file", migrations, "sqlite3", instance)
	if err != nil {
		log.Fatal("failed to init migrator", zap.Error(err))
	}
	defer m.Close()

	log.Info("initialized migrator")

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("failed to migrate", zap.Error(err))
	}

	log.Info("migrated database")
}
