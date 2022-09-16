package virtual_test

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/virtual"
)

func TestExclude(t *testing.T) {
	is := is.New(t)
	tree := virtual.Tree{
		"view/a.txt": &virtual.File{Data: []byte("a")},
		"view/b.txt": &virtual.File{Data: []byte("b")},
		"bud/bud.go": &virtual.File{Data: []byte("bud")},
	}
	fsys := virtual.Exclude(tree, func(path string) bool {
		return path == "bud" || strings.HasPrefix(path, "bud/")
	})
	des, err := fs.ReadDir(fsys, ".")
	is.Equal(err, nil)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "view")
}
