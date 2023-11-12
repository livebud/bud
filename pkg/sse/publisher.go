package sse

import (
	"context"
	"fmt"
	"net/http"
)

type Publisher interface {
	Publish(ctx context.Context, event *Event) error
}

// Create a sender from a response writer
func Create(w http.ResponseWriter) (Publisher, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("sse: response writer is not a flusher")
	}
	// Set the appropriate response headers
	headers := w.Header()
	headers.Add(`Content-Type`, `text/event-stream`)
	headers.Add(`Cache-Control`, `no-cache`)
	headers.Add(`Connection`, `keep-alive`)
	headers.Add(`Access-Control-Allow-Origin`, "*")
	w.WriteHeader(http.StatusOK)
	// Flush the headers
	flusher.Flush()
	return &publisher{w, flusher}, nil
}

type publisher struct {
	w http.ResponseWriter
	f http.Flusher
}

var _ Publisher = (*publisher)(nil)

func (p *publisher) Publish(ctx context.Context, event *Event) error {
	if _, err := p.w.Write(event.Format().Bytes()); err != nil {
		return fmt.Errorf("sse: unable to publish event %s: %w", event, err)
	}
	p.f.Flush()
	return nil
}
