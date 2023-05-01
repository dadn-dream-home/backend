package main

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

var wg sync.WaitGroup
var log = logrus.New()

func main() {
	ctx := context.Background()

	InitLogging()

	service, err := NewBackendGrpcService()
	if err != nil {
		log.WithError(err).Fatal("failed to create backend grpc service")
	}

	service.Init(ctx)
	service.Serve()
	wg.Wait()
}
