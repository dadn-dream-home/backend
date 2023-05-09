package main

import (
	"context"

	"github.com/dadn-dream-home/x/server/startup"
	"github.com/dadn-dream-home/x/server/telemetry"
)

func main() {
	ctx := telemetry.InitLogger(context.Background())

	config := startup.OpenConfig(ctx)

	db := startup.OpenDatabase(ctx, config.DatabaseConfig)
	startup.Migrate(ctx, db, config.DatabaseConfig)

	mqtt := startup.ConnectMQTT(ctx, config.MQTTConfig)

	server := startup.NewServer(ctx, db, mqtt)
	lis := startup.Listen(ctx, config.ServerConfig)

	server.Serve(ctx, lis)
}
