package main

import (
	"context"
	"path"
	"runtime"
	"strconv"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type contextKey int

const (
	RequestIdKey contextKey = iota
	LoggerKey
)

func LoggerFromContext(ctx context.Context) *logrus.Entry {
	logger, ok := ctx.Value(LoggerKey).(*logrus.Entry)
	if !ok {
		logger = logrus.NewEntry(log)
	}
	return logger
}

func InitLogging() {
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
}

func requestIdUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// keep the request id short
	rid := uuid.New().String()[:6]
	logger := log.WithField("reqid", rid)
	ctx = context.WithValue(ctx, LoggerKey, logger)
	ctx = context.WithValue(ctx, RequestIdKey, rid)

	return handler(ctx, req)
}

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (ss *serverStream) Context() context.Context {
	return ss.ctx
}

func requestIdStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	// keep the request id short
	rid := uuid.New().String()[:6]
	logger := log.WithField("rid", rid)
	ctx = context.WithValue(ctx, LoggerKey, logger)
	ctx = context.WithValue(ctx, RequestIdKey, rid)

	return handler(srv, &serverStream{ss, ctx})
}
