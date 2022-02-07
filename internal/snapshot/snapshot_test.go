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
	is.Equal(key, `UPcEYUuxm78`)
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0644, ModTime: modTime}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `du1Yyvmk_Ks`)
	mapfs["e"] = &fstest.MapFile{Mode: fs.ModeDir, ModTime: modTime}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `QY1LgL2TFbE`)
	// Adjust data
	mapfs["main.go"] = &fstest.MapFile{Data: []byte(`package main; func main() {}`)}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `6Y0qe6ntDqs`)
	// Adjust mode
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0655, ModTime: modTime}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `k12v200Bmu4`)
	// Adjust modtime, shouldn't change anything
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0655, ModTime: modTime2}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `k12v200Bmu4`)
	// Hash with nothing changing
	mapfs["b/c/d.go"] = &fstest.MapFile{Data: []byte("package d"), Mode: 0655, ModTime: modTime2}
	key, err = snapshot.Hash(mapfs)
	is.NoErr(err)
	is.Equal(key, `k12v200Bmu4`)
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

func TestRestoreNotExist(t *testing.T) {
	is := is.New(t)
	original := fstest.MapFS{
		"bin.go": &fstest.MapFile{Data: []byte(`package bin`)},
	}
	hash, err := snapshot.Hash(original)
	is.NoErr(err)
	is.Equal(hash, "wAkDzu4jU2g")
	// Restore non-existent
	fsys, err := snapshot.Restore(hash)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(fsys, nil)
}
