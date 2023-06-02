package database

import (
	"context"
	"database/sql"

	"github.com/dadn-dream-home/x/server/telemetry"
	"go.uber.org/zap"
)

func OpenDatabase(ctx context.Context, config Config) *sql.DB {
	log := telemetry.GetLogger(ctx)

	db, err := sql.Open("sqlite3", config.ConnectionString)
	if err != nil {
		log.Fatal("failed to open database", zap.Error(err))
	}

	log.Info("opened database")

	return db
}
