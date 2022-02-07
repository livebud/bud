package snapshot_test

import (
	"io/fs"
	"testing"
	"testing/fstest"
	"time"

	"gitlab.com/mnm/bud/internal/snapshot"

	"github.com/matryer/is"
)

var modTime = time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)
var modTime2 = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

func TestHash(t *testing.T) {
	is := is.New(t)
	mapfs := fstest.MapFS{}
	mapfs["main.go"] = &fstest.MapFile{Data: []byte(`package main`)}
	key, err := snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `UPcEYUuxm78`)
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0644, ModTime: modTime}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `CWmhz9qgIFo`)
	mapfs["e"] = &fstest.MapFile{Mode: fs.ModeDir, ModTime: modTime}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `XQ4COBOPHtE`)
	// Adjust data
	mapfs["main.go"] = &fstest.MapFile{Data: []byte(`package main; func main() {}`)}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `xEpUlkb9G0E`)
	// Adjust mode
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0655, ModTime: modTime}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `86nQ-5k1RP4`)
	// Adjust modtime
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0655, ModTime: modTime2}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `ixro1htVtx8`)
	// Hash with nothing changing
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0655, ModTime: modTime2}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `ixro1htVtx8`)
}

func TestBackupRestore(t *testing.T) {
	is := is.New(t)
	original := fstest.MapFS{}
	current := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte(`package main`)},
	}
	hash, err := snapshot.Hash(original)
	is.NoErr(err)
	is.Equal(hash, "SFQ0mar4vsA")
	// Backup based on the hash
	err = snapshot.Backup(hash, current)
	is.NoErr(err)
	fsys, err := snapshot.Restore(hash)
	is.NoErr(err)
	data, err := fs.ReadFile(fsys, "main.go")
	is.NoErr(err)
	is.Equal(string(data), "package main")
	// Restore again
	fsys, err = snapshot.Restore(hash)
	is.NoErr(err)
	data, err = fs.ReadFile(fsys, "main.go")
	is.NoErr(err)
	is.Equal(string(data), "package main")
}
