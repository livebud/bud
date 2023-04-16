package transpiler

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/RyanCarrier/dijkstra"
)

// ErrNoPath is returned when there is no path between two extensions.
var ErrNoPath = dijkstra.ErrNoPath

type File struct {
	base string
	ext  string
	Data []byte
}

// Path returns the current file path that's being transpiled.
func (f *File) Path() string {
	return f.base + f.ext
}

// Interface for transpiling and testing if you can transpile from one extension
// to another. This interface is read-only. If you'd like to add extensions, use
// the Transpiler struct.
type Interface interface {
	Best(fromExt string, accepts []string) (string, error)
	Transpile(ctx context.Context, fromPath, toExt string, code []byte) ([]byte, error)
}

func New() *Transpiler {
	return &Transpiler{
		ids:   map[string]int{},
		exts:  map[int]string{},
		fns:   map[string][]func(ctx context.Context, file *File) error{},
		graph: dijkstra.NewGraph(),
	}
}

// Transpiler is a generic multi-step tool for transpiling code from one
// language to another.
type Transpiler struct {
	ids  map[string]int // ext -> id
	exts map[int]string // id -> ext

	// map["ext>ext"][]fns
	fns map[string][]func(ctx context.Context, file *File) error

	mu    sync.RWMutex
	graph *dijkstra.Graph
}

var _ Interface = (*Transpiler)(nil)

// edgekey returns a key for the edge between two extensions.
// (e.g. edgeKey("svelte", "html") => "svelte>html")
func edgeKey(fromExt, toExt string) string {
	return fromExt + ">" + toExt
}

// Add a tranpile function to go from one extension to another.
func (t *Transpiler) Add(fromExt, toExt string, transpile func(ctx context.Context, file *File) error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.add(fromExt, toExt, transpile)
}

func (t *Transpiler) add(fromExt, toExt string, transpile func(ctx context.Context, file *File) error) {
	// Add the "from" extension to the graph
	if _, ok := t.ids[fromExt]; !ok {
		id := len(t.ids)
		t.ids[fromExt] = id
		t.exts[id] = fromExt
		t.graph.AddVertex(id)
	}
	edge := edgeKey(fromExt, toExt)
	// If the "from" and "to" extensions are the same, add the function and return
	if fromExt == toExt {
		t.fns[edge] = append(t.fns[edge], transpile)
		return
	}
	// Add the "to" extension to the graph
	if _, ok := t.ids[toExt]; !ok {
		id := len(t.ids)
		t.ids[toExt] = id
		t.exts[id] = toExt
		t.graph.AddVertex(id)
	}
	// Add the edge with a cost of 1
	t.graph.AddArc(t.ids[fromExt], t.ids[toExt], 1)
	// Add the function to a list of transpilers
	t.fns[edge] = append(t.fns[edge], transpile)
}

func (t *Transpiler) Path(fromExt, toExt string) (hops []string, err error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.path(fromExt, toExt)
}

// Best returns the best extension to transpile to from the given extension
func (t *Transpiler) Best(fromExt string, accepts []string) (toExt string, err error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.best(fromExt, accepts)
}

// Candidate is a candidate extension we can transpile to with the number of
// hops it would take to get there.
// type candidate struct {
// 	Ext  string
// 	Hops int
// }

func (t *Transpiler) best(fromExt string, accepts []string) (string, error) {
	fromID, ok := t.ids[fromExt]
	if !ok {
		return "", fmt.Errorf("transpiler: %w for %q", ErrNoPath, fromExt)
	}
	buckets := map[int][]string{}
	for _, id := range t.ids {
		if id == fromID {
			continue
		}
		best, err := t.graph.Shortest(fromID, id)
		if err != nil {
			if errors.Is(err, ErrNoPath) {
				continue
			}
			return "", fmt.Errorf("transpiler: unable to get shorted path for %q. %w", fromExt, err)
		}
		distance := int(best.Distance)
		buckets[distance] = append(buckets[distance], t.exts[id])
	}
	distances := []int{}
	for distance := range buckets {
		distances = append(distances, distance)
	}
	sort.Ints(distances)
	hops := [][]string{}
	for _, distance := range distances {
		hops = append(hops, buckets[distance])
	}
	// Within each hop bucket, look for the most acceptable extension
	for _, exts := range hops {
		for _, accept := range accepts {
			for _, ext := range exts {
				if ext == accept {
					return ext, nil
				}
			}
		}
	}
	return "", fmt.Errorf("transpiler: no acceptable path for %q. %w", fromExt, ErrNoPath)
}

// Path to go from one extension to another.
func (t *Transpiler) path(fromExt, toExt string) (hops []string, err error) {
	if fromExt == toExt {
		return []string{fromExt}, nil
	}
	if _, ok := t.ids[fromExt]; !ok {
		return nil, fmt.Errorf("transpiler: %w from %q to %q", ErrNoPath, fromExt, toExt)
	}
	if _, ok := t.ids[toExt]; !ok {
		return nil, fmt.Errorf("transpiler: %w from %q to %q", ErrNoPath, fromExt, toExt)
	}
	best, err := t.graph.Shortest(t.ids[fromExt], t.ids[toExt])
	if err != nil {
		return nil, fmt.Errorf("transpiler: %w", err)
	}
	for _, id := range best.Path {
		hops = append(hops, t.exts[id])
	}
	return hops, nil
}

func (t *Transpiler) Transpile(ctx context.Context, fromPath, toExt string, code []byte) ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.transpile(ctx, fromPath, toExt, code)
}

// Transpile the code from one extension to another.
func (t *Transpiler) transpile(ctx context.Context, fromPath, toExt string, code []byte) ([]byte, error) {
	fromExt := filepath.Ext(fromPath)
	// Find the shortest path
	hops, err := t.path(fromExt, toExt)
	if err != nil {
		return nil, err
	}
	// Create the file
	file := &File{
		base: strings.TrimSuffix(fromPath, filepath.Ext(fromPath)),
		ext:  fromExt,
		Data: code,
	}
	// For each hop run the functions
	for i, ext := range hops {
		// Call the transition functions (e.g. svelte => html)
		if i > 0 {
			prevExt := hops[i-1]
			edge := edgeKey(prevExt, ext)
			for _, fn := range t.fns[edge] {
				if err := fn(ctx, file); err != nil {
					return nil, err
				}
			}
		}
		file.ext = ext
		// Call the loops (e.g. svelte => svelte)
		for _, fn := range t.fns[edgeKey(ext, ext)] {
			if err := fn(ctx, file); err != nil {
				return nil, err
			}
		}
	}
	return file.Data, nil
}
