package sse

import (
	"context"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/livebud/bud/pkg/locals"
	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/request"
)

// New server-sent event (SSE) handler
func New(log logs.Log) *Handler {
	return &Handler{
		pub: newPublishers(log),
		log: log,
	}
}

type Handler struct {
	id  atomic.Int64
	pub *publishers
	log logs.Log
}

var _ http.Handler = (*Handler)(nil)
var _ Publisher = (*Handler)(nil)

func (s *Handler) Publish(ctx context.Context, event *Event) error {
	return s.pub.Publish(ctx, event)
}

const defaultUrl = "/live.js"

// Middleware for serving live reload events and passing the live url along
// TODO: middleware should support options with configurable url
func (s *Handler) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If the request is for the live route, serve it
			if r.Method == http.MethodGet && r.URL.Path == defaultUrl {
				s.ServeHTTP(w, r)
				return
			}
			// Otherwise, pass the live url along
			ctx := locals.Set(r.Context(), "live_url", defaultUrl)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch request.Accept(r, "text/javascript", "text/event-stream") {
	case "text/javascript":
		s.serveJS(w, r)
	case "text/event-stream":
		s.acceptClient(w, r)
	default:
		http.Error(w, "sse: unsupportedaccept header", http.StatusNotAcceptable)
		return
	}
}

func (s *Handler) serveJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	w.Write([]byte(`(function() {
		var source = new EventSource("` + r.URL.String() + `");
		source.addEventListener("message", function(e) {
			const { created, updated, deleted } = JSON.parse(e.data)
			if (!created.length && !deleted.length && updated.length === 1) {
				for (const link of document.getElementsByTagName("link")) {
					const url = new URL(link.href)
					const pathname = url.pathname.replace(/^\//, '')
					if (url.host === location.host && pathname === updated[0]) {
						const next = link.cloneNode()
						next.href = updated[0] + '?' + Math.random().toString(36).slice(2)
						next.onload = () => link.remove()
						link.parentNode.insertBefore(next, link.nextSibling)
						return
					}
				}
			}
			window.location.reload()
		});
	})();`))
}

func (s *Handler) acceptClient(w http.ResponseWriter, r *http.Request) {
	publisher, err := Create(w)
	if err != nil {
		s.log.Errorf("sse: unable to create publisher: %w", err)
		http.Error(w, err.Error(), 500)
		return
	}
	// Add the client to the publisher
	clientID := strconv.FormatInt(s.id.Add(1), 10)
	eventCh := s.pub.Set(clientID, publisher)
	defer s.pub.Remove(clientID)
	// Wait for the client to disconnect
	ctx := r.Context()
	for {
		select {
		// Send events to the client
		case event := <-eventCh:
			publisher.Publish(ctx, event)
		// Client disconnected
		case <-ctx.Done():
			return
		}
	}
}
