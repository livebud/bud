package svelte

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/internal/gotemplate"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/js"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/router"
	"github.com/livebud/bud/package/transpiler"
	"github.com/livebud/bud/package/viewer"
)

type FS = fs.FS

func New(fsys FS, log log.Log, module *gomod.Module, transpiler transpiler.Interface, vm js.VM) *Viewer {
	return &Viewer{fsys, log, module, transpiler, vm}
}

type Viewer struct {
	fsys       FS
	log        log.Log
	module     *gomod.Module
	transpiler transpiler.Interface
	vm         js.VM
}

var _ viewer.Viewer = (*Viewer)(nil)

func (v *Viewer) Render(ctx context.Context, page *viewer.Page) ([]byte, error) {
	code, err := v.compileSSR(page)
	if err != nil {
		return nil, err
	}
	props, err := json.Marshal(page)
	if err != nil {
		return nil, err
	}
	expr := fmt.Sprintf("%s\nbud.render(%s)", code, props)
	html, err := v.vm.Eval(page.Path, expr)
	if err != nil {
		return nil, err
	}
	return []byte(html), nil
}

func (v *Viewer) RenderError(ctx context.Context, page *viewer.Page) []byte {
	code, err := v.Render(ctx, page)
	if err != nil {
		return []byte(err.Error())
	}
	return code
}

func (v *Viewer) compileSSR(page *viewer.Page) (code []byte, err error) {
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPointsAdvanced: []esbuild.EntryPoint{
			{
				InputPath:  page.Path + ".js",
				OutputPath: page.Path + ".js",
			},
		},
		AbsWorkingDir: v.module.Directory(),
		Outdir:        "./",
		Format:        esbuild.FormatIIFE,
		Platform:      esbuild.PlatformBrowser,
		GlobalName:    "bud",
		Bundle:        true,
		// Metafile:      true,
		Plugins: []esbuild.Plugin{
			v.ssrSveltePage(page),
			v.ssrSvelteRuntime(),
			v.ssrSvelteTransform(),
			v.httpPlugin(),
			v.esmPlugin(),
		},
	})
	if len(result.Errors) > 0 {
		msgs := esbuild.FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
			Color:         true,
			Kind:          esbuild.ErrorMessage,
			TerminalWidth: 80,
		})
		return nil, fmt.Errorf(strings.Join(msgs, "\n"))
	}
	// Expect exactly 1 output file
	if len(result.OutputFiles) != 1 {
		return nil, fmt.Errorf("expected exactly 1 output file but got %d", len(result.OutputFiles))
	}
	ssrCode := result.OutputFiles[0].Contents
	return ssrCode, nil
}

//go:embed ssr.gotext
var ssrCode string

var ssrTemplate = gotemplate.MustParse("ssr.gotext", ssrCode)

// page plugin creates a virtual page point for the page
func (v *Viewer) ssrSveltePage(page *viewer.Page) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "ssrSveltePage",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^` + page.Path + `\.js$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = page.Path + `.js`
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: page.Path + `.js`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				code, err := ssrTemplate.Generate(page)
				if err != nil {
					return result, err
				}
				contents := string(code)
				result.ResolveDir = v.module.Directory()
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}

//go:embed ssr.ts
var ssrRuntime string

func (v *Viewer) ssrSvelteRuntime() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "ssrSvelteRuntime",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^ssrSvelteRuntime$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "ssrSvelteRuntime"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `ssrSvelteRuntime`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				result.ResolveDir = v.module.Directory()
				result.Contents = &ssrRuntime
				result.Loader = esbuild.LoaderTS
				return result, nil
			})
		},
	}
}

// Svelte plugin transforms Svelte imports to server-side JS
func (v *Viewer) ssrSvelteTransform() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "ssrSvelteTransform",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `\.svelte$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "svelte"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `svelte`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				code, err := fs.ReadFile(v.fsys, filepath.Clean(args.Path))
				if err != nil {
					return result, err
				}
				transpiled, err := v.transpiler.Transpile(args.Path, ".ssr.js", code)
				if err != nil {
					return result, err
				}
				contents := string(transpiled)
				result.ResolveDir = v.module.Directory()
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}

func (v *Viewer) RegisterClient(r *router.Router, page *viewer.Page) {
	r.Get(page.Client, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code, err := v.compileDOM(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(code)
	}))
}

func (v *Viewer) compileDOM(page *viewer.Page) (code []byte, err error) {
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPointsAdvanced: []esbuild.EntryPoint{
			{
				InputPath:  page.Client,
				OutputPath: page.Client,
			},
		},
		AbsWorkingDir:     v.module.Directory(),
		Outdir:            "./",
		Format:            esbuild.FormatIIFE,
		Platform:          esbuild.PlatformBrowser,
		GlobalName:        "bud",
		MinifyIdentifiers: false,
		MinifySyntax:      false,
		MinifyWhitespace:  false,
		Bundle:            true,
		// Metafile:      true,
		Plugins: []esbuild.Plugin{
			v.domSveltePage(page),
			v.domSvelteRuntime(),
			v.domSvelteTransform(),
			v.httpPlugin(),
			v.esmPlugin(),
		},
	})
	if len(result.Errors) > 0 {
		msgs := esbuild.FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
			Color:         true,
			Kind:          esbuild.ErrorMessage,
			TerminalWidth: 80,
		})
		return nil, fmt.Errorf(strings.Join(msgs, "\n"))
	}
	// Expect exactly 1 output file
	if len(result.OutputFiles) != 1 {
		return nil, fmt.Errorf("expected exactly 1 output file but got %d", len(result.OutputFiles))
	}
	ssrCode := result.OutputFiles[0].Contents
	return ssrCode, nil
}

//go:embed dom.gotext
var domCode string

var domTemplate = gotemplate.MustParse("dom.gotext", domCode)

// page plugin creates a virtual page point for the page
func (v *Viewer) domSveltePage(page *viewer.Page) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "domSveltePage",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^` + page.Client + `$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "domSveltePage"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `^` + page.Client + `$`, Namespace: "domSveltePage"}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				code, err := domTemplate.Generate(page)
				if err != nil {
					return result, err
				}
				contents := string(code)
				result.ResolveDir = v.module.Directory()
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}

//go:embed dom.ts
var domRuntime string

func (v *Viewer) domSvelteRuntime() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "domSvelteRuntime",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^domSvelteRuntime$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "domSvelteRuntime"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `domSvelteRuntime`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				result.ResolveDir = v.module.Directory()
				result.Contents = &domRuntime
				result.Loader = esbuild.LoaderTS
				return result, nil
			})
		},
	}
}

// Svelte plugin transforms Svelte imports to server-side JS
func (v *Viewer) domSvelteTransform() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "domSvelteTransform",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `\.svelte$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "svelte"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `svelte`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				code, err := fs.ReadFile(v.fsys, filepath.Clean(args.Path))
				if err != nil {
					return result, err
				}
				transpiled, err := v.transpiler.Transpile(args.Path, ".js", code)
				if err != nil {
					return result, err
				}
				contents := string(transpiled)
				result.ResolveDir = v.module.Directory()
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}

func (v *Viewer) httpPlugin() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "bud-http",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^http[s]?://`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "http"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `http`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				res, err := http.Get(args.Path)
				if err != nil {
					return result, err
				}
				defer res.Body.Close()
				body, err := io.ReadAll(res.Body)
				if err != nil {
					return result, err
				}
				contents := string(body)
				result.ResolveDir = v.module.Directory()
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}

func (v *Viewer) esmPlugin() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "bud-esm",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^[a-z0-9@]`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "esm"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `esm`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				res, err := http.Get("https://esm.sh/" + args.Path)
				if err != nil {
					return result, err
				}
				defer res.Body.Close()
				body, err := io.ReadAll(res.Body)
				if err != nil {
					return result, err
				}
				contents := string(body)
				result.ResolveDir = v.module.Directory()
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}
