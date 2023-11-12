package glob

import "github.com/livebud/bud/pkg/u"

// Bases returns all non-magical parts of the glob
func Bases(pattern string) ([]string, error) {
	expands, err := Expand(pattern)
	if err != nil {
		return nil, err
	}
	bases := make([]string, len(expands))
	for i, expand := range expands {
		bases[i] = Base(expand)
	}
	return u.Set(bases...), nil
}
