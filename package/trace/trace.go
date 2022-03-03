package trace

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

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

type keyValues []attribute.KeyValue

func (kv keyValues) Len() int {
	return len(kv)
}

func (kv keyValues) Less(i, j int) bool {
	return kv[i].Key < kv[j].Key
}

func (kv keyValues) Swap(i, j int) {
	kv[i], kv[j] = kv[j], kv[i]
}

func attributes(kvs []interface{}) (list keyValues) {
	size := len(kvs)
	// Special cases
	if size == 0 {
		return nil
	} else if size == 1 {
		return []attribute.KeyValue{attribute.String(fmt.Sprintf("%s", kvs[0]), "")}
	}
	for i := 1; i < size; i += 2 {
		list = append(list, attribute.String(
			fmt.Sprintf("%s", kvs[i-1]),
			fmt.Sprintf("%v", kvs[i]),
		))
	}
	// Sort the fields by key
	sort.Sort(list)
	return list
}

// Start a trace
// TODO: support key-value attributes
func (t *Tracer) Start(ctx context.Context, label string, attrs ...interface{}) (context.Context, *Span) {
	ctx, span := t.t.Start(ctx, label)
	return ctx, &Span{s: span, kvs: attrs}
}

type Span struct {
	s   trace.Span
	kvs []interface{}
}

// End the span
func (s *Span) End(err *error) {
	// Set the attributes
	attrs := attributes(s.kvs)
	s.s.SetAttributes(attrs...)
	// Record an error if it occurs
	if *err != nil {
		s.s.SetStatus(codes.Error, (*err).Error())
		s.s.RecordError(*err)
		s.s.End()
		return
	}
	// Otherwise mark okay and end the span
	s.s.SetStatus(codes.Ok, "")
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

// ToSpanData turns a ReadOnlySpan into *SpanData
func ToSpanData(span ReadOnlySpan) *SpanData {
	data := &SpanData{
		ID:       span.SpanContext().SpanID().String(),
		Name:     span.Name(),
		ParentID: span.Parent().SpanID().String(),
		Duration: span.EndTime().Sub(span.StartTime()).String(),
		Attrs:    map[string]string{},
	}
	if span.Status().Code == codes.Error {
		data.Error = span.Status().Description
	}
	for _, attr := range span.Attributes() {
		data.Attrs[string(attr.Key)] = attr.Value.AsString()
	}
	return data
}
