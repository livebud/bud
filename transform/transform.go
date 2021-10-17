package transform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/go-duo/bud/internal/dag"
)

type File struct {
	path string
	ext  string
	Code []byte
}

func (f *File) Path() string {
	base := strings.TrimSuffix(f.path, filepath.Ext(f.path))
	return base + f.ext
}

// Platform we're transforming to. Based on esbuild's Platform settings.
type Platform int

const (
	PlatformNeutral Platform = iota
	PlatformBrowser
	PlatformNode
)

type Map map[Platform]func(file *File) error

type Transformable struct {
	To   string
	From string
	Map  Map
}

func MustLoad(transformables ...*Transformable) *Transformer {
	transformer, err := Load(transformables...)
	if err != nil {
		panic("transform: unable to load the transformer: " + err.Error())
	}
	return transformer
}

func Load(transformables ...*Transformable) (*Transformer, error) {
	browser, err := load(PlatformBrowser, transformables)
	if err != nil {
		return nil, err
	}
	node, err := load(PlatformNode, transformables)
	if err != nil {
		return nil, err
	}
	return &Transformer{browser, node}, nil
}

func getTransform(transformable *Transformable, platform Platform) (func(file *File) error, bool) {
	tr, ok := transformable.Map[platform]
	if ok {
		return tr, true
	}
	tr, ok = transformable.Map[PlatformNeutral]
	if ok {
		return tr, true
	}
	return nil, false
}

func load(platform Platform, transformables []*Transformable) (*transformer, error) {
	graph := dag.New()
	tmap := map[string][]func(file *File) error{}
	froms := map[string]struct{}{}
	// Build a dependency graph of how the transforms transform (from -> to)
	for _, transformable := range transformables {
		transform, ok := getTransform(transformable, platform)
		if !ok {
			continue
		}
		graph.Link(transformable.From, transformable.To)
		key := transformable.From + ">" + transformable.To
		froms[transformable.From] = struct{}{}
		// We can compose transforms of the same type. For example, two
		// svelte-to-svelte transforms. We cannot compose different types though.
		// For example, two svelte-to-jsx transforms.
		// TODO: Figure out what to do with the ignored transform.
		if len(tmap[key]) > 0 && transformable.From != transformable.To {
			continue
		}
		// TODO: this can be undefined
		tmap[key] = append(tmap[key], transform)
	}
	// Build the full pathmap to generate the plugins
	pathmap := map[string]string{}
	for from := range froms {
		to, err := graph.ShortestPathOf(from, []string{".js", ".jsx"})
		if err != nil {
			return nil, err
		}
		pathmap[from] = to[len(to)-1]
	}
	// Build the index to efficiently access handlers
	// Compose multiple transforms
	index := map[string]func(file *File) error{}
	for key, transforms := range tmap {
		index[key] = compose(transforms)
	}
	return &transformer{graph, index, pathmap}, nil
}

func compose(fns []func(file *File) error) func(file *File) error {
	return func(file *File) error {
		for _, fn := range fns {
			if err := fn(file); err != nil {
				return err
			}
		}
		return nil
	}
}

// Transformer aggregates all the platform-specific transformers
type Transformer struct {
	Browser *transformer
	Node    *transformer
}

// Transformer is specific to a platform
type transformer struct {
	graph   *dag.Graph
	index   map[string]func(file *File) error
	pathmap map[string]string
}

// TODO: support context
func (t *transformer) Transform(fromPath, toPath string, code []byte) ([]byte, error) {
	fromExt := filepath.Ext(fromPath)
	hops, err := t.graph.ShortestPath(fromExt, filepath.Ext(toPath))
	if err != nil {
		return nil, err
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
		Code: code,
	}
	// Apply transformations over the transform pairs
	for _, pair := range pairs {
		// Handle .svelte -> .svelte transformations
		key := pair[0] + ">" + pair[1]
		if transform, ok := t.index[key]; ok {
			if err := transform(file); err != nil {
				return nil, err
			}
			// Update the extension
			file.ext = pair[1]
		}
	}
	return file.Code, nil
}

func (t *transformer) Plugins() (plugins []esbuild.Plugin) {
	for from, to := range t.pathmap {
		from := from
		plugins = append(plugins, esbuild.Plugin{
			Name: strings.TrimPrefix(from, ".") + "_to_" + strings.TrimPrefix(to, "."),
			Setup: func(epb esbuild.PluginBuild) {
				// Load svelte files. Add import if not present
				epb.OnLoad(esbuild.OnLoadOptions{Filter: `\` + from + `$`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
					// Read the code in
					code, err := os.ReadFile(args.Path)
					if err != nil {
						return result, err
					}
					fromPath := args.Path
					toPath := strings.TrimSuffix(args.Path, from) + "." + to
					// Transform the code
					// TODO: We wouldn't need to get the shortest path in Transform
					// everytime, we could pre-compute these shortest paths.
					newCode, err := t.Transform(fromPath, toPath, code)
					if err != nil {
						return result, err
					}
					// Update the file contents
					contents := string(newCode)
					result.ResolveDir = filepath.Dir(args.Path)
					result.Contents = &contents
					// Use an appropriate loader that esbuild understands
					switch to {
					case ".js":
						result.Loader = esbuild.LoaderJS
					case ".jsx":
						result.Loader = esbuild.LoaderJSX
					default:
						return result, fmt.Errorf("transform: unhandled loader type %q", to)
					}
					return result, nil
				})
			},
		})
	}
	return plugins
}
