package tracer

import (
	"context"
)

type discard struct{}

var _ Trace = (*discard)(nil)

func (discard) Start(ctx context.Context, label string, attrs ...interface{}) (context.Context, Span) {
	return ctx, discardSpan{}
}

type discardSpan struct{}

func (discardSpan) End(*error) {}
