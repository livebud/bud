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
