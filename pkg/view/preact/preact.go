package preact

//go:generate go run github.com/evanw/esbuild/cmd/esbuild@v0.19.5 --bundle ssr.tsx

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/internal/esb"
	"github.com/livebud/bud/internal/js"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/view"
)

func Embed(embed bool) func(*Viewer) {
	return func(v *Viewer) {
		v.embed = embed
	}
}

func Minify(minify bool) func(*Viewer) {
	return func(v *Viewer) {
		v.minify = minify
	}
}

// func WithJS(vm )

func New(module *mod.Module, options ...func(*Viewer)) *Viewer {
	p := &Viewer{
		module: module,
		embed:  false,
		minify: false,
	}
	for _, option := range options {
		option(p)
	}
	return p
}

var _ view.Viewer = (*Viewer)(nil)

type Viewer struct {
	module *mod.Module
	embed  bool
	minify bool
}

func (v *Viewer) Render(w io.Writer, path string, data *view.Data) error {
	return v.renderHTML(w, path, data)
}

func (v *Viewer) Routes(router mux.Router) {
	router.Get("/.preact/{path*}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		viewPath := strings.TrimSuffix(r.URL.Query().Get("path"), ".js")
		if viewPath == "" {
			http.Error(w, fmt.Sprintf("preact: client view %q not found", r.URL.Path), http.StatusNotFound)
			return
		}
		entry, err := v.CompileDOM("./" + viewPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("preact: unable to compile %q. %s", viewPath, err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(entry.Contents)
	}))
}

func (v *Viewer) renderHTML(w io.Writer, path string, data *view.Data) error {
	entry, err := v.CompileSSR("./" + path)
	if err != nil {
		return err
	}
	// TODO: optimize
	program, err := goja.Compile(path, string(entry.Contents), false)
	if err != nil {
		return err
	}
	vm := js.New()
	_, err = vm.RunProgram(program)
	if err != nil {
		return err
	}
	props, err := json.Marshal(data.Props)
	if err != nil {
		return err
	}
	result, err := vm.RunString(fmt.Sprintf("bud.render(%s)", props))
	if err != nil {
		return err
	}
	var ssr SSR
	err = json.Unmarshal([]byte(result.String()), &ssr)
	if err != nil {
		return err
	}
	if ssr.Head != "" {
		data.Slots.Slot("head").Write([]byte(ssr.Head))
	}
	data.Slots.Slot("head").Write([]byte(fmt.Sprintf(`<script id=".bud_props" type="text/template" defer>%s</script>`, props)))
	data.Slots.Slot("head").Write([]byte(fmt.Sprintf(`<script src="/.preact/%s.js" defer></script>`, path)))
	w.Write([]byte(ssr.HTML))
	return nil
}

type SSR struct {
	HTML string `json:"html"`
	Head string `json:"head"`
}

func (v *Viewer) Middleware(next http.Handler) http.Handler {
	return view.Middleware(v).Middleware(next)
}

//go:embed ssr_runtime.tsx
var ssrRuntime string

func (v *Viewer) ssrRuntime(entry string) func(esbuild.OnLoadArgs) (esbuild.OnLoadResult, error) {
	return func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
		result.Contents = &ssrRuntime
		result.ResolveDir = filepath.Dir(args.Path)
		result.Loader = esbuild.LoaderTSX
		return result, nil
	}
}

//go:embed ssr_entry.gotext
var ssrEntry string

func (v *Viewer) ssrEntry(entry string) func(esbuild.OnLoadArgs) (esbuild.OnLoadResult, error) {
	return func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
		result.Contents = &ssrEntry
		result.ResolveDir = filepath.Dir(args.Path)
		result.Loader = esbuild.LoaderJS
		return result, nil
	}
}

func spliceExt(path, ext string) string {
	end := filepath.Ext(path)
	return strings.TrimSuffix(path, end) + ext + end
}

func (v *Viewer) CompileSSR(entry string) (*esb.File, error) {
	if !isRelativeEntry(entry) {
		return nil, fmt.Errorf("entry must be relative %q", entry)
	}
	absDir, err := filepath.Abs(v.module.Directory())
	if err != nil {
		return nil, err
	}
	ssrEntry := spliceExt(entry, ".ssr")
	options := esbuild.BuildOptions{
		AbsWorkingDir:   absDir,
		EntryPoints:     []string{ssrEntry},
		Format:          esbuild.FormatIIFE,
		Platform:        esbuild.PlatformNeutral,
		JSXImportSource: "preact",
		JSX:             esbuild.JSXAutomatic,
		Bundle:          true,
		GlobalName:      "bud",
		Plugins: []esbuild.Plugin{
			esb.Virtual("ssr_runtime.tsx", v.ssrRuntime(entry)),
			esb.Virtual(ssrEntry, v.ssrEntry(entry)),
			esb.HTTP(http.DefaultClient),
		},
	}
	result := esbuild.Build(options)
	if result.Errors != nil {
		return nil, &esb.Error{Messages: result.Errors}
	} else if len(result.OutputFiles) == 0 {
		return nil, fmt.Errorf("esb: no output files")
	}
	return &result.OutputFiles[0], nil
}

func isRelativeEntry(entry string) bool {
	return strings.HasPrefix(entry, "./")
}

//go:embed dom_runtime.tsx
var domRuntime string

func (v *Viewer) domRuntime(entry string) func(esbuild.OnLoadArgs) (esbuild.OnLoadResult, error) {
	return func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
		result.Contents = &domRuntime
		result.ResolveDir = filepath.Dir(args.Path)
		result.Loader = esbuild.LoaderTSX
		return result, nil
	}
}

//go:embed dom_entry.gotext
var domEntry string

func (v *Viewer) domEntry(entry string) func(esbuild.OnLoadArgs) (esbuild.OnLoadResult, error) {
	return func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
		result.Contents = &domEntry
		result.ResolveDir = filepath.Dir(args.Path)
		result.Loader = esbuild.LoaderJS
		return result, nil
	}
}

func (v *Viewer) CompileDOM(entry string) (*esb.File, error) {
	if !isRelativeEntry(entry) {
		return nil, fmt.Errorf("entry must be relative %q", entry)
	}
	absDir, err := filepath.Abs(v.module.Directory())
	if err != nil {
		return nil, err
	}
	domEntry := spliceExt(entry, ".dom")
	options := esbuild.BuildOptions{
		AbsWorkingDir:   absDir,
		EntryPoints:     []string{domEntry},
		Format:          esbuild.FormatIIFE,
		Platform:        esbuild.PlatformBrowser,
		JSXImportSource: "preact",
		JSX:             esbuild.JSXAutomatic,
		Bundle:          true,
		GlobalName:      "bud",
		Plugins: []esbuild.Plugin{
			esb.Virtual("dom_runtime.tsx", v.domRuntime(entry)),
			esb.Virtual(domEntry, v.domEntry(entry)),
			esb.HTTP(http.DefaultClient),
		},
	}
	result := esbuild.Build(options)
	if result.Errors != nil {
		return nil, &esb.Error{Messages: result.Errors}
	} else if len(result.OutputFiles) == 0 {
		return nil, fmt.Errorf("esb: no output files")
	}
	return &result.OutputFiles[0], nil
}

// var jsxPlugin = esbuild.Plugin{
// 	Name: "prepend-create-element",
// 	Setup: func(epb esbuild.PluginBuild) {
// 		epb.OnLoad(esbuild.OnLoadOptions{Filter: `\.[jt]sx`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
// 			contents := "import { createElement as h } from \"preact\"\n"
// 			if result.Contents != nil {
// 				contents += *result.Contents
// 			}
// 			code, err := esbuild.LoaderDefault
// 			fmt.Println(args.Path)
// 			fmt.Println(contents)
// 			result.Contents = &contents
// 			return result, nil
// 		})
// 	},
// }
