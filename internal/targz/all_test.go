package targz

import (
	"io/fs"
	"testing"
	"testing/fstest"
	"time"

	"github.com/livebud/bud/internal/is"
)

var modTime = time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)

func TestZipUnzip(t *testing.T) {
	is := is.New(t)
	gzip, err := Zip(fstest.MapFS{
		"a.go":     &fstest.MapFile{Data: []byte("package a")},
		"b/c/d.go": &fstest.MapFile{Data: []byte("package d"), Mode: 0644, ModTime: modTime},
		"e":        &fstest.MapFile{Mode: fs.ModeDir, ModTime: modTime},
	})
	is.NoErr(err)
	is.True(len(gzip) != 0)
	fsys, err := Unzip(gzip)
	is.NoErr(err)
	data, err := fs.ReadFile(fsys, "a.go")
	is.NoErr(err)
	is.Equal(string(data), "package a")
	fi, err := fs.Stat(fsys, "b/c/d.go")
	is.NoErr(err)
	is.Equal(fi.Mode(), fs.FileMode(0644))
	is.True(fi.ModTime().Equal(modTime))
	data, err = fs.ReadFile(fsys, "b/c/d.go")
	is.NoErr(err)
	is.Equal(string(data), "package d")
	fi, err = fs.Stat(fsys, "e")
	is.NoErr(err)
	is.Equal(fi.Mode(), fs.ModeDir)
	is.True(fi.ModTime().Equal(modTime))
}
