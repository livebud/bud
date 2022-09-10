package linkmap

func New() *Map {
	return &Map{
		m: map[string]*List{},
	}
}

type Map struct {
	m map[string]*List
}

func (m *Map) Map() map[string]func(path string) bool {
	out := map[string]func(path string) bool{}
	for path, list := range m.m {
		out[path] = list.merge()
	}
	return out
}

func (m *Map) Scope(path string) *List {
	list := &List{}
	m.m[path] = list
	return list
}

type List struct {
	fns []func(path string) bool
}

func (l *List) Add(fn func(path string) bool) {
	l.fns = append(l.fns, fn)
}

func (l *List) merge() func(path string) bool {
	return func(path string) bool {
		for _, fn := range l.fns {
			if fn(path) {
				return true
			}
		}
		return false
	}
}
