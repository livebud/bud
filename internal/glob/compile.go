package glob

import "github.com/gobwas/glob"

type Matcher = glob.Glob

// Comple
func Compile(pattern string) (Matcher, error) {
	// Expand patterns like {a,b} into multiple globs a & b. This avoids an
	// infinite loop described in this comment:
	// https://github.com/gobwas/glob/issues/50#issuecomment-1330182417
	patterns, err := Expand(pattern)
	if err != nil {
		return nil, err
	}
	globs := make(globs, len(patterns))
	for i, pattern := range patterns {
		glob, err := glob.Compile(pattern)
		if err != nil {
			return nil, err
		}
		globs[i] = glob
	}
	return globs, nil
}

type globs []glob.Glob

func (globs globs) Match(path string) bool {
	for _, glob := range globs {
		if glob.Match(path) {
			return true
		}
	}
	return false
}
