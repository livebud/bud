package watcher_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/livebud/bud/package/watcher"
	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/vfs"
)

var waitForEvents = 500 * time.Millisecond

func writeFiles(dir string, files map[string]string) error {
	eg := new(errgroup.Group)
	for path, data := range files {
		path := filepath.Join(dir, path)
		data := data
		eg.Go(func() error {
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			return os.WriteFile(path, []byte(data), 0644)
		})
	}
	return eg.Wait()
}

func getEvent(event <-chan string) (string, error) {
	select {
	case path := <-event:
		return path, nil
	case <-time.After(1 * time.Second):
		return "", errors.New("timed out while waiting for watcher")
	}
}

func TestChange(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	err := vfs.Write(dir, vfs.Map{
		"a.txt": []byte(`a`),
	})
	is.NoErr(err)
	ctx := context.Background()
	event := make(chan string, 1)
	ctx, cancel := context.WithCancel(context.Background())
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(path string) error {
			select {
			case event <- path:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err = os.WriteFile(filepath.Join(dir, "a.txt"), []byte("b"), 0644)
	is.NoErr(err)
	path, err := getEvent(event)
	is.NoErr(err)
	is.Equal(path, filepath.Join(dir, "a.txt"))
	cancel()
	is.NoErr(eg.Wait())
}

func TestDelete(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	err := vfs.Write(dir, vfs.Map{
		"a.txt": []byte(`a`),
	})
	is.NoErr(err)
	ctx := context.Background()
	event := make(chan string, 1)
	ctx, cancel := context.WithCancel(context.Background())
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(path string) error {
			select {
			case event <- path:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err = os.RemoveAll(filepath.Join(dir, "a.txt"))
	is.NoErr(err)
	path, err := getEvent(event)
	is.NoErr(err)
	is.Equal(path, filepath.Join(dir, "a.txt"))
	cancel()
	is.NoErr(eg.Wait())
}

func TestCreate(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	event := make(chan string, 1)
	ctx, cancel := context.WithCancel(context.Background())
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(path string) error {
			select {
			case event <- path:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("b"), 0644)
	is.NoErr(err)
	path, err := getEvent(event)
	is.NoErr(err)
	is.Equal(path, filepath.Join(dir, "a.txt"))
	cancel()
	is.NoErr(eg.Wait())
}

func TestCreateRecursive(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	event := make(chan string, 1)
	ctx, cancel := context.WithCancel(context.Background())
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(path string) error {
			select {
			case event <- path:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err := os.MkdirAll(filepath.Join(dir, "b"), 0755)
	is.NoErr(err)
	path, err := getEvent(event)
	is.NoErr(err)
	is.Equal(path, filepath.Join(dir, "b"))
	err = os.WriteFile(filepath.Join(dir, "b", "a.txt"), []byte("b"), 0644)
	is.NoErr(err)
	path, err = getEvent(event)
	is.NoErr(err)
	is.Equal(path, filepath.Join(dir, "b", "a.txt"))
	cancel()
	is.NoErr(eg.Wait())
}

func TestWithScaffold(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	event := make(chan string, 4)
	ctx, cancel := context.WithCancel(context.Background())
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(path string) error {
			select {
			case event <- path:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err := writeFiles(dir, map[string]string{
		"controller/controller.go": `package controller`,
		"view/index.svelte":        `<h1>index</h1>`,
		"view/show.svelte":         `<h1>show</h1>`,
	})
	is.NoErr(err)
	paths := map[string]bool{}
	for i := 1; i <= 5; i++ {
		path, err := getEvent(event)
		is.NoErr(err)
		rel, err := filepath.Rel(dir, path)
		is.NoErr(err)
		paths[rel] = true
	}
	is.True(paths["controller"])
	is.True(paths["controller/controller.go"])
	is.True(paths["view"])
	is.True(paths["view/index.svelte"])
	is.True(paths["view/show.svelte"])
	// While there should be no more events, testing this can be flaky in CI.
	// Instead test that we don't have any events with unexpected paths.
	// An extra event isn't the end of the world, it'll just reload one more time.
	select {
	case path := <-event:
		rel, err := filepath.Rel(dir, path)
		is.NoErr(err)
		if _, ok := paths[rel]; !ok {
			t.Fatalf("unexpected event: %q", path)
		}
	case <-time.Tick(waitForEvents):
	}
	cancel()
	is.NoErr(eg.Wait())
}
