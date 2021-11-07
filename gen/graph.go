package gen

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/gobwas/glob"
)

func newGraph() *graph {
	return &graph{
		nodes: map[string]struct{}{},
		ins:   map[string]map[string]Event{},
		outs:  map[string]map[string]Event{},
	}
}

type graph struct {
	mu    sync.RWMutex
	nodes map[string]struct{}
	ins   map[string]map[string]Event
	outs  map[string]map[string]Event
}

func (g *graph) Link(from, to string, event Event) {
	if from == to {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	// Create nodes if we haven't yet
	g.nodes[from] = struct{}{}
	g.nodes[to] = struct{}{}
	// Set the dependencies
	if g.outs[from] == nil {
		g.outs[from] = map[string]Event{}
	}
	g.outs[from][to] = event
	// Link the other way too
	if g.ins[to] == nil {
		g.ins[to] = map[string]Event{}
	}
	g.ins[to][from] = event
}

// Return the links in (dependants)
func (g *graph) Ins(to string, event Event) (froms []string) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	tos := g.match(to)
	set := map[string]struct{}{}
	for _, to := range tos {
		for from, evt := range g.ins[to] {
			if evt&event > 0 {
				set[from] = struct{}{}
			}
		}
	}
	for from := range set {
		froms = append(froms, from)
	}
	sort.Strings(froms)
	return froms
}

// Match node patterns based on a path
// TODO: move glob compiling on write & improve performance
func (g *graph) match(path string) (nodes []string) {
	for node := range g.nodes {
		matcher := glob.MustCompile(node)
		if matcher.Match(path) {
			nodes = append(nodes, node)
		}
	}
	sort.Strings(nodes)
	return nodes
}

func (g *graph) String() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	s := strings.Builder{}
	s.WriteString("digraph g {\n")
	for from := range g.nodes {
		tos := g.outs[from]
		if len(tos) == 0 {
			s.WriteString(fmt.Sprintf("  %q\n", from))
			continue
		}
		for to, event := range tos {
			s.WriteString(fmt.Sprintf("  %q -> %q (%s)\n", from, to, event))
		}
	}
	s.WriteString("}\n")
	return s.String()
}

// Ins recursively returns dependants, dependants of dependants, etc.
// Ins includes path.
func (g *graph) DeepIns(path string, event Event) (des []string) {
	des = append(des, g.deepIns(path, event)...)
	return des
}

func (g *graph) deepIns(path string, event Event) (des []string) {
	froms := g.Ins(path, event)
	for _, from := range froms {
		des = append(des, from)
		des = append(des, g.deepIns(from, event)...)
	}
	return des
}
