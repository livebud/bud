package hot

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"gitlab.com/mnm/bud/internal/pubsub"
)

func New() *Server {
	ps := pubsub.New()
	return &Server{
		pubsub: ps,
		server: &http.Server{
			Handler: handler(ps),
		},
	}
}

type Server struct {
	server *http.Server
	pubsub pubsub.Client
}

func (s *Server) Reload(path string) {
	s.pubsub.Publish(path, nil)
}

// Support shutting down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Serve(l net.Listener) error {
	return s.server.Serve(l)
}

// type pubsub struct {
// 	subs map[string][]func()
// }

// func (p *pubsub) Subscribe(ctx context.Context, path string, fn func()) {
// 	p.subs[path] = append(p.subs[path], fn)
// 	p.subs["*"] = append(p.subs["*"], fn)
// 	// Wait for the request to finish
// 	<-ctx.Done()
// 	// TODO: unsubscribe the subscribers
// }

// // Publish to reload
// func (p *pubsub) Publish(path string) {
// 	if path == "*" {
// 		for _, fns := range p.subs {
// 			for _, fn := range fns {
// 				fn()
// 			}
// 		}
// 		return
// 	}
// 	fns, ok := p.subs[path]
// 	if !ok {
// 		return
// 	}
// 	for _, fn := range fns {
// 		fn()
// 	}
// }

func handler(ps pubsub.Subscriber) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			err := fmt.Errorf("hot: response writer is not a flusher")
			fmt.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		// Set the appropriate response headers
		headers := w.Header()
		headers.Add(`Content-Type`, `text/event-stream`)
		headers.Add(`Cache-Control`, `no-cache`)
		headers.Add(`Connection`, `keep-alive`)
		headers.Add(`Access-Control-Allow-Origin`, `*`)
		// Listen for this page path
		pagePath := r.URL.Query().Get("page")
		if pagePath == "" {
			err := fmt.Errorf("hot: missing page query")
			http.Error(w, err.Error(), 500)
			return
		}
		// Flush the headers
		flusher.Flush()
		fmt.Println("subscribed to", pagePath[1:])
		// Start the subscription
		subscription := ps.Subscribe("*", pagePath[1:])
		ctx := r.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case <-subscription.Wait():
				fmt.Println("reloading", pagePath)
				payload := fmt.Sprintf("data: {\"scripts\":[%q]}\n\n", fmt.Sprintf("%s?ts=%d", pagePath, time.Now().UnixMilli()))
				w.Write([]byte(payload))
				flusher.Flush()
			}
		}
	})
}
