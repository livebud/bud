package budfs_test

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/log"
)

type tailwind struct {
}

func (t *tailwind) GenerateFile(fsys budfs.FS, file *budfs.File) error {
	file.Data = []byte("/* tailwind */")
	return nil
}

type svelte struct {
}

func (s *svelte) GenerateFile(fsys budfs.FS, file *budfs.File) error {
	file.Data = []byte("/* svelte */")
	return nil
}

func TestFS(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"bud/public/index.html": &fstest.MapFile{Data: []byte("<h1>hello</h1>")},
	}
	log := log.Discard
	budfs := budfs.New(fsys, log)
	budfs.FileGenerator("bud/public/tailwind/tailwind.css", &tailwind{})
	budfs.FileGenerator("bud/view/index.svelte", &svelte{})

	// .
	des, err := fs.ReadDir(budfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)

	// bud
	is.Equal(des[0].Name(), "bud")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)
	fi, err := des[0].Info()
	is.NoErr(err)
	is.Equal(fi.Mode(), fs.ModeDir)
	stat, err := fs.Stat(budfs, "bud")
	is.NoErr(err)
	is.Equal(stat.Mode(), fs.ModeDir)

	// bud/public
	des, err = fs.ReadDir(budfs, "bud")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "public")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)
	fi, err = des[0].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "public")
	stat, err = fs.Stat(budfs, "bud/public")
	is.NoErr(err)
	is.Equal(stat.Name(), "public")

	// return errors for non-existent files
	_, err = budfs.Open("bud\\public")
	is.True(errors.Is(err, fs.ErrNotExist))

	// bud/public/tailwind
	des, err = fs.ReadDir(budfs, "bud/public/tailwind")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "tailwind.css")
	is.Equal(des[0].IsDir(), false)

	// read generated data
	data, err := fs.ReadFile(budfs, "bud/public/index.html")
	is.NoErr(err)
	is.Equal(string(data), "<h1>hello</h1>")
	data, err = fs.ReadFile(budfs, "bud/public/tailwind/tailwind.css")
	is.NoErr(err)
	is.Equal(string(data), "/* tailwind */")
	data, err = fs.ReadFile(budfs, "bud/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(data), "/* svelte */")

	// run the TestFS compliance test suite
	is.NoErr(fstest.TestFS(budfs, "bud/public/index.html", "bud/public/tailwind/tailwind.css", "bud/view/index.svelte"))
}
