package telemetry

import (
	"context"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

func InitLogger(ctx context.Context) context.Context {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: false,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// log filename and linenumber only
			fileName := path.Base(f.File) + ":" + strconv.Itoa(f.Line)
			return "", fileName
		},
	})
	log.SetReportCaller(true)
	return ContextWithLogger(ctx, log)
}
