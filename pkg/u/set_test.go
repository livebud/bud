package u_test

import (
	"testing"

	"github.com/livebud/bud/pkg/u"
	"github.com/matryer/is"
)

func TestOrderedStringSet(t *testing.T) {
	is := is.New(t)
	is.Equal(
		u.Set([]string{"a", "c", "b", "b", "d", "a", "c", "d", "e"}...),
		[]string{"a", "c", "b", "d", "e"},
	)
	is.Equal(
		u.Set([]string{"a", "a", "a"}...),
		[]string{"a"},
	)
}
