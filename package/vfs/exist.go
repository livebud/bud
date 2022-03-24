package vfs

import (
	"errors"
	"io/fs"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Exist returns an error if any of the paths don't exist
func Exist(fsys fs.FS, paths ...string) (err error) {
	eg := new(errgroup.Group)
	for _, path := range paths {
		path := path
		eg.Go(func() error {
			if _, err := fs.Stat(fsys, path); err != nil {
				return err
			}
			return nil
		})
	}
	return eg.Wait()
}

// Exists will check if files exist at once, returning a map of the results.
func SomeExist(fsys fs.FS, paths ...string) (map[string]bool, error) {
	m := map[string]bool{}
	mu := sync.Mutex{}
	eg := new(errgroup.Group)
	for _, path := range paths {
		path := path
		eg.Go(func() error {
			if _, err := fs.Stat(fsys, path); err != nil {
				if !errors.Is(err, fs.ErrNotExist) {
					return err
				}
				return nil
			}
			mu.Lock()
			m[path] = true
			mu.Unlock()
			return nil
		})
	}
	return m, eg.Wait()
}
