package u

func Set[T comparable](list ...T) []T {
	seen := map[T]struct{}{}
	set := make([]T, 0, len(list))
	for _, item := range list {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		set = append(set, item)
	}
	return set
}
