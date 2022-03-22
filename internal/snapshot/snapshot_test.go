package snapshot_test

import (
	"errors"
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
	is.Equal(key, `Fiyf9IKN3Y0`)
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0644, ModTime: modTime}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `kdsFZWQqPcc`)
	// New dir doesn't change anything
	mapfs["e"] = &fstest.MapFile{Data: []byte(``), Mode: fs.ModeDir, ModTime: modTime}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `kdsFZWQqPcc`)
	// Adjust data
	mapfs["main.go"] = &fstest.MapFile{Data: []byte(`package main; func main() {}`)}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `5qEzdxZ33Ws`)
	// Adjust mode doesn't change anything
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0655, ModTime: modTime}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `5qEzdxZ33Ws`)
	// Adjust modtime, shouldn't change anything
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0655, ModTime: modTime2}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `5qEzdxZ33Ws`)
	// Hash with no changes
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0655, ModTime: modTime2}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `5qEzdxZ33Ws`)
}

func TestBackupRestore(t *testing.T) {
	is := is.New(t)
	original := fstest.MapFS{}
	current := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte(`package main`)},
	}
	hash, err := snapshot.Hash(original)
	is.NoErr(err)
	is.Equal(hash, "70bbN1HY6Zk")
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

func TestRestoreNotExist(t *testing.T) {
	is := is.New(t)
	original := fstest.MapFS{
		"bin.go": &fstest.MapFile{Data: []byte(`package bin`)},
	}
	hash, err := snapshot.Hash(original)
	is.NoErr(err)
	is.Equal(hash, "FD4w4ZnukkU")
	// Restore non-existent
	fsys, err := snapshot.Restore(hash)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(fsys, nil)
}
