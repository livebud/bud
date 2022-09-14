package linkmap

import "sync"

// New linkmap. Linkmap is safe for concurrent use.
func New() *Map {
	return &Map{}
}

type Map struct {
	sm sync.Map
}

func (m *Map) Get(path string) (*List, bool) {
	list, ok := m.sm.Load(path)
	if !ok {
		return nil, false
	}
	return list.(*List), true
}

func (m *Map) Scope(path string) *List {
	list := &List{}
	m.sm.Store(path, list)
	return list
}

func (m *Map) Range(fn func(path string, list *List) bool) {
	m.sm.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*List))
	})
}

type List struct {
	mu  sync.RWMutex
	fns []func(path string) bool
}

func (l *List) Add(fn func(path string) bool) {
	l.mu.Lock()
	l.fns = append(l.fns, fn)
	l.mu.Unlock()
}

func (l *List) Check(path string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, fn := range l.fns {
		if fn(path) {
			return true
		}
	}
	return false
}
