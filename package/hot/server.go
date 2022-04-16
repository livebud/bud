package hot

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"gitlab.com/mnm/bud/runtime/web"

	"gitlab.com/mnm/bud/internal/pubsub"
)

func New() *Server {
	return &Server{pubsub.New()}
}

type Server struct {
	ps pubsub.Client
}

func (s *Server) Reload(path string) {
	s.ps.Publish(path, nil)
}

// Start listening on addr
func (s *Server) ListenAndServe(ctx context.Context, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return web.Serve(ctx, listener, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Take control of flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		err := fmt.Errorf("hot: response writer is not a flusher")
		http.Error(w, err.Error(), 500)
		return
	}
	// Set the appropriate response headers
	headers := w.Header()
	headers.Add(`Content-Type`, `text/event-stream`)
	headers.Add(`Cache-Control`, `no-cache`)
	headers.Add(`Connection`, `keep-alive`)
	headers.Add(`Access-Control-Allow-Origin`, "*")
	// Flush the headers
	flusher.Flush()
	// Subscribe to a specific page path or all pages
	pagePath := r.URL.Query().Get("page")
	topics := []string{"*"}
	if pagePath != "" {
		topics = append(topics, pagePath[1:])
	}
	subscription := s.ps.Subscribe(topics...)
	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case topic := <-subscription.Wait():
			_ = topic
			payload := fmt.Sprintf("data: {\"scripts\":[%q]}\n\n", fmt.Sprintf("%s?ts=%d", pagePath, time.Now().UnixMilli()))
			w.Write([]byte(payload))
			flusher.Flush()

		// TODO: rethink this, this was just the easiest way I could think of to add
		// full-reloading. Using the exclamation point so it doesn't conflict with a
		// file path
		case <-s.ps.Subscribe("!").Wait():
			payload := fmt.Sprintf("data: {\"reload\":true}\n\n")
			w.Write([]byte(payload))
			flusher.Flush()
		}
	}
}
