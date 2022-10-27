package json

import log "github.com/livebud/bud/package/log2"

func New() *Handler {
	return &Handler{}
}

type Handler struct {
}

func (h *Handler) Log(entry *log.Entry) error {
	return nil
}
