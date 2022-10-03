package gitignore_test

import (
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/gitignore"
	"github.com/livebud/bud/internal/is"
)

func TestBudRoot(t *testing.T) {
	is := is.New(t)
	ignore := gitignore.FromFS(fstest.MapFS{
		".gitignore": &fstest.MapFile{Data: []byte(`/bud`)},
	})
	is.True(ignore("bud/internal/web/web.go"))
	is.True(!ignore("main.go"))
}

func TestGitDir(t *testing.T) {
	is := is.New(t)
	ignore := gitignore.FromFS(fstest.MapFS{
		".gitignore": &fstest.MapFile{Data: []byte(``)},
	})
	is.True(ignore(".git"))
	is.True(ignore(".git/objects"))
}

func TestNodeModules(t *testing.T) {
	is := is.New(t)
	ignore := gitignore.FromFS(fstest.MapFS{
		".gitignore": &fstest.MapFile{Data: []byte(``)},
	})
	is.True(ignore("node_modules"))
	is.True(ignore("node_modules/svelte/internal/compiler.js"))
}
