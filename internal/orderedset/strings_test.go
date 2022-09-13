package orderedset_test

import (
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/orderedset"
)

func TestOrderedStringSet(t *testing.T) {
	is := is.New(t)
	is.Equal(orderedset.Strings(
		[]string{"a", "c", "b", "b", "d", "a", "c", "d", "e"}...),
		[]string{"a", "c", "b", "d", "e"},
	)
	is.Equal(orderedset.Strings(
		[]string{"a", "a", "a"}...),
		[]string{"a"},
	)
	is.Equal(orderedset.Strings(), []string{})
}

func TestNewSet(t *testing.T) {
	is := is.New(t)
	// Empty
	s := orderedset.New[string]()
	is.Equal(s.List(), []string{})
	// With duplicates
	s = orderedset.New[string]()
	s.Add("a", "a", "a")
	is.Equal(s.List(), []string{"a"})
	// With others
	s = orderedset.New[string]()
	s.Add("a", "c", "b", "b", "d", "a", "c", "d", "e")
	is.Equal(s.List(), []string{"a", "c", "b", "d", "e"})
}
