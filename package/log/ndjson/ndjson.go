package json

import log "github.com/livebud/bud/package/log"

func New() *Handler {
	return &Handler{}
}

type Handler struct {
}

func (h *Handler) Log(entry *log.Entry) error {
	return nil
}
