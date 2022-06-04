package current_test

import (
	"path/filepath"
	"testing"

	"github.com/livebud/bud/internal/current"
	"github.com/livebud/bud/internal/is"
)

func TestDir(t *testing.T) {
	is := is.New(t)
	dirname, err := current.Directory()
	is.NoErr(err)
	is.Equal(filepath.Base(dirname), "current")
}

func TestFile(t *testing.T) {
	is := is.New(t)
	filename, err := current.Filename()
	is.NoErr(err)
	is.Equal(filepath.Base(filename), "current_test.go")
}
