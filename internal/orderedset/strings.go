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
