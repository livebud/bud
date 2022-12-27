package virt_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/virt"
)

func TestOSRead(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	is.NoErr(os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644))
	is.NoErr(os.MkdirAll(filepath.Join(dir, "c"), 0755))
	is.NoErr(os.WriteFile(filepath.Join(dir, "c/c.txt"), []byte("d"), 0644))
	// Try reading the directory
	fsys := virt.OS(dir)
	err := fstest.TestFS(fsys, "a.txt", "b.txt", "c/c.txt")
	is.NoErr(err)
}

func TestOSRemoveAllOutsideFail(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	is.NoErr(os.MkdirAll(filepath.Join(dir, "b"), 0755))
	fsys := virt.OS(filepath.Join(dir, "b"))
	err := fsys.RemoveAll("../a.txt")
	is.True(errors.Is(err, fs.ErrInvalid))
}

func TestOSRemoveAll(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	fsys := virt.OS(dir)
	err := fsys.RemoveAll("a.txt")
	is.NoErr(err)
}
