package logfmt

import (
	"io"
	"sync"

	"github.com/go-logfmt/logfmt"
	"github.com/livebud/bud/pkg/logs"
)

func New(w io.Writer) *Handler {
	return &Handler{enc: logfmt.NewEncoder(w)}
}

type Handler struct {
	mu  sync.Mutex
	enc *logfmt.Encoder
}

func (h *Handler) Log(entry *logs.Entry) error {
	keys := entry.Fields.Keys()
	h.mu.Lock()
	defer h.mu.Unlock()
	h.enc.EncodeKeyval("ts", entry.Time)
	h.enc.EncodeKeyval("lvl", entry.Level.String())
	h.enc.EncodeKeyval("msg", entry.Message)
	for _, key := range keys {
		h.enc.EncodeKeyval(key, entry.Fields.Get(key))
	}
	h.enc.EndRecord()
	return nil
}
