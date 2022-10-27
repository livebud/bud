package filter

import (
	log "github.com/livebud/bud/package/log"
)

func New(handler log.Handler, level log.Level) *Handler {
	return &Handler{handler, level}
}

// Handler logs by level
type Handler struct {
	handler log.Handler
	level   log.Level
}

func (f *Handler) Log(entry *log.Entry) error {
	if entry.Level < f.level {
		return nil
	}
	return f.handler.Log(entry)
}
