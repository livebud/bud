package filter

import (
	"github.com/livebud/bud/package/log"
)

// Load console handler
func Load(handler log.Handler, pattern string) (log.Handler, error) {
	level, err := log.ParseLevel(pattern)
	if err != nil {
		return nil, err
	}
	return &Filter{handler, level}, nil
}

// Filter logs by level. Can be initialized manually or by the Load function.
type Filter struct {
	Handler log.Handler
	Level   log.Level
}

func (f *Filter) Log(entry log.Entry) {
	if entry.Level < f.Level {
		return
	}
	f.Handler.Log(entry)
}
