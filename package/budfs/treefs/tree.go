package treefs

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/xlab/treeprint"
)

type Generator interface {
	Generate(target string) (fs.File, error)
}

func New(name string) *Node {
	return &Node{
		Name:      name,
		Mode:      fs.ModeDir,
		Generator: nil,
		parent:    nil,
		children:  map[string]*Node{},
	}
}

type Node struct {
	Name      string
	Mode      fs.FileMode
	Generator Generator
	parent    *Node
	children  map[string]*Node
}

func (n *Node) Path() (path string) {
	if n == nil {
		return ""
	} else if n.parent == nil {
		return n.Name
	}
	names := []string{n.Name}
	parent := n.parent
	for parent != nil {
		names = append(names, parent.Name)
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

// Insert a new child.
func (n *Node) Insert(name string, mode fs.FileMode, generator Generator) (child *Node) {
	child = &Node{
		Name:      name,
		Mode:      mode,
		Generator: generator,
		parent:    n,
		children:  map[string]*Node{},
	}
	n.children[name] = child
	return child
}

func (n *Node) Upsert(name string, mode fs.FileMode, generator Generator) (child *Node) {
	child, found := n.children[name]
	if !found {
		return n.Insert(name, mode, generator)
	}
	child.Mode = mode
	child.Generator = generator
	return child
}

func (n *Node) Parent() *Node {
	return n.parent
}

func (n *Node) Child(name string) (*Node, bool) {
	child, found := n.children[name]
	return child, found
}

// Children returns a list of children, ordered alphanumerically.
func (n *Node) Children() (children []*Node) {
	children = make([]*Node, len(n.children))
	i := 0
	for _, child := range n.children {
		children[i] = child
		i++
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].Name < children[j].Name
	})
	return children
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
		node, found = node.children[name]
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
	segments := strings.Split(path, "/")
	for i, name := range segments {
		next, ok := node.children[name]
		if !ok {
			if i == 0 {
				return nil, "", false
			}
			return node, strings.Join(segments[:i], "/"), true
		}
		node = next
	}
	return node, path, true
}

func (n *Node) Delete(path ...string) (node *Node, found bool) {
	var parent *Node
	node = n
	// Traverse the children keyed by segments
	for _, name := range path {
		child, found := node.children[name]
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
	for name, child := range node.children {
		parent.children[name] = child
		child.parent = parent
	}
	// Then delete the parent's link to node
	delete(parent.children, node.Name)
	// Return the deleted node
	return node, true
}

func formatNode(node *Node) string {
	return fmt.Sprintf("%s generator=%v mode=%s", node.Name, node.Generator, node.Mode)
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
	return
}
