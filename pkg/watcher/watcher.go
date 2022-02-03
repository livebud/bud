package watcher

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"gitlab.com/mnm/bud/internal/gitignore"

	"github.com/fsnotify/fsnotify"
)

// Watch function
func Watch(ctx context.Context, dir string, fn func(path string) error) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	gitIgnore := gitignore.From(dir)
	walker := func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		isDir := de.IsDir()
		if gitIgnore(path, isDir) || filepath.Base(path) == ".git" {
			if isDir {
				return filepath.SkipDir
			}
			return nil
		}
		watcher.Add(path)
		return nil
	}
	if err := filepath.WalkDir(dir, walker); err != nil {
		return err
	}
	// watch for file events!
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-watcher.Errors:
			return err
		case evt := <-watcher.Events:
			switch evt.Op {
			// Ignore "CHMOD" events.
			case fsnotify.Chmod:
				continue
				// For some reason it appears that renames are often emitted instead of
				// Remove. Check it and correct
			case fsnotify.Rename:
				if _, err := os.Stat(evt.Name); err != nil {
					// If it's a different error, ignore
					if !errors.Is(err, fs.ErrNotExist) {
						continue
					}
					// Remove the path and emit an update
					watcher.Remove(evt.Name)
					// Trigger an update
					if err := fn(evt.Name); err != nil {
						return err
					}
					continue
				}

			// Remove the file or directory from the watcher.
			// We intentionally ignore errors for this case.
			case fsnotify.Remove:
				watcher.Remove(evt.Name)
				// Trigger an update
				if err := fn(evt.Name); err != nil {
					return err
				}
				continue
			// Try watching a the file as long as it's not inside .gitignore.
			// Ignore most errors since missing a file isn't the end of the world.
			case fsnotify.Create:
				fi, err := os.Stat(evt.Name)
				if err != nil {
					continue
				}
				if gitIgnore(evt.Name, fi.IsDir()) {
					continue
				}
				err = watcher.Add(evt.Name)
				if err != nil {
					return err
				}
				// Trigger an update
				if err := fn(evt.Name); err != nil {
					return err
				}
				continue
			// A file has been updated. Notify our matchers.
			case fsnotify.Write:
				fi, err := os.Stat(evt.Name)
				if err != nil {
					continue
				}
				if gitIgnore(evt.Name, fi.IsDir()) {
					continue
				}
				// Trigger an update
				if err := fn(evt.Name); err != nil {
					return err
				}
			}
		}
	}
}
