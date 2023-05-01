package main

import (
	"path"
	"runtime"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: false,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// log filename and linenumber only
			fileName := path.Base(f.File) + ":" + strconv.Itoa(f.Line)
			return f.Function, fileName
		},
	})
	log.SetReportCaller(true)

	service, err := NewBackendGrpcService()
	if err != nil {
		log.WithError(err).Fatal("failed to create backend grpc service")
	}

	service.Init()
	service.Serve()
	wg.Wait()
}
