package genfs

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/livebud/bud/internal/genfs/fstree"
	"github.com/livebud/bud/internal/virtual"
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

func New() *Dir {
	tree := fstree.New(".")
	tree.Generator = &fillerDir{tree}
	// TODO: use a filesystem instead
	return &Dir{tree, ""}
}

type Dir struct {
	node   *fstree.Node
	target string
}

func (d *Dir) Target() string {
	return d.target
}

func (d *Dir) Path() string {
	return d.node.Path()
}

func (d *Dir) Rel() string {
	rel := strings.TrimPrefix(d.target, d.Path())
	if rel == "" {
		return "."
	} else if rel[0] == '/' {
		rel = rel[1:]
	}
	return rel
}

func (d *Dir) Mode() fs.FileMode {
	return d.node.Mode
}

type File struct {
	node *fstree.Node
	Data []byte
}

func (f *File) Path() string {
	return f.node.Path()
}

func (f *File) Mode() fs.FileMode {
	return f.node.Mode
}

func (d *Dir) Open(target string) (fs.File, error) {
	// Find the closest match in the tree
	node, prefix, ok := d.node.FindByPrefix(target)
	if !ok {
		return nil, fs.ErrNotExist
	}
	// File matches that aren't exact are not allowed.
	if prefix != target && node.Mode.IsRegular() {
		return nil, fs.ErrNotExist
	}
	// Run the generators
	return node.Generator.Generate(target)
}

func (d *Dir) FileGenerator(filepath string, generator FileGenerator) {
	d.GenerateFile(filepath, generator.GenerateFile)
}

func (d *Dir) DirGenerator(dirpath string, generator DirGenerator) {
	d.GenerateDir(dirpath, generator.GenerateDir)
}

func (d *Dir) GenerateFile(path string, fn func(file *File) error) {
	segments := strings.Split(path, "/")
	last := len(segments) - 1
	node := mkdirAll(d.node, segments[:last])
	// Finally add the base path with it's file generator to the tree.
	child := node.Insert(segments[last], fs.FileMode(0), nil)
	child.Generator = &fileGenerator{child, fn}
}

func (d *Dir) GenerateDir(path string, fn func(dir *Dir) error) {
	segments := strings.Split(path, "/")
	last := len(segments) - 1
	node := mkdirAll(d.node, segments[:last])
	// Finally add the base path with it's file generator to the tree.
	child := node.Insert(segments[last], fs.ModeDir, nil)
	child.Generator = &dirGenerator{child, fn}
}

// mkdirAll creates all the parent directories leading up to the path
func mkdirAll(node *fstree.Node, segments []string) *fstree.Node {
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
	child    *fstree.Node
	generate func(file *File) error
}

func (g *fileGenerator) Generate(target string) (fs.File, error) {
	file := &File{node: g.child}
	if err := g.generate(file); err != nil {
		return nil, formatError(err, "error generating %q", target)
	}
	return &virtual.File{
		Path: file.Path(),
		Mode: file.Mode(),
		Data: file.Data,
	}, nil
}

type dirGenerator struct {
	child    *fstree.Node
	generate func(dir *Dir) error
}

func (g *dirGenerator) Generate(target string) (fs.File, error) {
	// Call the generator function from the child
	dir := &Dir{g.child, target}
	if err := g.generate(dir); err != nil {
		return nil, err
	}
	// When targeting directory generators directly, they are simply a filler
	// directory for the sub-generators.
	if g.child.Path() == target {
		generator := &fillerDir{g.child}
		return generator.Generate(target)
	}
	// Progress towards the target with the new branches in the child
	return dir.Open(target)
}

func relativeTarget(base, target string) string {
	rel := strings.TrimPrefix(target, base)
	if rel == "" {
		return "."
	} else if rel[0] == '/' {
		rel = rel[1:]
	}
	return rel
}

type fillerDir struct {
	node *fstree.Node
}

func (f *fillerDir) Generate(target string) (fs.File, error) {
	path := f.node.Path()
	// Filler directories must be exact matches with the target, otherwise we'll
	// create files that aren't supposed to exist.
	if path != target {
		return nil, fs.ErrNotExist
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
	node *fstree.Node
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
	file, err := value.Generate(e.node.Path())
	if err != nil {
		return nil, err
	}
	return file.Stat()
}

func formatError(err error, format string, args ...interface{}) error {
	return fmt.Errorf("budfs: %s. %w", fmt.Sprintf(format, args...), err)
}
