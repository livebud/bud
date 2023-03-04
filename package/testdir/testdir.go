package testdir

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/virtual"
	"golang.org/x/sync/errgroup"
)

// Directory to copy into, if empty, we use $TMPDIR/bud_testdir_*.
var dirFlag = flag.String("dir", "", "dir to copy into")

// MakeTemp creates and returns the temporary directory for testing. It only
// creates the directory once, but can be called multiple times to return the
// directory.
var MakeTemp = once.String(func() (string, error) {
	dir, err := makeTemp(*dirFlag)
	if err != nil {
		return "", err
	}
	return filepath.Clean(dir), nil
})

func makeTemp(dir string) (string, error) {
	if dir != "" {
		if !strings.HasPrefix(dir, "_tmp") {
			return "", fmt.Errorf("testdir: dir path must be relative and start with _tmp")
		}
		if strings.Contains(dir, "*") {
			return os.MkdirTemp(".", dir)
		}
		return dir, os.MkdirAll(dir, 0755)
	}
	return os.MkdirTemp("", "bud_testdir_*")
}

// Copy files into a directory.
func Copy(fsys fs.FS) error {
	log := testlog.New()
	dir, err := MakeTemp()
	if err != nil {
		return err
	}
	log.Debug("testdir: copying into", dir)
	dirfs := virtual.OS(dir)
	return virtual.Copy(log, fsys, dirfs)
}

func Sync(fsys fs.FS) error {
	log := testlog.New()
	dir, err := MakeTemp()
	if err != nil {
		return err
	}
	log.Debug("testdir: syncing into", dir)
	dirfs := virtual.OS(dir)
	return virtual.Sync(log, fsys, dirfs)
}

// Exists returns nil if path exists. Should be called after Write.
func Exists(paths ...string) error {
	eg := new(errgroup.Group)
	for _, path := range paths {
		path := path
		eg.Go(func() error {
			return exist(path)
		})
	}
	return eg.Wait()
}

func exist(path string) error {
	if _, err := Stat(path); err != nil {
		return err
	}
	return nil
}

// NotExists returns nil if path doesn't exist. Should be called after Write.
func NotExists(paths ...string) error {
	eg := new(errgroup.Group)
	for _, path := range paths {
		path := path
		eg.Go(func() error {
			return notExist(path)
		})
	}
	return eg.Wait()
}

func notExist(path string) error {
	if _, err := Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	return fmt.Errorf("%s exists: %w", path, fs.ErrExist)
}

// Stat return the file info of path. Should be called after Write.
func Stat(path string) (fs.FileInfo, error) {
	dir, err := MakeTemp()
	if err != nil {
		return nil, err
	}
	return os.Stat(filepath.Join(dir, path))
}
