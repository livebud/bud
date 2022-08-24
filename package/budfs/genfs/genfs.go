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
	tree.Generator = &fillerDir{tree}
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

func (d *dir) Mount(dirpath string, fsys fs.FS) {
	segments := strings.Split(dirpath, "/")
	last := len(segments) - 1
	node := mkdirAll(d.node, segments[:last])
	// Finally add the base path with it's file generator to the tree.
	child := node.Upsert(segments[last], fs.ModeDir, nil)
	child.Generator = Generate(func(target string) (fs.File, error) {
		return fsys.Open(target)
	})
}

func (d *dir) GenerateFile(path string, fn func(file *File) error) {
	segments := strings.Split(path, "/")
	last := len(segments) - 1
	node := mkdirAll(d.node, segments[:last])
	// Finally add the base path with it's file generator to the tree.
	child := node.Insert(segments[last], fs.FileMode(0), nil)
	child.Generator = &fileGenerator{child, fn}
}

func (d *dir) GenerateDir(path string, fn func(dir *Dir) error) {
	segments := strings.Split(path, "/")
	last := len(segments) - 1
	node := mkdirAll(d.node, segments[:last])
	// Finally add the base path with it's file generator to the tree.
	child := node.Upsert(segments[last], fs.ModeDir, nil)
	child.Generator = &dirGenerator{child, fn}
}

type File struct {
	node *treefs.Node
	Data []byte
}

func (f *File) Path() string {
	return f.node.Path()
}

func (f *File) Mode() fs.FileMode {
	return f.node.Mode
}

func (d *dir) open(target string) (fs.File, error) {
	// Find the closest match in the tree
	node, prefix, ok := d.node.FindByPrefix(target)
	if !ok {
		return nil, formatError(fs.ErrNotExist, "%q target not found in %q node", target, d.node.Path())
	}
	// File matches that aren't exact are not allowed.
	if prefix != target && node.Mode.IsRegular() {
		return nil, formatError(fs.ErrNotExist, "%q file generator doesn't match %q target", d.node.Path(), target)
	}
	// Run the generators
	relPath := relativePath(prefix, target)
	return node.Generator.Generate(relPath)
}

// mkdirAll creates all the parent directories leading up to the path
func mkdirAll(node *treefs.Node, segments []string) *treefs.Node {
	// Create the branches in the directory tree, if they don't exist already.
	for _, segment := range segments {
		child, ok := node.Child(segment)
		if !ok {
			child = node.Insert(segment, fs.ModeDir, nil)
			child.Generator = &fillerDir{child}
		}
		node = child
	}
	return node
}

type fileGenerator struct {
	child    *treefs.Node
	generate func(file *File) error
}

func (g *fileGenerator) Generate(target string) (fs.File, error) {
	file := &File{node: g.child}
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
	child    *treefs.Node
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
	return d.node.Mode
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
	dir := &Dir{&dir{g.child}, target}
	if err := g.generate(dir); err != nil {
		return nil, err
	}
	// When targeting directory generators directly, they are simply a filler
	// directory for the sub-generators.
	rel := relativePath(g.child.Path(), target)
	if rel == "." {
		generator := &fillerDir{g.child}
		return generator.Generate(target)
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

type fillerDir struct {
	node *treefs.Node
}

func (f *fillerDir) Generate(target string) (fs.File, error) {
	path := f.node.Path()
	// Filler directories must be exact matches with the target, otherwise we'll
	// create files that aren't supposed to exist.
	if target != "." {
		return nil, formatError(fs.ErrNotExist, "path doesn't match target in filler directory %s != %s", path, target)
	}
	children := f.node.Children()
	entries := make([]fs.DirEntry, len(children))
	for i, child := range children {
		entries[i] = &dirEntry{child}
	}
	return &virtual.Dir{
		Path:    path,
		Mode:    fs.ModeDir,
		Entries: entries,
	}, nil
}

type dirEntry struct {
	node *treefs.Node
}

var _ fs.DirEntry = (*dirEntry)(nil)

func (e *dirEntry) Name() string {
	return e.node.Name
}

func (e *dirEntry) IsDir() bool {
	return e.node.Mode.IsDir()
}

func (e *dirEntry) Type() fs.FileMode {
	return e.node.Mode
}

func (e *dirEntry) Info() (fs.FileInfo, error) {
	value := e.node.Generator
	if value == nil {
		value = &fillerDir{e.node}
	}
	file, err := value.Generate(".")
	if err != nil {
		return nil, err
	}
	return file.Stat()
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
