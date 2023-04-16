package svelte

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"text/template"

	esbuild "github.com/evanw/esbuild/pkg/api"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/es"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/router"
	"github.com/livebud/bud/package/transpiler"
	"github.com/livebud/bud/package/viewer"
	"github.com/livebud/bud/package/virtual"
	"github.com/livebud/js"
	"github.com/matthewmueller/gotext"
)

func New(es es.Builder, flag *framework.Flag, js js.VM, log log.Log, module *gomod.Module, pages viewer.Pages, transpiler transpiler.Interface) *Viewer {
	return &Viewer{es, flag, js, log, module, pages, transpiler}
}

type Viewer struct {
	es         es.Builder
	flag       *framework.Flag
	js         js.VM
	log        log.Log
	module     *gomod.Module
	pages      viewer.Pages
	transpiler transpiler.Interface
}

var _ viewer.Viewer = (*Viewer)(nil)

func (v *Viewer) Mount(r *router.Router) error {
	for _, page := range v.pages {
		// Serve the entrypoints (for hydrating)
		r.Get(page.Client.Route, v.serveDOMEntry(page))
		// Serve the individual views themselves (for hot reloads)
		r.Get(page.View.Client.Route, v.serveDOMView(page.View))
	}
	return nil
}

func (v *Viewer) Render(ctx context.Context, key string, propMap viewer.PropMap) ([]byte, error) {
	page, ok := v.pages[key]
	if !ok {
		return nil, fmt.Errorf("svelte: unable to find page from key %q", key)
	}
	v.log.Info("svelte: rendering", page.Path)
	file, err := v.compileSSR(ctx, page)
	if err != nil {
		return nil, err
	}
	propBytes, err := json.Marshal(propMap)
	if err != nil {
		return nil, err
	}
	expr := fmt.Sprintf(`%s; bud.render(%s)`, string(file.Contents), propBytes)
	html, err := v.js.Evaluate(ctx, page.Path, expr)
	if err != nil {
		return nil, err
	}
	return []byte(html), nil
}

func (v *Viewer) RenderError(ctx context.Context, key string, propMap viewer.PropMap, originalError error) []byte {
	page, ok := v.pages[key]
	if !ok {
		return []byte(fmt.Sprintf("svelte: unable to find page from key %q to render error. %s", key, originalError))
	}
	if page.Error == nil {
		return []byte(fmt.Sprintf("svelte: no error page for %q to render error. %s", key, originalError))
	}
	errorPage, ok := v.pages[page.Error.Key]
	if !ok {
		return []byte(fmt.Sprintf("svelte: unable to find error page for %q to render error. %s", page.Error.Key, originalError))
	}
	v.log.Info("svelte: rendering error", errorPage.Path)
	file, err := v.compileSSR(ctx, errorPage)
	if err != nil {
		return []byte(fmt.Sprintf("svelte: unable to serve error page %q to render error. %s. %s", errorPage.Path, err, originalError))
	}
	propMap[errorPage.Key] = viewer.Error(originalError)
	propBytes, err := json.Marshal(propMap)
	if err != nil {
		return []byte(fmt.Sprintf("svelte: unable to marshal props for %q to render error. %s. %s", errorPage.Path, err, originalError))
	}
	expr := fmt.Sprintf(`%s; bud.render(%s)`, string(file.Contents), propBytes)
	html, err := v.js.Evaluate(ctx, errorPage.Path, expr)
	if err != nil {
		return []byte(fmt.Sprintf("svelte: unable to evaluate javascript to render %q to render error. %s. %s", errorPage.Path, err, originalError))
	}
	return []byte(html)
}

func (v *Viewer) compileSSR(ctx context.Context, page *viewer.Page) (*es.File, error) {
	return v.es.Serve(&es.Serve{
		AbsDir:   v.module.Directory(),
		Entry:    "./" + page.Path + ".js",
		Platform: es.SSR,
		Plugins: []es.Plugin{
			v.ssrEntryPlugin(page),
			v.ssrRuntimePlugin(),
			v.ssrTranspile(ctx),
			es.HTTP(http.DefaultClient),
			es.ImportMap(v.log, map[string]string{
				"svelte":  "https://esm.run/svelte@" + versions.Svelte,
				"svelte/": "https://esm.run/svelte@" + versions.Svelte + "/",
			}),
		},
	})
}

func (v *Viewer) Bundle(ctx context.Context, embed virtual.Tree) error {
	for _, page := range v.pages {
		file, err := v.compileSSR(ctx, page)
		if err != nil {
			return err
		}
		embed[page.Path] = &virtual.File{
			Path: page.Path,
			Mode: 0644,
			Data: file.Contents,
		}
		file, err = v.compileDOMEntry(ctx, page)
		if err != nil {
			return err
		}
		embed[page.Client.Path] = &virtual.File{
			Path: page.Client.Path,
			Mode: 0644,
			Data: file.Contents,
		}
		file, err = v.compileDOMView(ctx, page.View)
		if err != nil {
			return err
		}
		embed[page.View.Client.Path] = &virtual.File{
			Path: page.View.Client.Path,
			Mode: 0644,
			Data: file.Contents,
		}
	}
	return nil
}

// Handler serves the page as a static view
func (v *Viewer) Handler(page *viewer.Page) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		propMap, err := viewer.StaticPropMap(page, r)
		if err != nil {
			v.renderError(ctx, w, page, propMap, err)
			return
		}
		html, err := v.Render(ctx, page.Key, propMap)
		if err != nil {
			v.renderError(ctx, w, page, propMap, err)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(html)
	})
}

func (v *Viewer) renderError(ctx context.Context, w http.ResponseWriter, page *viewer.Page, propMap map[string]interface{}, err error) {
	html := v.RenderError(ctx, page.Key, propMap, err)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(html)
}

//go:embed ssr_entry.gotext
var ssrEntryCode string

var ssrEntryTemplate = template.Must(template.New("ssr_entry.gotext").Parse(ssrEntryCode))

func (v *Viewer) ssrEntryPlugin(page *viewer.Page) es.Plugin {
	return es.Plugin{
		Name: "svelte_ssr_entry",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^./` + page.Path + `\.js$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = `svelte_ssr_entry`
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `svelte_ssr_entry`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				type View struct {
					Path      string
					Key       string
					Component string
					Client    *viewer.Client
				}
				type Page struct {
					*View
					Layout *View
					Error  *View
					Frames []*View
				}
				type State struct {
					// Note: we're slightly abusing imports.Import here, since those are meant
					// for Go imports, not JS imports. But it works out for this use case.
					Imports []*imports.Import
					Page    *Page
				}
				// Load the SSR state
				state := new(State)
				imports := imports.New()
				state.Page = &Page{
					View: &View{
						Path:      page.Path,
						Key:       page.Key,
						Component: imports.AddNamed(gotext.Pascal(page.Key), page.Path),
						Client:    page.Client,
					},
				}
				if page.Error != nil {
					state.Page.Error = &View{
						Path:      page.Error.Path,
						Key:       page.Error.Key,
						Component: imports.AddNamed(gotext.Pascal(page.Error.Key), page.Error.Path),
					}
				}
				if page.Layout != nil {
					state.Page.Layout = &View{
						Path:      page.Layout.Path,
						Key:       page.Layout.Key,
						Component: imports.AddNamed(gotext.Pascal(page.Layout.Key), page.Layout.Path),
					}
				}
				for _, frame := range page.Frames {
					state.Page.Frames = append(state.Page.Frames, &View{
						Path:      frame.Path,
						Key:       frame.Key,
						Component: imports.AddNamed(gotext.Pascal(frame.Key), frame.Path),
					})
				}
				state.Imports = imports.List()
				// Generate the SSR entry code
				code := new(bytes.Buffer)
				if err := ssrEntryTemplate.Execute(code, state); err != nil {
					return result, err
				}
				if err != nil {
					return result, err
				}
				contents := code.String()
				result.ResolveDir = v.module.Directory()
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
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^\.svelte_ssr_runtime$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "svelte_ssr_runtime"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `svelte_ssr_runtime`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				result.Contents = &ssrRuntimeCode
				result.ResolveDir = v.module.Directory()
				result.Loader = esbuild.LoaderTS
				return result, nil
			})
		},
	}
}

// Svelte plugin transforms Svelte imports to server-side JS
func (v *Viewer) ssrTranspile(ctx context.Context) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "ssr_transpile",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `\.svelte$`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				relPath, err := filepath.Rel(v.module.Directory(), args.Path)
				if err != nil {
					return result, err
				}
				code, err := fs.ReadFile(v.module, relPath)
				if err != nil {
					return result, err
				}
				ssrJsCode, err := v.transpiler.Transpile(ctx, args.Path, ".ssr.js", code)
				if err != nil {
					return result, err
				}
				contents := string(ssrJsCode)
				result.ResolveDir = v.module.Directory()
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}

// serveDOM serves the entrypoints (for hydrating)
func (v *Viewer) serveDOMEntry(page *viewer.Page) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.log.Info("svelte: serving client entry", r.URL.Path)
		domJSCode, err := v.compileDOMEntry(r.Context(), page)
		if err != nil {
			// TODO: hydrate a nice error message in the client
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.WriteHeader(http.StatusOK)
		w.Write(domJSCode.Contents)
	})
}

// Compile DOM entrypoint
func (v *Viewer) compileDOMEntry(ctx context.Context, page *viewer.Page) (*es.File, error) {
	return v.es.Serve(&es.Serve{
		AbsDir:   v.module.Directory(),
		Entry:    "./" + page.Path + ".js",
		Platform: es.DOM,
		Plugins: []es.Plugin{
			v.domEntryPlugin(page),
			v.domRuntimePlugin(),
			v.domExternals(),
			v.domTranspile(ctx),
			es.ExternalHTTP(),
			es.ExternalImportMap(v.log, map[string]string{
				"svelte":  "https://esm.run/svelte@" + versions.Svelte,
				"svelte/": "https://esm.run/svelte@" + versions.Svelte + "/",
			}),
		},
	})
}

// serveDOM serves the individual views themselves (for hot reloads)
func (v *Viewer) serveDOMView(view *viewer.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.log.Info("svelte: serving client view", r.URL.Path)
		domJsCode, err := v.compileDOMView(r.Context(), view)
		if err != nil {
			// TODO: hydrate a nice error message in the client
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.WriteHeader(http.StatusOK)
		w.Write(domJsCode.Contents)
	})
}

// Compile DOM view
func (v *Viewer) compileDOMView(ctx context.Context, view *viewer.View) (*es.File, error) {
	return v.es.Serve(&es.Serve{
		AbsDir:   v.module.Directory(),
		Entry:    "./" + view.Path,
		Platform: es.DOM,
		Plugins: []es.Plugin{
			v.domTranspile(ctx),
			es.ExternalHTTP(),
			es.ExternalImportMap(v.log, map[string]string{
				"svelte":  "https://esm.run/svelte@" + versions.Svelte,
				"svelte/": "https://esm.run/svelte@" + versions.Svelte + "/",
			}),
		},
	})
}

//go:embed dom_entry.gotext
var domEntryCode string

var domEntryTemplate = template.Must(template.New("dom_entry.gotext").Parse(domEntryCode))

func (v *Viewer) domEntryPlugin(page *viewer.Page) es.Plugin {
	return es.Plugin{
		Name: "svelte_dom_entry",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^./` + page.Path + `.js$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = page.Path + `.js`
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: page.Path + `.js`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				type View struct {
					Path      string
					Key       string
					Component string
				}
				type Page struct {
					*View
					Error  *View
					Frames []*View
				}
				type State struct {
					// Note: we're slightly abusing imports.Import here, since those are meant
					// for Go imports, not JS imports. But it works out for this use case.
					Imports []*imports.Import
					Page    *Page
					Hot     string
				}
				state := new(State)
				imports := imports.New()
				state.Page = &Page{
					View: &View{
						Path:      page.Path,
						Key:       page.Key,
						Component: imports.AddNamed(gotext.Pascal(page.Key), page.View.Client.Route),
					},
				}
				if page.Error != nil {
					state.Page.Error = &View{
						Path:      page.Error.Path,
						Key:       page.Error.Key,
						Component: imports.AddNamed(gotext.Pascal(page.Error.Key), page.Error.Client.Route),
					}
				}
				for _, frame := range page.Frames {
					state.Page.Frames = append(state.Page.Frames, &View{
						Path:      frame.Path,
						Key:       frame.Key,
						Component: imports.AddNamed(gotext.Pascal(frame.Key), frame.Client.Route),
					})
				}
				state.Imports = imports.List()
				if v.flag.Hot {
					// TODO: configurable
					state.Hot = `http://127.0.0.1:35729/bud/hot?path=` + page.Path
				}
				code := new(bytes.Buffer)
				if err := domEntryTemplate.Execute(code, state); err != nil {
					return result, err
				}
				if err != nil {
					return result, err
				}
				contents := code.String()
				result.ResolveDir = v.module.Directory()
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}

//go:embed dom_runtime.ts
var domRuntimeCode string

func (v *Viewer) domRuntimePlugin() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "svelte_dom_runtime",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^\.svelte_dom_runtime$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "svelte_dom_runtime"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `svelte_dom_runtime`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				result.ResolveDir = v.module.Directory()
				result.Contents = &domRuntimeCode
				result.Loader = esbuild.LoaderTS
				return result, nil
			})
		},
	}
}

func (v *Viewer) domExternals() es.Plugin {
	return es.Plugin{
		Name: "svelte_externals",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^/view/.*\.svelte\.js$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Path = args.Path
				result.External = true
				return result, nil
			})
		},
	}
}

// Svelte plugin transforms Svelte imports to client-side JS
func (v *Viewer) domTranspile(ctx context.Context) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "dom_transpile",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `\.svelte$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Path = args.Path
				result.Namespace = "dom_transpile"
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `\.svelte$`, Namespace: "dom_transpile"}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				code, err := fs.ReadFile(v.module, path.Clean(args.Path))
				if err != nil {
					return result, err
				}
				domJsCode, err := v.transpiler.Transpile(ctx, path.Clean(args.Path), ".dom.js", code)
				if err != nil {
					return result, err
				}
				contents := string(domJsCode)
				result.ResolveDir = v.module.Directory()
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}
