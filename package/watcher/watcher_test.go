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

func getEvent(event <-chan []watcher.Event) ([]watcher.Event, error) {
	select {
	case events := <-event:
		return events, nil
	case <-time.After(1 * time.Second):
		return nil, errors.New("timed out while waiting for watcher events")
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
	eventCh := make(chan []watcher.Event)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(events []watcher.Event) error {
			select {
			case eventCh <- events:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err = os.WriteFile(filepath.Join(dir, "a.txt"), []byte("b"), 0644)
	is.NoErr(err)
	events, err := getEvent(eventCh)
	is.NoErr(err)
	is.Equal(len(events), 1)
	is.Equal(events[0].Path, filepath.Join(dir, "a.txt"))
	is.Equal(events[0].Op, watcher.OpUpdate)
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
	eventCh := make(chan []watcher.Event)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(events []watcher.Event) error {
			select {
			case eventCh <- events:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err = os.RemoveAll(filepath.Join(dir, "a.txt"))
	is.NoErr(err)
	events, err := getEvent(eventCh)
	is.NoErr(err)
	is.Equal(len(events), 1)
	is.Equal(events[0].Path, filepath.Join(dir, "a.txt"))
	is.Equal(events[0].Op, watcher.OpDelete)
	cancel()
	is.NoErr(eg.Wait())
}

func TestCreate(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	eventCh := make(chan []watcher.Event)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(events []watcher.Event) error {
			select {
			case eventCh <- events:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("b"), 0644)
	is.NoErr(err)
	events, err := getEvent(eventCh)
	is.NoErr(err)
	is.Equal(len(events), 1)
	is.Equal(events[0].Path, filepath.Join(dir, "a.txt"))
	is.Equal(events[0].Op, watcher.OpCreate)
	cancel()
	is.NoErr(eg.Wait())
}

func TestCreateRecursive(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	eventCh := make(chan []watcher.Event)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(events []watcher.Event) error {
			select {
			case eventCh <- events:
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
	events, err := getEvent(eventCh)
	is.NoErr(err)
	is.Equal(len(events), 2)
	is.Equal(events[0].Path, filepath.Join(dir, "b"))
	is.Equal(events[0].Op, watcher.OpCreate)
	is.Equal(events[1].Path, filepath.Join(dir, "b/a.txt"))
	is.Equal(events[0].Op, watcher.OpCreate)
	cancel()
	is.NoErr(eg.Wait())
}

func TestWithScaffold(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	eventCh := make(chan []watcher.Event)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(events []watcher.Event) error {
			select {
			case eventCh <- events:
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
	events, err := getEvent(eventCh)
	is.NoErr(err)
	is.Equal(len(events), 5)
	is.Equal(events[0].Path, filepath.Join(dir, "controller"))
	is.Equal(events[0].Op, watcher.OpCreate)
	is.Equal(events[1].Path, filepath.Join(dir, "controller/controller.go"))
	is.Equal(events[1].Op, watcher.OpCreate)
	is.Equal(events[2].Path, filepath.Join(dir, "view"))
	is.Equal(events[2].Op, watcher.OpCreate)
	is.Equal(events[3].Path, filepath.Join(dir, "view/index.svelte"))
	is.Equal(events[3].Op, watcher.OpCreate)
	is.Equal(events[4].Path, filepath.Join(dir, "view/show.svelte"))
	is.Equal(events[4].Op, watcher.OpCreate)
	// Test that there's only been one event
	select {
	case <-eventCh:
		t.Fatalf("unexpected extra event")
	case <-time.Tick(waitForEvents):
	}
	cancel()
	is.NoErr(eg.Wait())
}

func TestWithRootDotFile(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	eventCh := make(chan []watcher.Event)
	err := writeFiles(dir, map[string]string{
		"controller/controller.go": `package controller`,
		".envrc":                   `export FOO=bar`,
		".gitignore":               `.envrc`,
	})
	is.NoErr(err)
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(events []watcher.Event) error {
			select {
			case eventCh <- events:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	// Update the event
	err = writeFiles(dir, map[string]string{
		"controller/controller.go": `package controller2`,
	})
	is.NoErr(err)
	// Get event
	events, err := getEvent(eventCh)
	is.NoErr(err)
	is.Equal(len(events), 1)
	is.Equal(events[0].Path, filepath.Join(dir, "controller/controller.go"))
	is.Equal(events[0].Op, watcher.OpUpdate)
}

func TestRename(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	err := vfs.Write(dir, vfs.Map{
		"a.txt": []byte(`a`),
	})
	is.NoErr(err)
	ctx := context.Background()
	eventCh := make(chan []watcher.Event)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(events []watcher.Event) error {
			select {
			case eventCh <- events:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err = os.Rename(filepath.Join(dir, "a.txt"), filepath.Join(dir, "b.txt"))
	is.NoErr(err)
	events, err := getEvent(eventCh)
	is.NoErr(err)
	is.Equal(len(events), 2)
	is.Equal(events[0].Path, filepath.Join(dir, "b.txt"))
	is.Equal(events[0].Op, watcher.OpCreate)
	is.Equal(events[1].Path, filepath.Join(dir, "a.txt"))
	is.Equal(events[1].Op, watcher.OpDelete)
	cancel()
	is.NoErr(eg.Wait())
}
