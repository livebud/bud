package svelte

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/livebud/bud/internal/gotemplate"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/framework/transform2/transformrt"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/js"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/router"
	"github.com/livebud/bud/package/viewer"
)

func New(fsys fs.FS, log log.Log, module *gomod.Module, transformer *transformrt.Transformer, vm js.VM) *Viewer {
	return &Viewer{fsys, log, module, transformer, vm}
}

type Viewer struct {
	fsys        fs.FS
	log         log.Log
	module      *gomod.Module
	transformer *transformrt.Transformer
	vm          js.VM
}

var _ viewer.Viewer = (*Viewer)(nil)

func (v *Viewer) Render(ctx context.Context, page *viewer.Page) ([]byte, error) {
	code, err := v.compile(page)
	if err != nil {
		return nil, err
	}
	props, err := json.Marshal(page)
	if err != nil {
		return nil, err
	}
	expr := fmt.Sprintf("%s\nbud.render(%s)", code, props)
	html, err := v.vm.Eval(page.Main.Path, expr)
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

func (v *Viewer) RegisterPage(r *router.Router, page *viewer.Page) {
}

func (v *Viewer) compile(page *viewer.Page) (code []byte, err error) {
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPointsAdvanced: []esbuild.EntryPoint{
			{
				InputPath:  page.Main.Path + ".js",
				OutputPath: path.Join(page.Main.Path + ".js"),
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
			v.pagePlugin(page),
			v.svelteRuntime(),
			v.sveltePlugin(),
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

//go:embed page.gotext
var pageCode string

var pageTemplate = gotemplate.MustParse("svelte/page.gotext", pageCode)

// page plugin creates a virtual page point for the page
func (v *Viewer) pagePlugin(page *viewer.Page) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "bud-svelte-page",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^` + page.Main.Path + `\.js$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = page.Main.Path + `.js`
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: page.Main.Path + `.js`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				code, err := pageTemplate.Generate(page)
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

//go:embed runtime.ts
var runtimeCode string

func (v *Viewer) svelteRuntime() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "bud-svelte-runtime",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^bud-svelte-runtime$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "svelte-runtime"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `svelte-runtime`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				result.ResolveDir = v.module.Directory()
				result.Contents = &runtimeCode
				result.Loader = esbuild.LoaderTS
				return result, nil
			})
		},
	}
}

// Svelte plugin transforms Svelte imports to server-side JS
func (v *Viewer) sveltePlugin() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "bud-svelte",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `\.svelte$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "svelte"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `svelte`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				code, err := v.transformer.Transform(v.fsys, path.Join("ssr.js", args.Path))
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
