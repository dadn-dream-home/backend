package main

import (
	"context"

	"github.com/dadn-dream-home/x/server/services"
	"github.com/dadn-dream-home/x/server/telemetry"
)

func main() {
	ctx := context.Background()
	ctx = telemetry.InitLogger(ctx)
	log := telemetry.GetLogger(ctx)

	grpcService, err := services.NewBackendService(ctx)
	if err != nil {
		log.WithError(err).Fatal("failed to create backend service")
	}
	grpcService.Serve(ctx)
}
