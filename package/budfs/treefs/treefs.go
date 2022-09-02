package treefs

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/xlab/treeprint"
)

// type Node interface {
// 	Parent() Node
// 	Children() []Node
// 	Name() string
// 	Generate(target string) (fs.File, error)
// }

func New(name string) *Node {
	root := &Node{
		name:      name,
		mode:      fs.ModeDir,
		kind:      kindFiller,
		childMap:  map[string]*Node{},
		generator: nil,
	}
	root.path = computePath(root)
	root.generator = &fillerDir{root}
	return root
}

type Generator interface {
	Generate(target string) (fs.File, error)
}

type nodeKind uint8

const (
	kindFiller nodeKind = iota
	kindGenerator
)

type Node struct {
	path      string
	name      string
	mode      fs.FileMode
	kind      nodeKind
	parent    *Node
	childMap  map[string]*Node
	generator Generator
}

func computePath(n *Node) (path string) {
	if n == nil {
		return ""
	} else if n.parent == nil {
		return n.name
	}
	names := []string{n.name}
	parent := n.parent
	for parent != nil {
		names = append(names, parent.name)
		parent = parent.parent
	}
	last := len(names) - 2
	for i := last; i >= 0; i-- {
		if i < last {
			path += "/"
		}
		path += names[i]
	}
	return path
}

func (n *Node) Path() string {
	return n.path
}

func (n *Node) Mode() fs.FileMode {
	return n.mode
}

func (n *Node) Generate(target string) (fs.File, error) {
	return n.generator.Generate(target)
}

// Entry returns node as a directory entry.
func (n *Node) Entry() fs.DirEntry {
	return &dirEntry{n}
}

func (n *Node) child(name string) (*Node, bool) {
	child, found := n.childMap[name]
	return child, found
}

// Children returns a list of children, ordered alphanumerically.
func (n *Node) Children() (children []*Node) {
	children = make([]*Node, len(n.childMap))
	i := 0
	for _, child := range n.childMap {
		children[i] = child
		i++
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].name < children[j].name
	})
	return children
}

func (n *Node) InsertDir(path string, generator Generator) *Node {
	return n.insert(path, fs.ModeDir, generator)
}

func (n *Node) InsertFile(path string, generator Generator) *Node {
	return n.insert(path, fs.FileMode(0), generator)
}

func (n *Node) insert(path string, mode fs.FileMode, generator Generator) *Node {
	segments := strings.Split(path, "/")
	last := len(segments) - 1
	parent := n.mkdirAll(segments[:last])
	// Add the base path with it's file generator to the tree.
	child, found := parent.childMap[segments[last]]
	if !found {
		child = &Node{
			name:     segments[last],
			parent:   parent,
			childMap: map[string]*Node{},
		}
		child.path = computePath(child)
		parent.childMap[segments[last]] = child
	}
	// Create or update the child's attributes
	child.mode = mode
	child.kind = kindGenerator
	child.generator = generator
	return child
}

func (n *Node) mkdirAll(segments []string) *Node {
	parent := n
	// Create the branches in the directory tree, if they don't exist already.
	for _, segment := range segments {
		child, ok := parent.child(segment)
		if !ok {
			child = &Node{
				name:      segment,
				mode:      fs.ModeDir,
				kind:      kindFiller,
				parent:    parent,
				childMap:  map[string]*Node{},
				generator: nil,
			}
			child.path = computePath(child)
			child.generator = &fillerDir{child}
			parent.childMap[segment] = child
		}
		parent = child
	}
	return parent
}

func formatNode(node *Node) string {
	if node.kind == kindGenerator {
		return fmt.Sprintf("%s generator=%v mode=%s", node.name, node.generator, node.mode)
	}
	return fmt.Sprintf("%s mode=%s", node.name, node.mode)
}

// Print the nodes in the tree.
func (n *Node) Print() string {
	tp := treeprint.NewWithRoot(formatNode(n))
	n.print(tp)
	return tp.String()
}

func (n *Node) print(tp treeprint.Tree) {
	for _, child := range n.Children() {
		cp := tp.AddBranch(formatNode(child))
		child.print(cp)
	}
}

func (n *Node) Find(path string) (node *Node, found bool) {
	// Special case to find the root node
	if path == "." {
		return n, true
	}
	// Traverse the children keyed by segments
	node = n
	segments := strings.Split(path, "/")
	for _, name := range segments {
		node, found = node.childMap[name]
		if !found {
			return nil, false
		}
	}
	return node, true
}

func (n *Node) FindByPrefix(path string) (node *Node, prefix string, found bool) {
	// Special case to find the root node
	if path == "." {
		return n, path, true
	}
	// Traverse the children keyed by segments
	node = n
	names := strings.Split(path, "/")
	for i, name := range names {
		next, found := node.childMap[name]
		if !found {
			if i == 0 {
				return nil, "", false
			}
			// Find generator dirs that match the prefix
			node, nth, found := n.findAncestor(node, func(n *Node) bool {
				return n.kind == kindGenerator && n.mode.IsDir()
			})
			if !found {
				return nil, "", false
			}
			return node, strings.Join(names[:i-nth], "/"), true
		}
		node = next
	}
	// Try finding an ancestor generator
	parent, nth, found := n.findAncestor(node, func(n *Node) bool {
		return n.kind == kindGenerator
	})
	if !found {
		// Otherwise just return a filler node
		return node, path, true
	}
	return parent, strings.Join(names[:len(names)-nth], "/"), true
}

// Find the first ancestor that's a generator.
func (n *Node) findAncestor(child *Node, match func(n *Node) bool) (*Node, int, bool) {
	node := child
	nth := 0
	// Scope the search to the node itself to avoid potential infinite loops.
	for node != n && node != nil {
		if match(node) {
			return node, nth, true
		}
		node = node.parent
		nth++
	}
	return nil, nth, false
}

func (n *Node) Delete(path ...string) (node *Node, found bool) {
	var parent *Node
	node = n
	// Traverse the children keyed by segments
	for _, name := range path {
		child, found := node.childMap[name]
		if !found {
			return nil, false
		}
		parent = node
		node = child
	}
	// Path appears to be the root node, just do nothing
	if parent == nil {
		return nil, true
	}
	// Link the children of node to the parent
	for name, child := range node.childMap {
		parent.childMap[name] = child
		child.parent = parent
	}
	// Then delete the parent's link to node
	delete(parent.childMap, node.name)
	// Return the deleted node
	return node, true
}
