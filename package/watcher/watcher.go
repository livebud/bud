package watcher

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/gitignore"

	"github.com/fsnotify/fsnotify"
)

var Stop = errors.New("stop watching")

// Watch function
func Watch(ctx context.Context, dir string, fn func(path string) error) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	// Avoid duplicate events by checking the stamp of the file.
	// TODO: bound this map
	duplicates := map[string]struct{}{}
	isDuplicate := func(path string, stat fs.FileInfo) bool {
		stamp, err := computeStamp(path, stat)
		if err != nil {
			return false
		}
		// Duplicate check
		if _, ok := duplicates[stamp]; ok {
			return true
		}
		duplicates[stamp] = struct{}{}
		return false
	}
	// Note which paths we've seen already. This is similar to deduping, but for
	// rename and remove events where the files don't exist anymore.
	// TODO: bound this map
	seen := map[string]struct{}{}
	hasSeen := func(path string) bool {
		if _, ok := seen[path]; ok {
			return true
		}
		seen[path] = struct{}{}
		return false
	}
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
		// Don't do anything if we've already seen this path
		if hasSeen(path) {
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
		// Don't do anything if we've already seen this path
		if hasSeen(path) {
			return nil
		}
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
		// Reset seen path
		delete(seen, path)
		// Stat the file
		stat, err := os.Stat(path)
		if err != nil {
			return nil
		}
		if gitIgnore(path, stat.IsDir()) {
			return nil
		}
		if isDuplicate(path, stat) {
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
		// Reset seen path
		delete(seen, path)
		// Stat the file
		stat, err := os.Stat(path)
		if err != nil {
			return nil
		}
		if gitIgnore(path, stat.IsDir()) {
			return nil
		}
		if isDuplicate(path, stat) {
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
	// Wait for the watcher to complete
	if err := eg.Wait(); err != nil {
		if !errors.Is(err, Stop) {
			return err
		}
		return nil
	}
	return nil
}

// computeStamp uses path, size, mode and modtime to try and ensure this is a
// unique event.
func computeStamp(path string, stat fs.FileInfo) (stamp string, err error) {
	mtime := stat.ModTime().UnixNano()
	mode := stat.Mode()
	size := stat.Size()
	stamp = path + ":" + strconv.Itoa(int(size)) + ":" + mode.String() + ":" + strconv.Itoa(int(mtime))
	return stamp, nil
}
