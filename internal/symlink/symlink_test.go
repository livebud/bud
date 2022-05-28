package symlink_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/symlink"
)

func TestLink(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	from := filepath.Join(dir, "a.txt")
	err := os.WriteFile(from, []byte("a"), 0644)
	is.NoErr(err)
	to := filepath.Join(dir, "b.txt")
	err = symlink.Link(from, to)
	is.NoErr(err)
	data, err := os.ReadFile(to)
	is.NoErr(err)
	is.Equal(string(data), "a")
	// Test the override
	err = symlink.Link(from, to)
	is.NoErr(err)
	data, err = os.ReadFile(to)
	is.NoErr(err)
	is.Equal(string(data), "a")

}
