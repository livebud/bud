package dag2

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

func New() *Graph {
	return &Graph{
		nodes: map[string]struct{}{},
		ins:   map[string]map[string]struct{}{},
		outs:  map[string]map[string]struct{}{},
	}
}

type Graph struct {
	mu    sync.RWMutex
	nodes map[string]struct{}
	ins   map[string]map[string]struct{}
	outs  map[string]map[string]struct{}
}

func (g *Graph) Nodes() (nodes []string) {
	for path := range g.nodes {
		nodes = append(nodes, path)
	}
	sort.Strings(nodes)
	return nodes
}

func (g *Graph) Set(path string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.nodes[path] = struct{}{}
}

func (g *Graph) Link(from, to string) {
	// if from == to {
	// 	return
	// }
	g.mu.Lock()
	defer g.mu.Unlock()
	// Create nodes if we haven't yet
	g.nodes[from] = struct{}{}
	g.nodes[to] = struct{}{}
	// Set the dependencies
	if g.outs[from] == nil {
		g.outs[from] = map[string]struct{}{}
	}
	g.outs[from][to] = struct{}{}
	// Link the other way too
	if g.ins[to] == nil {
		g.ins[to] = map[string]struct{}{}
	}
	g.ins[to][from] = struct{}{}
}

// Remove a node
func (g *Graph) Remove(paths ...string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, path := range paths {
		g.remove(path)
	}
}

// Unlink and remove the node
func (g *Graph) remove(path string) {
	// Remove dependant links
	for from := range g.ins[path] {
		delete(g.ins[path], from)
		delete(g.outs[from], path)
	}
	// Remove dependency links
	for to := range g.outs[path] {
		delete(g.outs[path], to)
		delete(g.ins[to], path)
	}
	// Remove node
	delete(g.nodes, path)
}

// Return the links out (dependencies)
func (g *Graph) Children(from string) (tos []string) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.children(from)
}

func (g *Graph) children(from string) (froms []string) {
	for from := range g.outs[from] {
		froms = append(froms, from)
	}
	sort.Strings(froms)
	return froms
}

// Descendants recursively returns children, children of children, etc.
// Descendants includes self.
func (g *Graph) Descendants(path string) (descendants []string) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	descendants = append(descendants, path)
	descendants = append(descendants, g.children(path)...)
	return descendants
}

// Return the links in (parents)
func (g *Graph) Parents(to string) (froms []string) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.parents(to)
}

func (g *Graph) parents(to string) (tos []string) {
	for to := range g.ins[to] {
		tos = append(tos, to)
	}
	sort.Strings(tos)
	return tos
}

// Ancestors recursively returns parents, parents of parents, etc.
func (g *Graph) Ancestors(path string) (ancestors []string) {
	ancestors = append(ancestors, g.parents(path)...)
	return ancestors
}

// String returns a digraph
func (g *Graph) String() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	s := strings.Builder{}
	s.WriteString("digraph g {\n")
	// Rank from the bottom to the top
	s.WriteString("  rankdir = BT\n")
	for from := range g.nodes {
		tos := g.outs[from]
		if len(tos) == 0 {
			s.WriteString(fmt.Sprintf("  %q\n", from))
			continue
		}
		for to := range tos {
			s.WriteString(fmt.Sprintf("  %q -> %q\n", from, to))
		}
	}
	s.WriteString("}\n")
	return s.String()
}

func (g *Graph) ShortestPath(start, end string) (nodes []string, err error) {
	return g.shortestPath(start, end, nil)
}

func (g *Graph) shortestPath(start, end string, paths []string) (shortest []string, err error) {
	if _, ok := g.nodes[start]; !ok {
		return nil, fmt.Errorf("dag: %q doesn't exist", start)
	}
	paths = append(paths, start)
	if start == end {
		return paths, nil
	}
	for _, node := range g.children(start) {
		if hasNode(paths, node) {
			continue
		}
		newPath, err := g.shortestPath(node, end, paths)
		if err != nil {
			return nil, err
		}
		lenNewPath := len(newPath)
		lenShortest := len(shortest)
		if lenNewPath == 0 {
			continue
		}
		if lenShortest == 0 || lenNewPath < lenShortest {
			shortest = newPath
		}
	}
	if len(shortest) == 0 {
		return nil, fmt.Errorf("dag: no path between %q and %q", start, end)
	}
	return shortest, nil
}

func hasNode(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func (g *Graph) ShortestPathOf(start string, ends []string) (nodes []string, err error) {
	for _, end := range ends {
		shortest, err := g.shortestPath(start, end, nil)
		if err != nil {
			continue
		}
		if len(nodes) == 0 || len(shortest) < len(nodes) {
			nodes = shortest
		}
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("dag: no path between %q and %v", start, ends)
	}
	return nodes, nil
}
