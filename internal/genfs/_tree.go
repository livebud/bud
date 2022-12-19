package genfs

import (
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/xlab/treeprint"
)

func newTree() *tree {
	return &tree{}
}

type tree struct {
	node *node
}

func (t *tree) Insert(fpath string, mode fs.FileMode, generator generator) {
	segments := strings.Split(fpath, "/")
	last := len(segments) - 1
	name := segments[last]
	parent := t.mkdirAll(segments[:last])
	// Add the base path with it's file generator to the tree.
	child, found := parent.children[name]
	if !found {
		child = &node{
			generator,
			fpath,
			name,
			mode,
			map[string]*node{},
			parent,
		}
		// Add child to parent
		parent.children[name] = child
	}
	// Create or update the child's attributes
	child.Mode = mode
	child.Generator = generator
}

func (t *tree) mkdirAll(segments []string) *node {
	parent := t.node
	// Create the branches in the directory tree, if they don't exist already.
	for i, segment := range segments {
		child, ok := parent.children[segment]
		if !ok {
			child = &node{
				nil,
				path.Join(segments[0 : i+1]...),
				segment,
				fs.ModeDir,
				map[string]*node{},
				parent,
			}
			parent.children[segment] = child
		}
		parent = child
	}
	return parent
}

type node struct {
	Generator generator
	Path      string // path from root
	Name      string // node name
	Mode      fs.FileMode

	// Internal to the tree
	children map[string]*node
	parent   *node
}

// Print the nodes in the tree.
func (t *tree) Print() string {
	tp := treeprint.NewWithRoot(formatNode(t.node))
	print(tp, t.node)
	return tp.String()
}

func (t *tree) Get(path string, mode fs.FileMode) (n *node, ok bool) {
	// Special case to find the root node
	if path == "." && t.node.Mode&mode != 0 {
		return t.node, true
	}
	// Traverse the children keyed by segments
	node := t.node
	segments := strings.Split(path, "/")
	for _, name := range segments {
		node, ok = node.children[name]
		if !ok {
			return nil, false
		}
	}
	// Also match the mode
	if node.Mode&mode == 0 {
		return nil, false
	}
	return node, true
}

func (n *node) Generate(target string) (fs.File, error) {
	if n.Generator == nil {
		// n.Generator = &filler{n}
	}
	return n.Generator.Generate(target)
}

func (n *node) Children() (children []*node) {
	for _, child := range n.children {
		children = append(children, child)
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].Name < children[j].Name
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
	if node.Generator == nil {
		return fmt.Sprintf("%s mode=%s", node.Name, node.Mode)
	}
	return fmt.Sprintf("%s mode=%s generator=%v", node.Name, node.Mode, node.Generator)
}

func print(tp treeprint.Tree, node *node) {
	for _, child := range node.Children() {
		cp := tp.AddBranch(formatNode(child))
		print(cp, child)
	}
}
