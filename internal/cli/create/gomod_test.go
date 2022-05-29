package create

import (
	"testing"

	"github.com/livebud/bud/internal/is"
)

func TestGoVersion(t *testing.T) {
	// TODO: does test if test dependencies end up in the final binary if we
	// aren't using the package *_test
	is := is.New(t)
	is.Equal(goVersion("1.18.1"), "1.18")
	is.Equal(goVersion("1.18"), "1.18")
	is.Equal(goVersion("1"), "1.0")
}
