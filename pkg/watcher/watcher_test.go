package watcher_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/livebud/bud/pkg/watcher"
	"github.com/matryer/is"
	"golang.org/x/sync/errgroup"
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

func getEvent(event <-chan *watcher.Events) (*watcher.Events, error) {
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
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	ctx := context.Background()
	eventCh := make(chan *watcher.Events)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(events *watcher.Events) error {
			select {
			case eventCh <- events:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("b"), 0644))
	events, err := getEvent(eventCh)
	is.NoErr(err)
	is.Equal(len(events.Updated), 1)
	is.Equal(events.Updated[0], "a.txt")
	cancel()
	is.NoErr(eg.Wait())
}

func TestDelete(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	ctx := context.Background()
	eventCh := make(chan *watcher.Events)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(events *watcher.Events) error {
			select {
			case eventCh <- events:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	is.NoErr(os.RemoveAll(filepath.Join(dir, "a.txt")))
	events, err := getEvent(eventCh)
	is.NoErr(err)
	is.Equal(len(events.Deleted), 1)
	is.Equal(events.Deleted[0], "a.txt")
	cancel()
	is.NoErr(eg.Wait())
}

func TestCreates(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	eventCh := make(chan *watcher.Events)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(events *watcher.Events) error {
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
	is.Equal(len(events.Created), 1)
	is.Equal(events.Created[0], "a.txt")
	cancel()
	is.NoErr(eg.Wait())
}

func TestCreatesRecursive(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	eventCh := make(chan *watcher.Events)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(events *watcher.Events) error {
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
	is.Equal(len(events.Created), 2)
	is.Equal(events.Created[0], "b")
	is.Equal(events.Created[1], "b/a.txt")
	cancel()
	is.NoErr(eg.Wait())
}

func TestWithScaffold(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	eventCh := make(chan *watcher.Events)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(events *watcher.Events) error {
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
	is.Equal(len(events.Created), 5)
	is.Equal(events.Created[0], "controller")
	is.Equal(events.Created[1], "controller/controller.go")
	is.Equal(events.Created[2], "view")
	is.Equal(events.Created[3], "view/index.svelte")
	is.Equal(events.Created[4], "view/show.svelte")
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
	eventCh := make(chan *watcher.Events)
	err := writeFiles(dir, map[string]string{
		"controller/controller.go": `package controller`,
		".envrc":                   `export FOO=bar`,
		".gitignore":               `.envrc`,
	})
	is.NoErr(err)
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(events *watcher.Events) error {
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
	is.Equal(len(events.Updated), 1)
	is.Equal(events.Updated[0], "controller/controller.go")
}

func TestRename(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	ctx := context.Background()
	eventCh := make(chan *watcher.Events)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return watcher.Watch(ctx, dir, func(events *watcher.Events) error {
			select {
			case eventCh <- events:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	is.NoErr(os.Rename(filepath.Join(dir, "a.txt"), filepath.Join(dir, "b.txt")))
	events, err := getEvent(eventCh)
	is.NoErr(err)
	is.Equal(len(events.Deleted), 1)
	is.Equal(events.Deleted[0], "a.txt")
	is.Equal(len(events.Created), 1)
	is.Equal(events.Created[0], "b.txt")
	cancel()
	is.NoErr(eg.Wait())
}
