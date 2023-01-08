package versions_test

import (
	"errors"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/versions"
)

func TestGoVersion(t *testing.T) {
	is := is.New(t)
	is.NoErr(versions.CheckGo("go1.17"))
	is.NoErr(versions.CheckGo("go1.18"))
	is.True(errors.Is(versions.CheckGo("go1.16"), versions.ErrMinGoVersion))
	is.True(errors.Is(versions.CheckGo("go1.16.5"), versions.ErrMinGoVersion))
	is.True(errors.Is(versions.CheckGo("go1.8"), versions.ErrMinGoVersion))
	is.NoErr(versions.CheckGo("abc123"))
}
