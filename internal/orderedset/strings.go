package orderedset

func Strings(list ...string) []string {
	seen := map[string]struct{}{}
	set := make([]string, 0, len(list))
	for _, item := range list {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		set = append(set, item)
	}
	return set
}

func New[T comparable]() *Set[T] {
	return &Set[T]{
		seen: map[T]struct{}{},
		set:  []T{},
	}
}

type Set[T comparable] struct {
	seen map[T]struct{}
	set  []T
}

func (s *Set[T]) Add(items ...T) {
	for _, item := range items {
		s.add(item)
	}
}

func (s *Set[T]) add(item T) {
	if _, ok := s.seen[item]; ok {
		return
	}
	s.seen[item] = struct{}{}
	s.set = append(s.set, item)
}

func (s *Set[T]) Has(item T) bool {
	_, ok := s.seen[item]
	return ok
}

func (s *Set[T]) List() []T {
	return s.set
}
