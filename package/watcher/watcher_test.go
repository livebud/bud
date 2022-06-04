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

func getEvent(event <-chan []string) ([]string, error) {
	select {
	case paths := <-event:
		return paths, nil
	case <-time.After(1 * time.Second):
		return nil, errors.New("timed out while waiting for watcher")
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
	event := make(chan []string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(paths []string) error {
			select {
			case event <- paths:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err = os.WriteFile(filepath.Join(dir, "a.txt"), []byte("b"), 0644)
	is.NoErr(err)
	paths, err := getEvent(event)
	is.NoErr(err)
	is.Equal(len(paths), 1)
	is.Equal(paths[0], filepath.Join(dir, "a.txt"))
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
	event := make(chan []string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(paths []string) error {
			select {
			case event <- paths:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err = os.RemoveAll(filepath.Join(dir, "a.txt"))
	is.NoErr(err)
	paths, err := getEvent(event)
	is.NoErr(err)
	is.Equal(len(paths), 1)
	is.Equal(paths[0], filepath.Join(dir, "a.txt"))
	cancel()
	is.NoErr(eg.Wait())
}

func TestCreate(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	event := make(chan []string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(paths []string) error {
			select {
			case event <- paths:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("b"), 0644)
	is.NoErr(err)
	paths, err := getEvent(event)
	is.NoErr(err)
	is.Equal(len(paths), 1)
	is.Equal(paths[0], filepath.Join(dir, "a.txt"))
	cancel()
	is.NoErr(eg.Wait())
}

func TestCreateRecursive(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	event := make(chan []string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(paths []string) error {
			select {
			case event <- paths:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err := os.MkdirAll(filepath.Join(dir, "b"), 0755)
	is.NoErr(err)
	err = os.WriteFile(filepath.Join(dir, "b", "a.txt"), []byte("b"), 0644)
	is.NoErr(err)
	paths, err := getEvent(event)
	is.NoErr(err)
	is.Equal(len(paths), 2)
	is.Equal(paths[0], filepath.Join(dir, "b"))
	is.Equal(paths[1], filepath.Join(dir, "b/a.txt"))
	cancel()
	is.NoErr(eg.Wait())
}

func TestWithScaffold(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	event := make(chan []string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(paths []string) error {
			select {
			case event <- paths:
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
	paths, err := getEvent(event)
	is.NoErr(err)
	is.Equal(len(paths), 5)
	is.Equal(paths[0], filepath.Join(dir, "controller"))
	is.Equal(paths[1], filepath.Join(dir, "controller/controller.go"))
	is.Equal(paths[2], filepath.Join(dir, "view"))
	is.Equal(paths[3], filepath.Join(dir, "view/index.svelte"))
	is.Equal(paths[4], filepath.Join(dir, "view/show.svelte"))
	// Test that there's only been one event
	select {
	case <-event:
		t.Fatalf("unexpected extra event")
	case <-time.Tick(waitForEvents):
	}
	cancel()
	is.NoErr(eg.Wait())
}
