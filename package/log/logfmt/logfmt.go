package logfmt

import (
	"io"
	"sync"

	"github.com/go-logfmt/logfmt"
	log "github.com/livebud/bud/package/log"
)

func New(w io.Writer) *Handler {
	return &Handler{enc: logfmt.NewEncoder(w)}
}

type Handler struct {
	mu  sync.Mutex
	enc *logfmt.Encoder
}

func (h *Handler) Log(entry *log.Entry) error {
	keys := entry.Fields.Keys()
	h.mu.Lock()
	defer h.mu.Unlock()
	h.enc.EncodeKeyval("timestamp", entry.Timestamp)
	h.enc.EncodeKeyval("level", entry.Level.String())
	h.enc.EncodeKeyval("message", entry.Message)
	for _, key := range keys {
		h.enc.EncodeKeyval(key, entry.Fields.Get(key))
	}
	h.enc.EndRecord()
	return nil
}
