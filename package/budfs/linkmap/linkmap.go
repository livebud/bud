package linkmap

import (
	"sync"

	"github.com/livebud/bud/package/log"
)

// New linkmap. Linkmap is safe for concurrent use.
func New(log log.Log) *Map {
	return &Map{log: log}
}

type Map struct {
	log log.Log
	sm  sync.Map
}

func (m *Map) Get(path string) (*List, bool) {
	list, ok := m.sm.Load(path)
	if !ok {
		return nil, false
	}
	return list.(*List), true
}

func (m *Map) Scope(path string) *List {
	list := &List{log: m.log, from: path, tos: map[string]struct{}{}}
	m.sm.Store(path, list)
	return list
}

func (m *Map) Range(fn func(path string, list *List) bool) {
	m.sm.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*List))
	})
}

type List struct {
	log  log.Log
	mu   sync.RWMutex
	from string
	fns  []func(path string) bool
	tos  map[string]struct{}
}

func (l *List) Link(caller string, tos ...string) {
	l.log.Fields(log.Fields{
		"caller": caller,
		"from":   l.from,
		"to":     tos,
	}).Debug("linkmap: link")
	l.mu.Lock()
	for _, to := range tos {
		l.tos[to] = struct{}{}
	}
	l.mu.Unlock()
}

func (l *List) Select(caller string, fn func(path string) bool) {
	l.log.Fields(log.Fields{
		"caller": caller,
		"from":   l.from,
	}).Debug("linkmap: select")
	l.mu.Lock()
	l.fns = append(l.fns, fn)
	l.mu.Unlock()
}

func (l *List) Check(path string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for to := range l.tos {
		if to == path {
			return true
		}
	}
	for _, fn := range l.fns {
		if fn(path) {
			return true
		}
	}
	return false
}
