package svelte

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/livebud/bud/internal/imports"
	"github.com/matthewmueller/gotext"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/internal/es"
	"github.com/livebud/bud/package/viewer"
	"github.com/livebud/bud/runtime/transpiler"
	"github.com/livebud/js"
)

func New(esbuilder *es.Builder, fsys fs.FS, js js.VM, transpiler transpiler.Interface, pages viewer.Pages) *Viewer {
	return &Viewer{esbuilder, fsys, js, pages, transpiler}
}

type Viewer struct {
	esbuilder  *es.Builder
	fsys       fs.FS
	js         js.VM
	pages      viewer.Pages
	transpiler transpiler.Interface
}

var _ viewer.Viewer = (*Viewer)(nil)

func (v *Viewer) Render(ctx context.Context, key string, props viewer.Props) ([]byte, error) {
	page, ok := v.pages[key]
	if !ok {
		return nil, fmt.Errorf("svelte: %q. %w", key, viewer.ErrPageNotFound)
	}
	ssrCode, err := v.esbuilder.Build(&es.Build{
		Entrypoint: page.Path + ".js",
		Plugins: []es.Plugin{
			v.ssrEntryPlugin(page),
			v.ssrRuntimePlugin(),
			v.ssrTranspile(),
		},
	})
	if err != nil {
		return nil, err
	}
	propBytes, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}
	expr := fmt.Sprintf(`%s; bud.render(%q, %s)`, ssrCode, key, propBytes)
	html, err := v.js.Evaluate(ctx, page.Path, expr)
	if err != nil {
		return nil, err
	}
	return []byte(html), nil
}

func (v *Viewer) RenderError(ctx context.Context, key string, err error, props viewer.Props) []byte {
	return []byte("RenderError not implemented yet")
}

func (v *Viewer) Register(router viewer.Router) {
	fmt.Println("register called")
}

func (v *Viewer) Bundle(ctx context.Context, out viewer.FS) error {
	for _, page := range v.pages {
		fmt.Println("bundling", page.Path)
	}
	return nil
}

//go:embed ssr_entry.gotext
var ssrEntryCode string

var ssrEntryTemplate = template.Must(template.New("ssr_entry.gotext").Parse(ssrEntryCode))

type State struct {
	// Note: we're slightly abusing imports.Import here, since those are meant
	// for Go imports, not JS imports. But it works out for this use case.
	Imports []*imports.Import
	Pages   []*Page
}

type Page struct {
	*View
	Layout *View
	Error  *View
	Frames []*View
}

type View struct {
	Path      string
	Key       string
	Component string
}

func newState(pages ...*viewer.Page) *State {
	state := new(State)
	imports := imports.New()
	for _, p := range pages {
		page := new(Page)
		page.View = &View{
			Path:      p.Path,
			Key:       p.Key,
			Component: imports.AddNamed(gotext.Pascal(p.Key), p.Path),
		}
		if p.Error != nil {
			page.Error = &View{
				Path:      p.Error.Path,
				Key:       p.Error.Key,
				Component: imports.AddNamed(gotext.Pascal(p.Error.Key), p.Error.Path),
			}
		}
		if p.Layout != nil {
			page.Layout = &View{
				Path:      p.Layout.Path,
				Key:       p.Layout.Key,
				Component: imports.AddNamed(gotext.Pascal(p.Layout.Key), p.Layout.Path),
			}
		}
		for _, frame := range p.Frames {
			page.Frames = append(page.Frames, &View{
				Path:      frame.Path,
				Key:       frame.Key,
				Component: imports.AddNamed(gotext.Pascal(frame.Key), frame.Path),
			})
		}
		state.Pages = append(state.Pages, page)
	}
	state.Imports = imports.List()
	return state
}

func (v *Viewer) ssrEntryPlugin(page *viewer.Page) es.Plugin {
	return es.Plugin{
		Name: "svelte_ssr_entry",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^` + page.Path + `\.js$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = page.Path + `.js`
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: page.Path + `.js`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				code := new(bytes.Buffer)
				state := newState(page)
				if err := ssrEntryTemplate.Execute(code, state); err != nil {
					return result, err
				}
				if err != nil {
					return result, err
				}
				contents := code.String()
				result.ResolveDir = v.esbuilder.Directory()
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}

//go:embed ssr_runtime.ts
var ssrRuntimeCode string

func (v *Viewer) ssrRuntimePlugin() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "svelte_ssr_runtime",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^ssrSvelteRuntime$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "ssrSvelteRuntime"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `ssrSvelteRuntime`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				result.ResolveDir = v.esbuilder.Directory()
				result.Contents = &ssrRuntimeCode
				result.Loader = esbuild.LoaderTS
				return result, nil
			})
		},
	}
}

// Svelte plugin transforms Svelte imports to server-side JS
func (v *Viewer) ssrTranspile() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "ssr_transpile",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `\.svelte$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "svelte"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `svelte`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				filePath := filepath.Clean(args.Path)
				svelteCode, err := fs.ReadFile(v.fsys, filePath)
				if err != nil {
					return result, err
				}
				ssrJSCode, err := v.transpiler.Transpile(filePath, ".ssr.js", svelteCode)
				if err != nil {
					return result, err
				}
				contents := string(ssrJSCode)
				result.ResolveDir = v.esbuilder.Directory()
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}
