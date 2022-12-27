package genfs

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/livebud/bud/internal/fsmode"
	"github.com/xlab/treeprint"
)

// generate function
// type generate func(target string) (fs.File, error)

// func (fn generate) Generate(target string) (fs.File, error) {
// 	return fn(target)
// }

func newTree() *tree {
	return &tree{
		node: &node{
			children:  map[string]*node{},
			Generator: nil,
			Path:      ".",
			Name:      ".",
			Mode:      fsmode.Dir,
			parent:    nil,
		},
	}
}

type tree struct {
	node *node
}

type node struct {
	Path string      // path from root
	Name string      // basename
	Mode fsmode.Mode // mode of the file

	Generator generator
	children  map[string]*node
	parent    *node
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

func (t *tree) Insert(fpath string, mode fsmode.Mode, generator generator) {
	segments := strings.Split(fpath, "/")
	last := len(segments) - 1
	name := segments[last]
	parent := t.mkdirAll(segments[:last])
	// Add the base path with it's file generator to the tree.
	child, found := parent.children[name]
	if !found {
		child = &node{
			children:  map[string]*node{},
			Generator: generator,
			Path:      fpath,
			Name:      name,
			Mode:      mode,
			parent:    parent,
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
	for _, segment := range segments {
		child, ok := parent.children[segment]
		if !ok {
			child = &node{
				children:  map[string]*node{},
				Generator: nil,
				Path:      path.Join(parent.Path, segment),
				Name:      segment,
				Mode:      fsmode.Dir,
				parent:    parent,
			}
			parent.children[segment] = child
		}
		parent = child
	}
	return parent
}

// Find an exact match the provided path
func (t *tree) Find(path string) (n *node, ok bool) {
	// Special case to find the root node
	if path == "." {
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
	return node, true
}

// Get the closest match to the provided path
func (t *tree) FindPrefix(path string) (n *node, ok bool) {
	// Special case to find the root node
	if path == "." {
		return t.node, true
	}
	// Traverse the children keyed by segments
	node := t.node
	segments := strings.Split(path, "/")
	for _, name := range segments {
		child, ok := node.children[name]
		if !ok {
			// nodes that aren't dirs must be an exact match
			if !node.Mode.IsDir() {
				return nil, false
			}
			return node, true
		}
		node = child
	}
	return node, true
}

func (t *tree) Delete(paths ...string) {
	for _, path := range paths {
		if node, ok := t.Find(path); ok {
			// We're trying to delete the root, ignore for now
			if node.parent == nil {
				continue
			}
			// Remove node from parent, deleting all descendants
			delete(node.parent.children, node.Name)
			node.parent = nil
		}
	}
}

func formatNode(node *node) string {
	if node.Generator == nil {
		return fmt.Sprintf("%s mode=%s", node.Name, node.Mode)
	}
	return fmt.Sprintf("%s mode=%s generator=%v", node.Name, node.Mode, node.Generator)
}

func (t *tree) Print() string {
	tp := treeprint.NewWithRoot(formatNode(t.node))
	print(tp, t.node)
	return tp.String()
}

func print(tp treeprint.Tree, node *node) {
	for _, child := range node.Children() {
		cp := tp.AddBranch(formatNode(child))
		print(cp, child)
	}
}
