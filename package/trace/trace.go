package trace

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/trace"
)

type ReadOnlySpan = sdktrace.ReadOnlySpan
type Exporter = sdktrace.SpanExporter
type Provider = trace.TracerProvider
type SpanID = trace.SpanID

func New(exporter Exporter) *Tracer {
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter)),
	)
	return &Tracer{t: provider.Tracer("")}
}

type Tracer struct {
	t trace.Tracer
}

// Start a trace
// TODO: support key-value attributes
func (t *Tracer) Start(ctx context.Context, label string) (context.Context, *Span) {
	ctx, span := t.t.Start(ctx, label)
	return ctx, &Span{s: span}
}

type Span struct {
	s trace.Span
}

// End the span
func (s *Span) End(err *error) {
	if err != nil {
		s.s.RecordError(*err)
	}
	s.s.End()
}

func Encode(ctx context.Context) ([]byte, error) {
	tp := propagation.TraceContext{}
	carrier := propagation.MapCarrier{}
	tp.Inject(ctx, carrier)
	return json.Marshal(carrier)
}

func Decode(ctx context.Context, data []byte) (context.Context, error) {
	tp := propagation.TraceContext{}
	carrier := propagation.MapCarrier{}
	if err := json.Unmarshal(data, &carrier); err != nil {
		return nil, err
	}
	return tp.Extract(ctx, carrier), nil
}
