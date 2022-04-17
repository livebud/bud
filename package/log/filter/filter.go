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
	return &filter{handler, level}, nil
}

type filter struct {
	handler log.Handler
	level   log.Level
}

func (f *filter) Log(entry log.Entry) {
	if entry.Level < f.level {
		return
	}
	f.handler.Log(entry)
}
