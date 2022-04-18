package logfilter

import (
	"errors"
	"fmt"

	"github.com/livebud/bud/package/log"
)

//go:generate peg -switch -inline filter.peg

var ErrParsing = errors.New("urlx: unable to parse")

type filters []filter

func (filters filters) Match(entry log.Entry) bool {
	for _, filter := range filters {
		if !filter.Match(entry) {
			return false
		}
	}
	return true
}

// Used in the parser
type filter struct {
	level    string
	packages []string
}

func (f *filter) Match(entry log.Entry) bool {
	fmt.Println("matching", entry.Level)
	return false
}

type Matcher interface {
	Match(entry log.Entry) bool
}

func Parse(pattern string) (Matcher, error) {
	parser := &parser{Buffer: pattern}
	parser.Init()
	err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("%w %q", ErrParsing, pattern)
	}
	parser.Execute()
	return parser.filters, nil
}
