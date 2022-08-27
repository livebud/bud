package genfs

import (
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/livebud/bud/internal/virtual"
	"github.com/livebud/bud/package/budfs/treefs"
)

type Generator interface {
	Generate(target string) (fs.File, error)
}

type FileGenerator interface {
	GenerateFile(file *File) error
}

type DirGenerator interface {
	GenerateDir(dir *Dir) error
}

type Generate func(target string) (fs.File, error)

func (fn Generate) Generate(target string) (fs.File, error) {
	return fn(target)
}

func New() *FileSystem {
	tree := treefs.New(".")
	return &FileSystem{&dir{tree}}
}

type FileSystem struct {
	*dir
}

func (f *FileSystem) Open(target string) (fs.File, error) {
	if !fs.ValidPath(target) {
		return nil, formatError(fs.ErrInvalid, "invalid target path %q", target)
	}
	return f.dir.open(target)
}

type dir struct {
	node *treefs.Node
}

func (d *dir) Print() string {
	return d.node.Print()
}

func (d *dir) FileGenerator(path string, generator FileGenerator) {
	d.GenerateFile(path, generator.GenerateFile)
}

func (d *dir) DirGenerator(path string, generator DirGenerator) {
	d.GenerateDir(path, generator.GenerateDir)
}

func (d *dir) FileServer(path string, generator FileGenerator) {
	d.ServeFile(path, generator.GenerateFile)
}

func (d *dir) Mount(path string, fsys fs.FS) {
	d.node.InsertDir(path, Generate(func(target string) (fs.File, error) {
		return fsys.Open(target)
	}))
}

func (d *dir) GenerateFile(path string, fn func(file *File) error) {
	fileg := &fileGenerator{nil, fn}
	fileg.node = d.node.InsertFile(path, fileg)
}

func (d *dir) GenerateDir(path string, fn func(dir *Dir) error) {
	dirg := &dirGenerator{nil, fn}
	dirg.node = d.node.InsertDir(path, dirg)
}

func (d *dir) ServeFile(path string, fn func(file *File) error) {
	servef := &fileServer{nil, fn}
	servef.node = d.node.InsertDir(path, servef)
}

type File struct {
	path string
	mode fs.FileMode
	Data []byte
}

func (f *File) Path() string {
	return f.path
}

func (f *File) Mode() fs.FileMode {
	return f.mode
}

func (d *dir) open(target string) (fs.File, error) {
	// Find the closest match in the tree
	node, prefix, ok := d.node.FindByPrefix(target)
	if !ok {
		return nil, formatError(fs.ErrNotExist, "%q target not found in %q node", target, d.node.Path())
	}
	// File matches that aren't exact are not allowed.
	if prefix != target && node.Mode().IsRegular() {
		return nil, formatError(fs.ErrNotExist, "%q file generator doesn't match %q target", d.node.Path(), target)
	}
	// Run the generators
	relPath := relativePath(prefix, target)
	return node.Generate(relPath)
}

type fileGenerator struct {
	node     *treefs.Node
	generate func(file *File) error
}

func (g *fileGenerator) Generate(target string) (fs.File, error) {
	file := &File{
		path: g.node.Path(),
		mode: g.node.Mode(),
	}
	if err := g.generate(file); err != nil {
		return nil, formatError(err, "error generating %q", file.Path())
	}
	return &virtual.File{
		Path: file.Path(),
		Mode: file.Mode(),
		Data: file.Data,
	}, nil
}

type dirGenerator struct {
	node     *treefs.Node
	generate func(dir *Dir) error
}

type Dir struct {
	*dir
	target string
}

func (d *Dir) Path() string {
	return d.node.Path()
}

func (d *Dir) Mode() fs.FileMode {
	return d.node.Mode()
}

func (d *Dir) Target() string {
	return path.Join(d.Path(), d.target)
}

func (d *Dir) Relative() string {
	rel := strings.TrimPrefix(d.target, d.Path())
	if rel == "" {
		return "."
	} else if rel[0] == '/' {
		rel = rel[1:]
	}
	return rel
}

func (g *dirGenerator) Generate(target string) (fs.File, error) {
	// Call the generator function from the child
	dir := &Dir{&dir{g.node}, target}
	if err := g.generate(dir); err != nil {
		return nil, err
	}
	// When targeting directory generators directly, they are simply a virtual
	// directory that contains sub-generators.
	rel := relativePath(g.node.Path(), target)
	if rel == "." {
		children := g.node.Children()
		entries := make([]fs.DirEntry, len(children))
		for i, child := range children {
			entries[i] = child.Entry()
		}
		return &virtual.Dir{
			Path:    g.node.Path(),
			Mode:    g.node.Mode(),
			Entries: entries,
		}, nil
	}
	// Progress towards the target with the new branches in the child
	return dir.open(rel)
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

type EmbedFile struct {
	Data []byte
}

var _ FileGenerator = (*EmbedFile)(nil)

func (e *EmbedFile) GenerateFile(file *File) error {
	file.Data = e.Data
	return nil
}

func formatError(err error, format string, args ...interface{}) error {
	return fmt.Errorf("genfs: %s. %w", fmt.Sprintf(format, args...), err)
}

type fileServer struct {
	node     *treefs.Node
	generate func(file *File) error
}

func (g *fileServer) Generate(target string) (fs.File, error) {
	if target == "." {
		return nil, &fs.PathError{
			Op:   "open",
			Path: g.node.Path(),
			Err:  fs.ErrInvalid,
		}
	}
	file := &File{
		path: path.Join(g.node.Path(), target),
		mode: fs.FileMode(0),
	}
	if err := g.generate(file); err != nil {
		return nil, formatError(err, "error serving file %q", file.Path())
	}
	return &virtual.File{
		Path: file.Path(),
		Mode: file.Mode(),
		Data: file.Data,
	}, nil
}
