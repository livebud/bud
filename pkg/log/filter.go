package log

import (
	"context"
	"log/slog"
)

// Filter logs by level
func Filter(level Level, handler Handler) Handler {
	return &filter{level, handler}
}

type filter struct {
	level   Level
	handler Handler
}

func (f *filter) Enabled(ctx context.Context, l Level) bool {
	return l >= f.level
}

func (f *filter) Handle(ctx context.Context, record slog.Record) error {
	return f.handler.Handle(ctx, record)
}

func (f *filter) WithAttrs(attrs []slog.Attr) Handler {
	return &filter{
		level:   f.level,
		handler: f.handler.WithAttrs(attrs),
	}
}

func (f *filter) WithGroup(group string) Handler {
	return &filter{
		level:   f.level,
		handler: f.handler.WithGroup(group),
	}
}
