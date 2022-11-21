package trpipe

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/albertorestifo/dijkstra"
	"github.com/livebud/bud/package/log"
)

func New(log log.Log) *Pipeline {
	return &Pipeline{
		fns:   map[string][]func(file *File) error{},
		graph: dijkstra.Graph{},
		log:   log,
	}
}

type Pipeline struct {
	mu    sync.RWMutex
	fns   map[string][]func(file *File) error
	graph dijkstra.Graph
	log   log.Log
}

type File struct {
	path string
	ext  string
	Data []byte
}

func (f *File) Path() string {
	base := strings.TrimSuffix(f.path, filepath.Ext(f.path))
	return base + f.ext
}

func (p *Pipeline) Add(from, to string, fn func(file *File) error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.fns[from]; !ok {
		p.fns[from] = []func(file *File) error{}
	}
	key := from + ">" + to
	p.fns[key] = append(p.fns[key], fn)
	if p.graph[from] == nil {
		p.graph[from] = map[string]int{}
	}
	if p.graph[to] == nil {
		p.graph[to] = map[string]int{}
	}
	p.graph[from][to] = 1
}

func (p *Pipeline) Run(fromPath, toExt string, code []byte) ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	fromExt := path.Ext(fromPath)
	hops, _, err := p.graph.Path(fromExt, toExt)
	if err != nil {
		return nil, fmt.Errorf("trpipe: no path to transform %q to %q for %q. %w", fromExt, toExt, fromPath, err)
	} else if len(hops) == 0 {
		return code, nil
	}
	// Turn the hops into pairs (e.g. [ [.svelte, .js], ...])
	pairs := [][2]string{[2]string{hops[0], hops[0]}}
	for i := 1; i < len(hops); i++ {
		pairs = append(pairs, [2]string{hops[i-1], hops[i]})
		pairs = append(pairs, [2]string{hops[i], hops[i]})
	}
	file := &File{
		path: fromPath,
		ext:  fromExt,
		Data: code,
	}
	// Apply transformations over the transform pairs
	for _, pair := range pairs {
		// Handle .svelte -> .svelte transformations
		key := pair[0] + ">" + pair[1]
		if transforms, ok := p.fns[key]; ok {
			p.log.Fields(log.Fields{
				"key":        key,
				"transforms": len(transforms),
			}).Debug("trpipe: running transforms")
			for _, transform := range transforms {
				if err := transform(file); err != nil {
					return nil, err
				}
			}
			// Update the extension
			file.ext = pair[1]
		}
	}
	return file.Data, nil
}
