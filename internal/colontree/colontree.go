package colontree

import (
	"sort"
	"strconv"
	"strings"
)

func New(path, value string) *Node {
	return &Node{
		parent:   nil,
		path:     path,
		children: map[string]*Node{},
		value:    value,
	}
}

type Node struct {
	parent   *Node
	path     string
	children map[string]*Node
	value    string
}

func (n *Node) Insert(path, value string) *Node {
	if path == "" {
		n.value = value
		return n
	}
	parts := strings.Split(path, ":")
	child, ok := n.children[parts[0]]
	if !ok {
		child = &Node{
			parent:   n,
			path:     parts[0],
			children: map[string]*Node{},
		}
		n.children[parts[0]] = child
	} else if len(parts) == 1 {
		child.value = value
	}
	return child.Insert(strings.Join(parts[1:], ":"), value)
}

func (n *Node) Find(path string) *Node {
	if path == "" {
		return n
	}
	parts := strings.Split(path, ":")
	child, ok := n.children[parts[0]]
	if !ok {
		return nil
	}
	return child.Find(strings.Join(parts[1:], ":"))
}

func (n *Node) Parent() *Node {
	return n.parent
}

func (n *Node) Path() string {
	return n.path
}

// Full path
func (n *Node) Full() string {
	if n.parent == nil || n.parent.path == "" {
		return n.path
	}
	return n.parent.Full() + ":" + n.path
}

func (n *Node) Value() string {
	return n.value
}

func (n *Node) Children() []*Node {
	children := []*Node{}
	for _, child := range n.children {
		children = append(children, child)
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].path < children[j].path
	})
	return children
}

// Print in graphviz dot format
func (n *Node) Print() string {
	out := "digraph {\n"
	out += print(n, "  ")
	out += "}"
	return out
}

func print(n *Node, prefix string) string {
	out := ""
	quotedPath := strconv.Quote(n.path)
	out += prefix + quotedPath + ";\n"
	for _, child := range n.Children() {
		out += prefix + quotedPath + " -> " + strconv.Quote(child.path) + ";\n"
		out += print(child, prefix)
	}
	return out
}
