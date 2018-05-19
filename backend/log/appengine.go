package log

import (
	"context"

	appenginelog "google.golang.org/appengine/log"
)

type AppEngineLogger struct{}

func (l *AppEngineLogger) Info(ctx context.Context, msg string) {
	appenginelog.Infof(ctx, msg)
}

func (l *AppEngineLogger) Debug(ctx context.Context, msg string) {
	appenginelog.Debugf(ctx, msg)
}

func (l *AppEngineLogger) Warn(ctx context.Context, msg string) {
	appenginelog.Warningf(ctx, msg)
}

func (l *AppEngineLogger) Error(ctx context.Context, msg string) {
	appenginelog.Errorf(ctx, msg)
}
