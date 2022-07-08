package watcher

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

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

func getEvent(event <-chan []UpdateEvent) ([]UpdateEvent, error) {
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
	eventChan := make(chan []UpdateEvent)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return Watch(ctx, dir, func(events []UpdateEvent) error {
			select {
			case eventChan <- events:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err = os.WriteFile(filepath.Join(dir, "a.txt"), []byte("b"), 0644)
	is.NoErr(err)
	events, err := getEvent(eventChan)
	is.NoErr(err)
	is.Equal(len(events), 1)
	is.Equal(events[0].Path, filepath.Join(dir, "a.txt"))
	is.Equal(events[0].EventType, EditEventType)
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
	eventChan := make(chan []UpdateEvent)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return Watch(ctx, dir, func(events []UpdateEvent) error {
			select {
			case eventChan <- events:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err = os.RemoveAll(filepath.Join(dir, "a.txt"))
	is.NoErr(err)
	events, err := getEvent(eventChan)
	is.NoErr(err)
	is.Equal(len(events), 1)
	is.Equal(events[0].Path, filepath.Join(dir, "a.txt"))
	is.Equal(events[0].EventType, RemoveEventType)
	cancel()
	is.NoErr(eg.Wait())
}

func TestCreate(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	eventChan := make(chan []UpdateEvent)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return Watch(ctx, dir, func(events []UpdateEvent) error {
			select {
			case eventChan <- events:
			case <-ctx.Done():
			}
			return nil
		})
	})
	time.Sleep(waitForEvents)
	err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("b"), 0644)
	is.NoErr(err)
	events, err := getEvent(eventChan)
	is.NoErr(err)
	is.Equal(len(events), 1)
	is.Equal(events[0].Path, filepath.Join(dir, "a.txt"))
	cancel()
	is.NoErr(eg.Wait())
}

func TestCreateRecursive(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	eventChan := make(chan []UpdateEvent)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return Watch(ctx, dir, func(events []UpdateEvent) error {
			select {
			case eventChan <- events:
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
	events, err := getEvent(eventChan)
	is.NoErr(err)
	is.Equal(len(events), 2)
	expectedUpdates := []string{
		filepath.Join(dir, "b"),
		filepath.Join(dir, "b/a.txt"),
	}

	is.In(expectedUpdates, events[0].Path)
	is.In(expectedUpdates, events[1].Path)
	if events[0].Path == events[1].Path {
		is.Fail("expected paths to not be the same")
	}
	cancel()
	is.NoErr(eg.Wait())
}

func TestWithScaffold(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	eventChan := make(chan []UpdateEvent)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return Watch(ctx, dir, func(events []UpdateEvent) error {
			select {
			case eventChan <- events:
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
	events, err := getEvent(eventChan)
	is.NoErr(err)
	is.Equal(len(events), 5)
	expectedPaths := []string{
		filepath.Join(dir, "controller"),
		filepath.Join(dir, "controller/controller.go"),
		filepath.Join(dir, "view"),
		filepath.Join(dir, "view/index.svelte"),
		filepath.Join(dir, "view/show.svelte"),
	}

	for _, event := range events {
		is.In(expectedPaths, event.Path)
	}
	// Test that there's only been one event
	select {
	case <-eventChan:
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
	eventChan := make(chan []UpdateEvent)
	err := writeFiles(dir, map[string]string{
		"controller/controller.go": `package controller`,
		".envrc":                   `export FOO=bar`,
		".gitignore":               `.envrc`,
	})
	is.NoErr(err)
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return Watch(ctx, dir, func(events []UpdateEvent) error {
			select {
			case eventChan <- events:
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
	events, err := getEvent(eventChan)
	is.NoErr(err)
	is.Equal(len(events), 1)
	is.Equal(events[0].Path, filepath.Join(dir, "controller/controller.go"))
}
