package watcher

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/gitignore"

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
	// Files to ignore while walking the directory
	shouldIgnore := func(path string, de fs.DirEntry) error {
		if gitIgnore(path, de.IsDir()) || filepath.Base(path) == ".git" {
			return filepath.SkipDir
		}
		return nil
	}
	// Walk the files, adding files that aren't ignored
	walkDir := func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if err := shouldIgnore(path, de); err != nil {
			return err
		}
		watcher.Add(path)
		return nil
	}
	if err := filepath.WalkDir(dir, walkDir); err != nil {
		return err
	}
	// Trigger takes the walkDir above but will also trigger the fn abovre
	trigger := func(walk func(path string, de fs.DirEntry, err error) error) func(path string, de fs.DirEntry, err error) error {
		return func(path string, de fs.DirEntry, err error) error {
			if err := walk(path, de, err); err != nil {
				return err
			}
			return fn(path)
		}
	}
	// For some reason renames are often emitted instead of
	// Remove. Check it and correct.
	rename := func(path string) error {
		_, err := os.Stat(path)
		if nil == err {
			return nil
		}
		// If it's a different error, ignore
		if !errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		// Remove the path and emit an update
		watcher.Remove(path)
		// Trigger an update
		if err := fn(path); err != nil {
			return err
		}
		return nil
	}
	// Remove the file or directory from the watcher.
	// We intentionally ignore errors for this case.
	remove := func(path string) error {
		watcher.Remove(path)
		// Trigger an update
		if err := fn(path); err != nil {
			return err
		}
		return nil
	}
	// Watching a file or directory as long as it's not inside .gitignore.
	// Ignore most errors since missing a file isn't the end of the world.
	// If a new directory is created, add and trigger all the files within
	// that directory.
	create := func(path string) error {
		stat, err := os.Stat(path)
		if err != nil {
			return nil
		}
		if gitIgnore(path, stat.IsDir()) {
			return nil
		}
		err = watcher.Add(path)
		if err != nil {
			return err
		}
		// If it's a directory, walk the dir and trigger creates
		// because those create events won't happen on their own
		if stat.IsDir() {
			if err := filepath.WalkDir(path, trigger(walkDir)); err != nil {
				return err
			}
			return nil
		}
		// Otherwise, trigger the create
		if err := fn(path); err != nil {
			return err
		}
		return nil
	}
	// A file or directory has been updated. Notify our matchers.
	write := func(path string) error {
		fi, err := os.Stat(path)
		if err != nil {
			return nil
		}
		if gitIgnore(path, fi.IsDir()) {
			return nil
		}
		// Trigger an update
		if err := fn(path); err != nil {
			return err
		}
		return nil
	}

	// Watch for file events!
	// Note: The FAQ currently says it needs to be in a separate Go routine
	// https://github.com/fsnotify/fsnotify#faq, so we'll do that.
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case err := <-watcher.Errors:
				return err
			case evt := <-watcher.Events:
				switch op := evt.Op; {

				// Handle rename events
				case op&fsnotify.Rename != 0:
					if err := rename(evt.Name); err != nil {
						return err
					}

				// Handle remove events
				case op&fsnotify.Remove != 0:
					if err := remove(evt.Name); err != nil {
						return err
					}

				// Handle create events
				case op&fsnotify.Create != 0:
					if err := create(evt.Name); err != nil {
						return err
					}

				// Handle write events
				case op&fsnotify.Write != 0:
					if err := write(evt.Name); err != nil {
						return err
					}
				}
			}
		}
	})
	return eg.Wait()
}
