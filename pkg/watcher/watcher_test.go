package watcher_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gitlab.com/mnm/bud/pkg/watcher"
	"golang.org/x/sync/errgroup"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/pkg/vfs"
)

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
			event <- path
			return nil
		})
	})
	time.Sleep(100 * time.Millisecond)
	err = os.WriteFile(filepath.Join(dir, "a.txt"), []byte("b"), 0644)
	is.NoErr(err)
	select {
	case path := <-event:
		is.Equal(path, filepath.Join(dir, "a.txt"))
	case <-time.After(time.Second):
		t.Fatal("timed out while waiting for watcher")
	}
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
			event <- path
			return nil
		})
	})
	time.Sleep(100 * time.Millisecond)
	err = os.RemoveAll(filepath.Join(dir, "a.txt"))
	is.NoErr(err)
	select {
	case path := <-event:
		is.Equal(path, filepath.Join(dir, "a.txt"))
	case <-time.After(time.Second):
		t.Fatal("timed out while waiting for watcher")
	}
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
			event <- path
			return nil
		})
	})
	time.Sleep(100 * time.Millisecond)
	err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("b"), 0644)
	is.NoErr(err)
	select {
	case path := <-event:
		is.Equal(path, filepath.Join(dir, "a.txt"))
	case <-time.After(time.Second):
		t.Fatal("timed out while waiting for watcher")
	}
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
			event <- path
			return nil
		})
	})
	time.Sleep(100 * time.Millisecond)
	err := os.MkdirAll(filepath.Join(dir, "b"), 0755)
	is.NoErr(err)
	select {
	case path := <-event:
		is.Equal(path, filepath.Join(dir, "b"))
	case <-time.After(time.Second):
		t.Fatal("timed out while waiting for watcher")
	}
	err = os.WriteFile(filepath.Join(dir, "b", "a.txt"), []byte("b"), 0644)
	is.NoErr(err)
	select {
	case path := <-event:
		is.Equal(path, filepath.Join(dir, "b", "a.txt"))
	case <-time.After(time.Second):
		t.Fatal("timed out while waiting for watcher")
	}
	cancel()
	is.NoErr(eg.Wait())
}
