package glob_test

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/glob"
	"github.com/livebud/bud/internal/is"
)

func TestWalk(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"controller/controller.go":       &fstest.MapFile{Data: []byte(`// /`), Mode: 0644},
		"controller/posts/controller.go": &fstest.MapFile{Data: []byte(`// /posts`), Mode: 0644},
	}
	matches := []string{}
	err := glob.Walk(fsys, "controller/**.go", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		matches = append(matches, path)
		return nil
	})
	is.NoErr(err)
	is.Equal(len(matches), 2)
	is.Equal(matches[0], "controller/controller.go")
	is.Equal(matches[1], "controller/posts/controller.go")
}
