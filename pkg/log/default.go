package log

import (
	"context"
	"os"

	"github.com/livebud/bud/pkg/color"
)

func DefaultLevel() Level {
	return LevelInfo
}

func Default() *Logger {
	return New(Filter(LevelInfo, &Console{
		Writer:    os.Stderr,
		Color:     color.Default(),
		AddSource: true,
	}))
}

var defaultLog = Default()

func Debug(msg string, args ...any) {
	defaultLog.Debug(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	defaultLog.DebugContext(ctx, msg, args...)
}

func Info(msg string, args ...any) {
	defaultLog.Info(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	defaultLog.InfoContext(ctx, msg, args...)
}

func Warn(msg string, args ...any) {
	defaultLog.Warn(msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	defaultLog.WarnContext(ctx, msg, args...)
}

func Error(msg string, args ...any) {
	defaultLog.Error(msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	defaultLog.ErrorContext(ctx, msg, args...)
}
