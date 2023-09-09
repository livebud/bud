package log

import (
	"context"

	"log/slog"

	"golang.org/x/sync/errgroup"
)

func Multi(handlers ...slog.Handler) *Logger {
	return New(&handler{handlers})
}

type handler struct {
	handlers []slog.Handler
}

var _ slog.Handler = (*handler)(nil)

func (h *handler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return true
}

func (h *handler) Handle(ctx context.Context, record slog.Record) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, handler := range h.handlers {
		handler := handler
		if !handler.Enabled(ctx, record.Level) {
			continue
		}
		eg.Go(func() error { return handler.Handle(ctx, record.Clone()) })
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	multi := &handler{}
	for _, handler := range h.handlers {
		multi.handlers = append(multi.handlers, handler.WithAttrs(attrs))
	}
	return multi
}

func (h *handler) WithGroup(group string) slog.Handler {
	multi := &handler{}
	for _, handler := range h.handlers {
		multi.handlers = append(multi.handlers, handler.WithGroup(group))
	}
	return multi
}
