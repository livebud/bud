package hot

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.com/mnm/bud/gen"
)

var hotPath = "/bud/hot"

func New(bf gen.FS) *Server {
	return &Server{bf}
}

type Server struct {
	bf gen.FS
}

// Middleware that handles refreshing the frontend
func (s *Server) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != hotPath {
			next.ServeHTTP(w, r)
			return
		}
		flusher, ok := w.(http.Flusher)
		if !ok {
			err := fmt.Errorf("hot: response writer is not a flusher")
			fmt.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		_ = flusher
		headers := w.Header()
		headers.Add(`Content-Type`, `text/event-stream`)
		headers.Add(`Cache-Control`, `no-cache`)
		headers.Add(`Connection`, `keep-alive`)

		pagePath := r.URL.Query().Get("page")
		if pagePath == "" {
			err := fmt.Errorf("hot: missing page query")
			http.Error(w, err.Error(), 500)
			return
		}
		sub, err := s.bf.Subscribe(pagePath[1:])
		if err != nil {
			err := fmt.Errorf("hot: unable to subscribe to updates: %w", err)
			http.Error(w, err.Error(), 500)
			return
		}
		defer sub.Close()
		// fmt.Println("subscribed to", pagePath[1:])

		ctx := r.Context()
		for {
			// Exit the loop if the request has gone away
			select {
			case <-ctx.Done():
				return
			case event := <-sub.Wait():
				fmt.Println("got new event", pagePath, string(event))
				payload := fmt.Sprintf("data: {\"scripts\":[%q]}\n\n", fmt.Sprintf("%s?ts=%d", pagePath, time.Now().UnixMilli()))
				w.Write([]byte(payload))
				flusher.Flush()
			}
		}
	})
}
