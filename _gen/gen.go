package gen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gitlab.com/mnm/bud/internal/pubsub"
	"golang.org/x/sync/errgroup"
)

// ErrSkipped allows you to skip generating files, without producing an error.
// TODO: consider moving to vfs.
// var ErrSkipped = errors.New("skipped")

type FS interface {
	Open(name string) (fs.File, error)
	Add(generators map[string]Generator)
	Subscribe(name string) (pubsub.Subscription, error)
}

type Generator interface {
	open(f F, key, relative, target string) (fs.File, error)
}

type F interface {
	fs.FS
	link(from, to string, event Event)
}

func New(dirfs fs.FS) *FileSystem {
	roots := map[string]bool{}
	dir := newDir(".")
	ps := pubsub.New()
	return &FileSystem{&innerFS{dir, dirfs, roots, ps, newGraph()}}
}

func root(path string) string {
	index := strings.Index(path, string(filepath.Separator))
	if index < 0 {
		return path
	}
	return path[0:index]
}

type FileSystem struct {
	ifs *innerFS
}

var _ FS = (*FileSystem)(nil)

func (d *FileSystem) Open(name string) (fs.File, error) {
	return d.ifs.Open(name)
}

// Add additional generators to GFS. This is not concurrency safe.
// TODO: merge generators if they exist already
func (d *FileSystem) Add(generators map[string]Generator) {
	for path, generator := range generators {
		d.ifs.roots[root(path)] = true
		d.ifs.dir.generators[path] = generator
	}
}

func (d *FileSystem) Subscribe(name string) (pubsub.Subscription, error) {
	if _, err := fs.Stat(d.ifs, name); err != nil {
		return nil, err
	}
	return d.ifs.ps.Subscribe(name), nil
}

func (d *FileSystem) Trigger(path string, event Event) {
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
	// Special case for the root. Synthesize the directory including real files
	// and generators
	if name == "." {
		return i.mergeEntries(name)
	}
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

// Entry set is used to ensure merged entries are unique by filename
func newEntrySet() *entrySet {
	return &entrySet{set: map[string]struct{}{}}
}

type entrySet struct {
	set     map[string]struct{}
	entries []fs.DirEntry
}

func (e *entrySet) Add(des ...fs.DirEntry) {
	for _, de := range des {
		name := de.Name()
		if _, ok := e.set[name]; !ok {
			e.entries = append(e.entries, de)
			e.set[name] = struct{}{}
		}
	}
}

func (e *entrySet) List() []fs.DirEntry {
	sort.Slice(e.entries, func(i, j int) bool {
		return e.entries[i].Name() < e.entries[j].Name()
	})
	return e.entries
}

// Merge the generator entries with the dirfs entries
// Currently only used for "."
func (i *innerFS) mergeEntries(name string) (fs.File, error) {
	file, err := i.dir.open(i, "", name, name)
	if err != nil {
		return nil, err
	}
	entries := newEntrySet()
	// Read all the entries from the generators
	if rdir, ok := file.(fs.ReadDirFile); ok {
		des, err := rdir.ReadDir(-1)
		if err != nil {
			return nil, err
		}
		entries.Add(des...)
	}
	// Read all the entries from dirfs
	if i.dirfs != nil {
		des, err := fs.ReadDir(i.dirfs, name)
		if err != nil {
			return nil, err
		}
		entries.Add(des...)
	}
	return &openDir{
		path:    name,
		entries: entries.List(),
	}, nil

}

func (i *innerFS) link(from, to string, event Event) {
	i.graph.Link(from, to, event)
}

// SkipUnless will return ErrSkipped unless all the paths exists
func SkipUnless(f fs.FS, paths ...string) error {
	eg := new(errgroup.Group)
	for _, path := range paths {
		path := path
		eg.Go(func() error {
			if _, err := fs.Stat(f, path); err != nil {
				return fmt.Errorf("%w %q", fs.ErrNotExist, path)
			}
			return nil
		})
	}
	return eg.Wait()
}
