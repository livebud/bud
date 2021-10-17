package bfs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/mnm/bud/internal/pubsub"
	"golang.org/x/sync/errgroup"
)

func New(dirfs fs.FS) *GFS {
	roots := map[string]bool{}
	dir := newDir(".")
	ps := pubsub.New()
	return &GFS{&innerFS{dir, dirfs, roots, ps, newGraph()}}
}

func root(path string) string {
	index := strings.Index(path, string(filepath.Separator))
	if index < 0 {
		return path
	}
	return path[0:index]
}

type GFS struct {
	ifs *innerFS
}

var _ BFS = (*GFS)(nil)

func (d *GFS) Open(name string) (fs.File, error) {
	return d.ifs.Open(name)
}

// Add additional generators to GFS. This is not concurrency safe.
// TODO: merge generators if they exist already
func (d *GFS) Add(generators map[string]Generator) {
	for path, generator := range generators {
		d.ifs.roots[root(path)] = true
		d.ifs.dir.generators[path] = generator
	}
}

func (d *GFS) Subscribe(name string) (pubsub.Subscription, error) {
	if _, err := fs.Stat(d.ifs, name); err != nil {
		return nil, err
	}
	return d.ifs.ps.Subscribe(name), nil
}

func (d *GFS) Trigger(path string, event Event) {
	nodes := d.ifs.graph.Ins(path, event)
	for _, node := range nodes {
		d.ifs.ps.Publish(node, []byte(event.String()))
	}
}

type innerFS struct {
	dir   *Dir
	dirfs fs.FS
	roots map[string]bool
	ps    pubsub.Client
	graph *graph
}

// Open the file
func (i *innerFS) Open(name string) (fs.File, error) {
	file, err := i.open(name)
	if err != nil {
		return nil, fmt.Errorf("open %s > %w", name, err)
	}
	return file, nil
}

func (i *innerFS) open(name string) (fs.File, error) {
	key := root(name)
	// Test if we should look within the generator filesystem or the real
	// filesystem.
	if _, ok := i.roots[key]; ok {
		return i.dir.open(i, "", name, name)
	}
	if i.dirfs == nil {
		return nil, fs.ErrNotExist
	}
	file, err := i.dirfs.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fs.ErrNotExist
		}
		return nil, err
	}
	return file, nil
}

func (i *innerFS) link(from, to string, event Event) {
	i.graph.Link(from, to, event)
}

func Exists(f fs.FS, paths ...string) error {
	eg := new(errgroup.Group)
	for _, path := range paths {
		path := path
		eg.Go(func() error {
			if _, err := fs.Stat(f, path); err != nil {
				return fmt.Errorf("exists %s > %w", path, err)
			}
			return nil
		})
	}
	return eg.Wait()
}
