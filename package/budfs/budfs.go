package budfs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/livebud/bud/package/budfs/linkmap"

	"github.com/livebud/bud/package/virtual/vcache"

	"github.com/livebud/bud/package/virtual"

	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/glob"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/internal/orderedset"
	"github.com/livebud/bud/internal/valid"
	"github.com/livebud/bud/package/budfs/mergefs"
	"github.com/livebud/bud/package/budfs/treefs"
	"github.com/livebud/bud/package/log"
)

func New(fsys fs.FS, log log.Log) *FileSystem {
	// Exclude the underlying filesystem (often os) from contributing bud/* files.
	// The bud/* directory is owned by the generator filesytem.
	fsys = virtual.Exclude(fsys, func(path string) bool {
		return path == "bud" || strings.HasPrefix(path, "bud/")
	})
	cache := vcache.New()
	node := treefs.New(".")
	mountfs := &mountFS{}
	merged := mergefs.Merge(node, mountfs, fsys)
	return &FileSystem{
		cache,
		new(once.Closer),
		mountfs,
		merged,
		node,
		linkmap.New(log),
		log,
	}
}

type FileSystem struct {
	cache   vcache.Cache
	closer  *once.Closer
	mountfs *mountFS
	fsys    fs.FS
	node    *treefs.Node
	lmap    *linkmap.Map
	log     log.Log
}

type File struct {
	Data   []byte
	path   string
	mode   fs.FileMode
	target string
}

func (f *File) Target() string {
	return f.target
}

func (f *File) Relative() string {
	return relativePath(f.path, f.target)
}

func (f *File) Path() string {
	return f.path
}

func (f *File) Mode() fs.FileMode {
	return f.mode
}

type FS interface {
	fs.FS
	fs.ReadDirFS
	fs.GlobFS
	Watch(paths ...string) error
	Context() context.Context
	Defer(func() error)
}

type Dir struct {
	fsys   *FileSystem
	node   *treefs.Node
	target string
}

func (d *Dir) Target() string {
	return d.target
}

func (d *Dir) Relative() string {
	return relativePath(d.node.Path(), d.target)
}

func (d *Dir) Path() string {
	return d.node.Path()
}

func (d *Dir) Mode() fs.FileMode {
	return d.node.Mode()
}

func (d *Dir) GenerateFile(path string, fn func(fsys FS, file *File) error) {
	fileg := &fileGenerator{d.fsys, fn, nil, path}
	fileg.node = d.node.FileGenerator(path, fileg)
}

func (d *Dir) FileGenerator(path string, generator FileGenerator) {
	d.GenerateFile(path, generator.GenerateFile)
}

func (d *Dir) GenerateDir(dir string, fn func(fsys FS, dir *Dir) error) {
	dirg := &dirGenerator{d.fsys, fn, nil}
	dirg.node = d.node.DirGenerator(dir, dirg)
}

func (d *Dir) DirGenerator(dir string, generator DirGenerator) {
	d.GenerateDir(dir, generator.GenerateDir)
}

func (d *Dir) ServeFile(dir string, fn func(fsys FS, file *File) error) {
	fileg := &fileServer{d.fsys, fn, nil, dir}
	fileg.node = d.node.DirGenerator(dir, fileg)
}

func (d *Dir) FileServer(dir string, generator FileGenerator) {
	d.ServeFile(dir, generator.GenerateFile)
}

type mountGenerator struct {
	dir  string
	fsys fs.FS
}

func (g *mountGenerator) Generate(target string) (fs.File, error) {
	return g.fsys.Open(relativePath(g.dir, target))
}

func (d *Dir) Mount(mount fs.FS) error {
	des, err := fs.ReadDir(mount, ".")
	if err != nil {
		return fmt.Errorf("budfs: mount error. %w", err)
	}
	// Wrap mount in the existing generator cache
	mountg := &mountGenerator{d.node.Path(), mount}
	// Loop over the first level and add the mount, allowing us to mount "."
	// on an existing directory
	for _, de := range des {
		if de.IsDir() {
			d.node.DirGenerator(de.Name(), mountg)
			continue
		}
		d.node.FileGenerator(de.Name(), mountg)
	}
	return nil
}

type FileGenerator interface {
	GenerateFile(fsys FS, file *File) error
}

type GenerateFile func(fsys FS, file *File) error

func (fn GenerateFile) GenerateFile(fsys FS, file *File) error {
	return fn(fsys, file)
}

type DirGenerator interface {
	GenerateDir(fsys FS, dir *Dir) error
}

type GenerateDir func(fsys FS, dir *Dir) error

func (fn GenerateDir) GenerateDir(fsys FS, dir *Dir) error {
	return fn(fsys, dir)
}

type EmbedFile struct {
	Data []byte
}

var _ FileGenerator = (*EmbedFile)(nil)

func (e *EmbedFile) GenerateFile(fsys FS, file *File) error {
	file.Data = e.Data
	return nil
}

func (f *FileSystem) Open(name string) (fs.File, error) {
	file, err := f.fsys.Open(name)
	if err != nil {
		return nil, fmt.Errorf("budfs: open %q. %w", name, err)
	}
	return file, nil
}

func (f *FileSystem) Close() error {
	return f.closer.Close()
}

func (f *FileSystem) Dir() *Dir {
	return &Dir{fsys: f, node: f.node, target: "."}
}

type fileGenerator struct {
	fsys *FileSystem
	fn   func(fsys FS, file *File) error
	node *treefs.Node
	path string
}

func (g *fileGenerator) Generate(target string) (fs.File, error) {
	if entry, ok := g.fsys.cache.Get(target); ok {
		return virtual.New(entry), nil
	}
	fctx := &fileSystem{context.TODO(), g.fsys, g.fsys.lmap.Scope(target)}
	file := &File{nil, g.path, g.node.Mode(), target}
	g.fsys.log.Fields(log.Fields{
		"target": target,
		"path":   g.path,
	}).Debug("budfs: running file generator function")
	if err := g.fn(fctx, file); err != nil {
		return nil, err
	}
	vfile := &virtual.File{
		Path: g.node.Path(),
		Mode: g.node.Mode(),
		Data: file.Data,
	}
	g.fsys.cache.Set(target, vfile)
	return virtual.New(vfile), nil
}

func (f *FileSystem) GenerateFile(path string, fn func(fsys FS, file *File) error) {
	fileg := &fileGenerator{f, fn, nil, path}
	fileg.node = f.node.FileGenerator(path, fileg)
}

func (f *FileSystem) FileGenerator(path string, generator FileGenerator) {
	f.GenerateFile(path, generator.GenerateFile)
}

type dirGenerator struct {
	fsys *FileSystem
	fn   func(fsys FS, dir *Dir) error
	node *treefs.Node
}

func (g *dirGenerator) Generate(target string) (fs.File, error) {
	if _, ok := g.fsys.cache.Get(g.node.Path()); ok {
		return g.node.Open(target)
	}
	// Clear the subdirectories
	g.node.Clear()
	fctx := &fileSystem{context.TODO(), g.fsys, g.fsys.lmap.Scope(target)}
	dir := &Dir{g.fsys, g.node, target}
	g.fsys.log.Fields(log.Fields{
		"target": target,
		"path":   g.node.Path(),
	}).Debug("budfs: running dir generator function")
	if err := g.fn(fctx, dir); err != nil {
		return nil, err
	}
	g.fsys.cache.Set(g.node.Path(), &virtual.Dir{
		Path:    g.node.Path(),
		Mode:    g.node.Mode(),
		Entries: g.node.Entries(),
	})
	return g.node.Open(target)
}

func (f *FileSystem) GenerateDir(path string, fn func(fsys FS, dir *Dir) error) {
	dirg := &dirGenerator{f, fn, nil}
	dirg.node = f.node.DirGenerator(path, dirg)
}

func (f *FileSystem) DirGenerator(path string, generator DirGenerator) {
	f.GenerateDir(path, generator.GenerateDir)
}

type fileServer struct {
	fsys *FileSystem
	fn   func(fsys FS, file *File) error
	node *treefs.Node
	path string
}

func (g *fileServer) Generate(target string) (fs.File, error) {
	if entry, ok := g.fsys.cache.Get(target); ok {
		return virtual.New(entry), nil
	}
	// Always return an empty directory if we request the root
	rel := relativePath(g.node.Path(), target)
	if rel == "." {
		return virtual.New(&virtual.Dir{
			Path: g.path,
			Mode: fs.ModeDir,
		}), nil
	}
	fctx := &fileSystem{context.TODO(), g.fsys, g.fsys.lmap.Scope(target)}
	// File differs slightly than others because g.node.Path() is the directory
	// path, but we want the target path for serving files.
	file := &File{nil, g.path, g.node.Mode(), target}
	g.fsys.log.Fields(log.Fields{
		"target": target,
		"path":   g.node.Path(),
	}).Debug("budfs: running file server function")
	if err := g.fn(fctx, file); err != nil {
		return nil, err
	}
	vfile := &virtual.File{
		Path: target,
		Mode: fs.FileMode(0),
		Data: file.Data,
	}
	g.fsys.cache.Set(target, vfile)
	return virtual.New(vfile), nil
}

func (f *FileSystem) ServeFile(dir string, fn func(fsys FS, file *File) error) {
	fileg := &fileServer{f, fn, nil, dir}
	fileg.node = f.node.DirGenerator(dir, fileg)
}

func (f *FileSystem) FileServer(dir string, generator FileGenerator) {
	f.ServeFile(dir, generator.GenerateFile)
}

type mountFS struct {
	mu   sync.RWMutex
	fsys fs.FS
}

func (m *mountFS) Set(fsys fs.FS) {
	m.mu.Lock()
	m.fsys = fsys
	m.mu.Unlock()
}

func (m *mountFS) Open(name string) (fs.File, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.fsys == nil {
		return nil, fmt.Errorf("budfs: open from mount %q. no filesystem mounted. %w", name, fs.ErrNotExist)
	}
	file, err := m.fsys.Open(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *FileSystem) Mount(fsys fs.FS) {
	f.mountfs.Set(fsys)
}

// Sync the overlay to the filesystem
func (f *FileSystem) Sync(writable virtual.FS, to string, options ...dsync.Option) error {
	// Temporarily replace the underlying fs.FS with a cached fs.FS
	cache := vcache.New()
	fsys := f.fsys
	f.fsys = vcache.Wrap(cache, fsys, f.log)
	err := dsync.To(f.fsys, writable, to, options...)
	f.fsys = fsys
	return err
}

// Change updates the cache
func (f *FileSystem) Change(paths ...string) {
	for i := 0; i < len(paths); i++ {
		path := paths[i]
		if f.cache.Has(path) {
			f.log.Debug("budfs: delete cache path %q", path)
			f.cache.Delete(path)
		}
		f.lmap.Range(func(genPath string, fns *linkmap.List) bool {
			if f.cache.Has(genPath) && fns.Check(path) {
				paths = append(paths, genPath)
			}
			return true
		})
	}
}

type fileSystem struct {
	ctx  context.Context
	fsys *FileSystem
	link *linkmap.List
}

var _ FS = (*fileSystem)(nil)

// Open implements fs.FS
func (f *fileSystem) Open(name string) (fs.File, error) {
	f.link.Link("open", name)
	file, err := f.fsys.Open(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Watch the paths for changes
func (f *fileSystem) Watch(paths ...string) error {
	for _, path := range paths {
		// Not a glob
		if glob.Base(path) == path {
			f.link.Link("watch", path)
			continue
		}
		// Compile the pattern into a glob matcher
		matcher, err := glob.Compile(path)
		if err != nil {
			return err
		}
		// Watch for changes to the pattern
		f.link.Select("watch", func(path string) bool {
			return matcher.Match(path)
		})
	}
	return nil
}

func (f *fileSystem) Context() context.Context {
	return f.ctx
}

// Defer a function until close is called. May be called multiple times if
// generators are triggered multiple times.
func (f *fileSystem) Defer(fn func() error) {
	f.fsys.closer.Add(fn)
}

// Glob implements fs.GlobFS
func (f *fileSystem) Glob(pattern string) (matches []string, err error) {
	// Compile the pattern into a glob matcher
	matcher, err := glob.Compile(pattern)
	if err != nil {
		return nil, err
	}
	// Watch for changes to the pattern
	f.link.Select("glob", func(path string) bool {
		return matcher.Match(path)
	})
	// Base is a minor optimization to avoid walking the entire tree
	bases, err := glob.Bases(pattern)
	if err != nil {
		return nil, err
	}
	// Compute the matches for each base
	for _, base := range bases {
		results, err := f.glob(matcher, base)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		matches = append(matches, results...)
	}
	// Deduplicate the matches
	return orderedset.Strings(matches...), nil
}

// ReadDir implements fs.ReadDirFS
func (f *fileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	f.link.Select("readdir", func(path string) bool {
		return path == name || filepath.Dir(path) == name
	})
	des, err := fs.ReadDir(f.fsys, name)
	if err != nil {
		return nil, err
	}
	return des, nil
}

func (f *fileSystem) glob(matcher glob.Matcher, base string) (matches []string, err error) {
	// Walk the directory tree, filtering out non-valid paths
	err = fs.WalkDir(f.fsys, base, valid.WalkDirFunc(func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// If the paths match, add it to the list of matches
		if matcher.Match(path) {
			matches = append(matches, path)
		}
		return nil
	}))
	if err != nil {
		return nil, err
	}
	// return the list of matches
	return matches, nil
}

func relativePath(base, target string) string {
	rel := strings.TrimPrefix(target, base)
	if rel == "" {
		return "."
	} else if rel[0] == '/' {
		rel = rel[1:]
	}
	return rel
}
