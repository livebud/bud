package memory

import log "github.com/livebud/bud/package/log"

func New() *Handler {
	return &Handler{}
}

type Handler struct {
	Entries []*log.Entry
}

var _ log.Handler = (*Handler)(nil)

func (h *Handler) Log(entry *log.Entry) error {
	h.Entries = append(h.Entries, entry)
	return nil
}
