package snapshot_test

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"
	"time"

	"github.com/livebud/bud/internal/snapshot"

	"github.com/livebud/bud/internal/is"
)

var modTime = time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)
var modTime2 = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

func TestHash(t *testing.T) {
	is := is.New(t)
	mapfs := fstest.MapFS{}
	mapfs["main.go"] = &fstest.MapFile{Data: []byte(`package main`)}
	h1, err := snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(len(h1), 11)
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0644, ModTime: modTime}
	h2, err := snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(len(h2), 11)
	is.True(h1 != h2)
	// New dir doesn't change anything
	mapfs["e"] = &fstest.MapFile{Data: []byte(``), Mode: fs.ModeDir, ModTime: modTime}
	h3, err := snapshot.Hash(mapfs)
	is.NoErr(err)
	is.True(h2 == h3)
	// Adjust data
	mapfs["main.go"] = &fstest.MapFile{Data: []byte(`package main; func main() {}`)}
	h4, err := snapshot.Hash(mapfs)
	is.NoErr(err)
	is.True(h3 != h4)
	// Adjust mode doesn't change anything
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0655, ModTime: modTime}
	h5, err := snapshot.Hash(mapfs)
	is.NoErr(err)
	is.True(h4 == h5)
	// Adjust modtime, shouldn't change anything
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0655, ModTime: modTime2}
	h6, err := snapshot.Hash(mapfs)
	is.NoErr(err)
	is.True(h5 == h6)
	// Hash with no changes
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0655, ModTime: modTime2}
	h7, err := snapshot.Hash(mapfs)
	is.NoErr(err)
	is.True(h6 == h7)
}

func TestBackupRestore(t *testing.T) {
	is := is.New(t)
	original := fstest.MapFS{}
	current := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte(`package main`)},
	}
	h1, err := snapshot.Hash(original)
	is.NoErr(err)
	is.Equal(len(h1), 11)
	// Backup based on the hash
	err = snapshot.Backup(h1, current)
	is.NoErr(err)
	fsys, err := snapshot.Restore(h1)
	is.NoErr(err)
	data, err := fs.ReadFile(fsys, "main.go")
	is.NoErr(err)
	is.Equal(string(data), "package main")
	// Restore again
	fsys, err = snapshot.Restore(h1)
	is.NoErr(err)
	data, err = fs.ReadFile(fsys, "main.go")
	is.NoErr(err)
	is.Equal(string(data), "package main")
}

func TestRestoreNotExist(t *testing.T) {
	is := is.New(t)
	original := fstest.MapFS{
		"bin.go": &fstest.MapFile{Data: []byte(`package bin`)},
	}
	h1, err := snapshot.Hash(original)
	is.NoErr(err)
	is.Equal(len(h1), 11)
	// Restore non-existent
	fsys, err := snapshot.Restore(h1)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(fsys, nil)
}

// TODO: test changed files within replaced modules
// triggering hash changes
