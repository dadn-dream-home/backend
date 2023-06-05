package main

import (
	"context"

	"github.com/dadn-dream-home/x/server/startup"
	"github.com/dadn-dream-home/x/server/startup/database"
	"github.com/dadn-dream-home/x/server/telemetry"
)

func main() {
	ctx := telemetry.InitLogger(context.Background())

	config := startup.OpenConfig(ctx)

	hooker := database.RegisterHook(ctx, config.DatabaseConfig)
	db := database.OpenDatabase(ctx, config.DatabaseConfig)
	database.Migrate(ctx, db, config.DatabaseConfig)

	mqtt := startup.ConnectMQTT(ctx, config.MQTTConfig)

	server := startup.NewServer(ctx, db, mqtt, hooker)
	lis := startup.Listen(ctx, config.ServerConfig)

	server.Serve(ctx, lis)
}
