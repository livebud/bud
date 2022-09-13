package linkmap

type Map map[string]*List

func (m Map) Scope(path string) *List {
	list := &List{}
	m[path] = list
	return list
}

type List []func(path string) bool

func (list *List) Add(fn func(path string) bool) {
	*list = append(*list, fn)
}

func (list List) Check(path string) bool {
	for _, fn := range list {
		if fn(path) {
			return true
		}
	}
	return false
}
