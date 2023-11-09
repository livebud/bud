package radix

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/livebud/bud/pkg/mux/ast"
	"github.com/livebud/bud/pkg/mux/internal/parser"
)

var ErrDuplicate = fmt.Errorf("route")
var ErrNoMatch = fmt.Errorf("no match")

func New() *Tree {
	return &Tree{}
}

type Tree struct {
	root *Node
}

func (t *Tree) Insert(route string, handler http.Handler) error {
	r, err := parser.Parse(trimTrailingSlash(route))
	if err != nil {
		return err
	}
	// Expand optional and wildcard routes
	for _, route := range r.Expand() {
		if err := t.insert(route, handler); err != nil {
			return err
		}
	}
	return nil
}

func (t *Tree) insert(route *ast.Route, handler http.Handler) error {
	if t.root == nil {
		t.root = &Node{
			Route:    route,
			Handler:  handler,
			sections: route.Sections,
		}
		return nil
	}
	return t.root.insert(route, handler, route.Sections)
}

type Node struct {
	Route    *ast.Route
	Handler  http.Handler
	sections ast.Sections
	children Nodes
}

func (n *Node) Priority() (priority int) {
	if len(n.sections) == 0 {
		return 0
	}
	return n.sections[0].Priority()
}

type Nodes []*Node

var _ sort.Interface = (*Nodes)(nil)

func (n Nodes) Len() int {
	return len(n)
}

func (n Nodes) Less(i, j int) bool {
	return n[i].Priority() > n[j].Priority()
}

func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n *Node) insert(route *ast.Route, handler http.Handler, sections ast.Sections) error {
	lcp := n.sections.LongestCommonPrefix(sections)
	if lcp < n.sections.Len() {
		// Split the node's sections
		parts := n.sections.Split(lcp)
		// Create a new node with the parent's sections after the lcp.
		splitChild := &Node{
			Route:    n.Route,
			Handler:  n.Handler,
			sections: parts[1],
			children: n.children,
		}
		n.sections = parts[0]
		n.children = Nodes{splitChild}
		// Add a new child if we have more sections left.
		if lcp < sections.Len() {
			newChild := &Node{
				Route:    route,
				Handler:  handler,
				sections: sections.Split(lcp)[1],
			}
			// Replace the parent's sections with the lcp.
			n.children = append(n.children, newChild)
			n.Route = nil
			n.Handler = nil
		} else {
			// Otherwise this route matches the parent. Update the parent's route and
			// handler.
			n.Route = route
			n.Handler = handler
		}
		sort.Sort(n.children)
		return nil
	}
	// Route already exists
	if lcp == sections.Len() {
		if n.Route == nil {
			n.Route = route
			n.Handler = handler
			return nil
		}
		oldRoute := n.Route.String()
		newRoute := route.String()
		if oldRoute == newRoute {
			return fmt.Errorf("%w already exists %q", ErrDuplicate, oldRoute)
		}
		return fmt.Errorf("%w %q is ambiguous with %q", ErrDuplicate, newRoute, oldRoute)
	}
	// Check children for a match
	remainingSections := sections.Split(lcp)[1]
	for _, child := range n.children {
		if child.sections.At(0) == remainingSections.At(0) {
			return child.insert(route, handler, remainingSections)
		}
	}
	n.children = append(n.children, &Node{
		Route:    route,
		Handler:  handler,
		sections: remainingSections,
	})
	sort.Sort(n.children)
	return nil
}

type Slot struct {
	Key   string
	Value string
}

func createSlots(route *ast.Route, slotValues []string) (slots []*Slot) {
	index := 0
	for _, section := range route.Sections {
		switch s := section.(type) {
		case ast.Slot:
			slots = append(slots, &Slot{
				Key:   s.Slot(),
				Value: slotValues[index],
			})
			index++
		}
	}
	return slots
}

type Match struct {
	Route   *ast.Route
	Path    string
	Slots   []*Slot
	Handler http.Handler
}

func (m *Match) String() string {
	s := new(strings.Builder)
	s.WriteString(m.Route.String())
	query := url.Values{}
	for _, slot := range m.Slots {
		query.Add(slot.Key, slot.Value)
	}
	if len(query) > 0 {
		s.WriteString(" ")
		s.WriteString(query.Encode())
	}
	return s.String()
}

func (t *Tree) Match(path string) (*Match, error) {
	path = trimTrailingSlash(path)
	// A tree without any routes shouldn't panic
	if t.root == nil || len(path) == 0 || path[0] != '/' {
		return nil, fmt.Errorf("%w for %q", ErrNoMatch, path)
	}
	match, ok := t.root.Match(path, []string{})
	if !ok {
		return nil, fmt.Errorf("%w for %q", ErrNoMatch, path)
	}
	match.Path = path
	return match, nil
}

func (n *Node) Match(path string, slotValues []string) (*Match, bool) {
	for _, section := range n.sections {
		if len(path) == 0 {
			return nil, false
		}
		index, slots := section.Match(path)
		if index <= 0 {
			return nil, false
		}
		path = path[index:]
		slotValues = append(slotValues, slots...)
	}
	if len(path) == 0 {
		// We've reached a non-routable node
		if n.Route == nil {
			return nil, false
		}
		return &Match{
			Route:   n.Route,
			Handler: n.Handler,
			Slots:   createSlots(n.Route, slotValues),
		}, true
	}
	for _, child := range n.children {
		if match, ok := child.Match(path, slotValues); ok {
			return match, true
		}
	}
	return nil, false
}

// Find by a route
func (t *Tree) Find(route string) (*Node, error) {
	r, err := parser.Parse(trimTrailingSlash(route))
	if err != nil {
		return nil, err
	} else if t.root == nil {
		return nil, fmt.Errorf("%w for %s", ErrNoMatch, route)
	}
	return t.root.find(route, r.Sections)
}

// Find by a route
func (n *Node) find(route string, sections ast.Sections) (*Node, error) {
	lcp := n.sections.LongestCommonPrefix(sections)
	if lcp < n.sections.Len() {
		return nil, fmt.Errorf("%w for %s", ErrNoMatch, route)
	}
	if lcp == sections.Len() {
		if n.Route == nil {
			return nil, fmt.Errorf("%w for %s", ErrNoMatch, route)
		}
		return n, nil
	}
	remainingSections := sections.Split(lcp)[1]
	for _, child := range n.children {
		if child.sections.At(0) == remainingSections.At(0) {
			return child.find(route, remainingSections)
		}
	}
	return nil, fmt.Errorf("%w for %s", ErrNoMatch, route)
}

func (t *Tree) FindByPrefix(prefix string) (*Node, error) {
	route, err := parser.Parse(trimTrailingSlash(prefix))
	if err != nil {
		return nil, err
	} else if t.root == nil {
		return nil, fmt.Errorf("%w for %s", ErrNoMatch, route)
	}
	return t.root.findByPrefix(prefix, route.Sections)
}

func (n *Node) findByPrefix(prefix string, sections ast.Sections) (*Node, error) {
	if n.sections.Len() > sections.Len() {
		return nil, fmt.Errorf("%w for %s", ErrNoMatch, prefix)
	}
	lcp := n.sections.LongestCommonPrefix(sections)
	if lcp == 0 {
		return nil, fmt.Errorf("%w for %s", ErrNoMatch, sections)
	}
	if lcp == sections.Len() {
		if n.Route == nil {
			return nil, fmt.Errorf("%w for %s", ErrNoMatch, prefix)
		}
		return n, nil
	}
	remainingSections := sections.Split(lcp)[1]
	for _, child := range n.children {
		if child.sections.At(0) == remainingSections.At(0) && child.sections.Len() <= remainingSections.Len() {
			return child.findByPrefix(prefix, remainingSections)
		}
	}
	if lcp < n.sections.Len() {
		return nil, fmt.Errorf("%w for %s", ErrNoMatch, prefix)
	}
	return n, nil
}

func (t *Tree) String() string {
	return t.string(t.root, "")
}

func (t *Tree) string(n *Node, indent string) string {
	route := n.sections.String()
	var mods []string
	if n.Route != nil {
		mods = append(mods, "routable="+n.Route.String())
	}
	mod := ""
	if len(mods) > 0 {
		mod = " [" + strings.Join(mods, ", ") + "]"
	}
	out := fmt.Sprintf("%s%s%s\n", indent, route, mod)
	for i := 0; i < len(route); i++ {
		indent += "•"
	}
	for _, child := range n.children {
		out += t.string(child, indent)
	}
	return out
}

// Traverse the tree in depth-first order
func (t *Tree) Each(fn func(n *Node) (next bool)) {
	t.each(t.root, fn)
}

func (t *Tree) each(n *Node, fn func(n *Node) (next bool)) {
	if !fn(n) {
		return
	}
	for _, child := range n.children {
		t.each(child, fn)
	}
}

// trimTrailingSlash strips any trailing slash (e.g. /users/ => /users)
func trimTrailingSlash(input string) string {
	if input == "/" {
		return input
	}
	return strings.TrimRight(input, "/")
}
