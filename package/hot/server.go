package hot

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/log"
)

// New server-sent event (SSE) server
func New(log log.Interface, ps pubsub.Subscriber) *Server {
	return &Server{log, ps, time.Now}
}

type Server struct {
	log log.Interface
	ps  pubsub.Subscriber
	Now func() time.Time // Used for testing
}

func pagePath(url string) string {
	return strings.TrimPrefix(strings.TrimPrefix(url, "/bud/hot"), "/")
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
	topics := []string{"frontend:update"}
	pagePath := pagePath(r.URL.Path)
	if pagePath != "" {
		topics = append(topics, `frontend:update:`+pagePath)
	}
	subscription := s.ps.Subscribe(topics...)
	s.log.Debug("hot: subscribed to topics", "topics", topics)
	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-subscription.Wait():
			s.log.Debug("hot: got event", "topic", "frontend:update")
			if pagePath == "" {
				s.log.Debug("hot: no page path, triggering a full reload")
				reload(flusher, w)
				continue
			}
			// Add /bud/ because we'll be requesting a generated file
			scriptPath := fmt.Sprintf("%s?ts=%d", "/bud/"+pagePath, s.Now().UnixMilli())
			event := &Event{
				Data: []byte(fmt.Sprintf(`{"scripts":[%q]}`, scriptPath)),
			}
			w.Write(event.Format().Bytes())
			flusher.Flush()

		// TODO: Create a new event type. EventSourcing has a concept of event types
		// which can be differentiated by the browser.
		//
		// See: https://html.spec.whatwg.org/multipage/server-sent-events.html#server-sent-events-intro
		case <-s.ps.Subscribe("backend:update").Wait():
			s.log.Debug("hot: got event", "topic", "page:reload")
			reload(flusher, w)
		}
	}
}

func reload(flusher http.Flusher, w http.ResponseWriter) {
	event := &Event{
		Data: []byte(`{"reload":true}`),
	}
	w.Write(event.Format().Bytes())
	flusher.Flush()
}

// https://html.spec.whatwg.org/multipage/server-sent-events.html#event-stream-interpretation
type Event struct {
	ID    string // id (optional)
	Type  string // event type (optional)
	Data  []byte // data
	Retry int    // retry (optional)
}

func (e *Event) Format() *bytes.Buffer {
	b := new(bytes.Buffer)
	if e.ID != "" {
		b.WriteString("id: " + e.ID + "\n")
	}
	if e.Type != "" {
		b.WriteString("event: " + e.Type + "\n")
	}
	if len(e.Data) > 0 {
		b.WriteString("data: ")
		b.Write(e.Data)
		b.WriteByte('\n')
	}
	if e.Retry > 0 {
		b.WriteString("retry: " + strconv.Itoa(e.Retry) + "\n")
	}
	b.WriteByte('\n')
	return b
}
