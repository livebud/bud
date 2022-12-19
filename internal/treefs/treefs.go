package treefs

import (
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/livebud/bud/package/virtual"
	"github.com/xlab/treeprint"
)

func New(path string) *Tree {
	return &Tree{&node{map[string]*node{}, nil, path, path, ModeDir, nil}}
}

type Generate func(target string) (fs.File, error)

func (fn Generate) Generate(target string) (fs.File, error) {
	return fn(target)
}

type Generator interface {
	Generate(target string) (fs.File, error)
}

type Node interface {
	Name() string
	Path() string
	Mode() Mode
	Generate(target string) (fs.File, error)
	Children() (entries []Node)
}

type Mode uint8

func (m Mode) String() string {
	out := ""
	if m&ModeDir != 0 {
		out += "d"
	} else {
		out += "-"
	}
	if m&ModeGenerator != 0 {
		out += "g"
	} else {
		out += "-"
	}
	return out
}

func (m Mode) IsDir() bool {
	return m&ModeDir != 0
}

func (m Mode) IsFile() bool {
	return m&ModeDir == 0
}

func (m Mode) IsGenerator() bool {
	return m&ModeGenerator != 0
}

func (m Mode) FileMode() fs.FileMode {
	mode := fs.FileMode(0)
	if m.IsDir() {
		mode |= fs.ModeDir
	}
	return mode
}

const (
	ModeDir Mode = 1 << iota
	ModeGenerator
)

type Tree struct {
	node *node
}

// TODO: reject inserting Mode(0) with children
func (t *Tree) Insert(fpath string, mode Mode, generator Generator) error {
	segments := strings.Split(fpath, "/")
	last := len(segments) - 1
	name := segments[last]
	parent := t.mkdirAll(segments[:last])
	// Add the base path with it's file generator to the tree.
	child, found := parent.childMap[name]
	if !found {
		child = &node{
			map[string]*node{},
			generator,
			fpath,
			name,
			mode,
			parent,
		}
		// Add child to parent
		parent.childMap[name] = child
	}
	// Create or update the child's attributes
	child.mode = mode
	child.generator = generator
	return nil
}

func (t *Tree) mkdirAll(segments []string) *node {
	parent := t.node
	// Create the branches in the directory tree, if they don't exist already.
	for i, segment := range segments {
		child, ok := parent.childMap[segment]
		if !ok {
			child = &node{
				map[string]*node{},
				nil,
				path.Join(segments[0 : i+1]...),
				segment,
				ModeDir,
				parent,
			}
			parent.childMap[segment] = child
		}
		parent = child
	}
	return parent
}

// Print the nodes in the tree.
func (t *Tree) Print() string {
	tp := treeprint.NewWithRoot(formatNode(t.node))
	print(tp, t.node)
	return tp.String()
}

func (t *Tree) get(path string) (n *node, ok bool) {
	// Special case to find the root node
	if path == "." {
		return t.node, true
	}
	// Traverse the children keyed by segments
	node := t.node
	segments := strings.Split(path, "/")
	for _, name := range segments {
		node, ok = node.childMap[name]
		if !ok {
			return nil, false
		}
	}
	return node, true
}

func (t *Tree) Get(path string) (n Node, ok bool) {
	return t.get(path)
}

func (t *Tree) FindByPrefix(fpath string) (n Node, ok bool) {
	// Found exact match
	if node, ok := t.Get(fpath); ok {
		return node, true
	}
	// Found prefix that matches the mode
	for isWithin(t.node.path, fpath) {
		if node, ok := t.Get(fpath); ok {
			if !node.Mode().IsDir() {
				return nil, false
			}
			return node, true
		}
		next := path.Dir(fpath)
		// Ensure we're not at the root
		if next == fpath {
			return nil, false
		}
		fpath = next
	}
	return nil, false
}

func (t *Tree) Delete(paths ...string) {
	for _, path := range paths {
		if node, ok := t.get(path); ok {
			delete(node.parent.childMap, node.name)
		}
	}
}

type node struct {
	childMap  map[string]*node
	generator Generator
	path      string // path from root
	name      string // node name
	mode      Mode
	parent    *node
}

var _ Node = (*node)(nil)

func (n *node) Name() string {
	return n.name
}

func (n *node) Path() string {
	return n.path
}

func (n *node) Mode() Mode {
	return n.mode
}

func (n *node) Generate(target string) (fs.File, error) {
	if n.generator == nil {
		n.generator = &filler{n}
	}
	return n.generator.Generate(target)
}

func (n *node) children() (children []*node) {
	for _, child := range n.childMap {
		children = append(children, child)
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].Name() < children[j].Name()
	})
	return children
}

func (n *node) Children() (children []Node) {
	for _, child := range n.childMap {
		children = append(children, child)
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].Name() < children[j].Name()
	})
	return children
}

func isWithin(parent, child string) bool {
	if parent == "." {
		// Everything is within the root
		return true
	}
	if parent == child {
		return true
	}
	return strings.HasPrefix(child, parent)
}

func formatNode(node *node) string {
	if node.generator == nil {
		return fmt.Sprintf("%s mode=%s", node.name, node.mode)
	}
	return fmt.Sprintf("%s mode=%s generator=%v", node.name, node.mode, node.generator)
}

func print(tp treeprint.Tree, node *node) {
	for _, child := range node.children() {
		cp := tp.AddBranch(formatNode(child))
		print(cp, child)
	}
}

// filler is a generator that creates virtual directories so that tree.Open()
// works like a filesystem.
type filler struct {
	node *node
}

func (f *filler) Generate(target string) (fs.File, error) {
	path := f.node.Path()
	// Filler directories must be exact matches with the target, otherwise we'll
	// create files that aren't supposed to exist.
	if target != path {
		return nil, fmt.Errorf("treefs: path doesn't match target in filler directory %s != %s", path, target)
	}
	children := f.node.children()
	var entries []fs.DirEntry
	for _, child := range children {
		de := &dirEntry{child}
		entries = append(entries, de)
	}
	return virtual.New(&virtual.Dir{
		Path:    path,
		Mode:    fs.ModeDir,
		Entries: entries,
	}), nil
}

type dirEntry struct {
	node *node
}

var _ fs.DirEntry = (*dirEntry)(nil)

func (e *dirEntry) Name() string {
	return e.node.name
}

func (e *dirEntry) IsDir() bool {
	return e.node.mode.IsDir()
}

func (e *dirEntry) Type() fs.FileMode {
	return e.node.mode.FileMode()
}

func (e *dirEntry) Info() (fs.FileInfo, error) {
	value := e.node.generator
	if value == nil {
		value = &filler{e.node}
	}
	file, err := value.Generate(e.node.Path())
	if err != nil {
		return nil, err
	}
	return file.Stat()
}
