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
	"text/template"

	"github.com/dop251/goja"
	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/internal/esb"
	"github.com/livebud/bud/internal/js"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/view"
)

func WithEmbed(embed bool) func(*Viewer) {
	return func(v *Viewer) {
		v.embed = embed
	}
}

func WithMinify(minify bool) func(*Viewer) {
	return func(v *Viewer) {
		v.minify = minify
	}
}

func WithLive(liveUrl string) func(*Viewer) {
	return func(v *Viewer) {
		v.liveUrl = liveUrl
	}
}

func WithEnv(obj any) func(*Viewer) {
	return func(v *Viewer) {
		v.env = obj
	}
}

func New(module *mod.Module, options ...func(*Viewer)) *Viewer {
	p := &Viewer{
		module:  module,
		embed:   false,
		minify:  false,
		liveUrl: "",
		env:     map[string]any{},
	}
	for _, option := range options {
		option(p)
	}
	return p
}

var _ view.Viewer = (*Viewer)(nil)

type Viewer struct {
	module  *mod.Module
	embed   bool
	minify  bool
	liveUrl string
	env     any
}

func (v *Viewer) Render(w io.Writer, path string, data *view.Data) error {
	html, err := v.evaluateSSR(path, data)
	if err != nil {
		return fmt.Errorf("preact: unable to evaluate html for %q. %w", path, err)
	}
	_, err = w.Write(html)
	return err
}

// var _ gen.FileServer = (*Viewer)(nil)

// func (v *Viewer) ServeFile(fsys gen.FS, file *gen.File) error {
// 	return nil
// }

func (v *Viewer) RenderHTML(w http.ResponseWriter, path string, data *view.Data) error {
	w.Header().Set("Content-Type", "text/html")
	return v.Render(w, path, data)
}

// RenderJS renders the client-side for hydrating the server-side rendered HTML.
func (v *Viewer) RenderJS(w http.ResponseWriter, path string, data *view.Data) error {
	entry, err := v.CompileDOM("./" + path)
	if err != nil {
		return fmt.Errorf("preact: unable to compile dom %q. %w", path, err)
	}
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(entry.Contents)
	return nil
}

// func (v *Viewer) Routes(router mux.Router) {
// 	router.Get("/.preact/{path*}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		viewPath := strings.TrimSuffix(r.URL.Query().Get("path"), ".js")
// 		if viewPath == "" {
// 			http.Error(w, fmt.Sprintf("preact: client view %q not found", r.URL.Path), http.StatusNotFound)
// 			return
// 		}
// 		entry, err := v.CompileDOM("./" + viewPath)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("preact: unable to compile %q. %s", viewPath, err), http.StatusInternalServerError)
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/javascript")
// 		w.Write(entry.Contents)
// 	}))
// }

func (v *Viewer) evaluateSSR(path string, data *view.Data) ([]byte, error) {
	entry, err := v.CompileSSR("./" + path)
	if err != nil {
		return nil, err
	}
	// TODO: optimize
	program, err := goja.Compile(path, string(entry.Contents), false)
	if err != nil {
		return nil, err
	}
	vm := js.New()
	_, err = vm.RunProgram(program)
	if err != nil {
		return nil, err
	}
	props, err := json.Marshal(data.Props)
	if err != nil {
		return nil, err
	}
	result, err := vm.RunString(fmt.Sprintf("bud.render(%s, { liveUrl: %q })", props, v.liveUrl))
	if err != nil {
		return nil, err
	}
	var ssr SSR
	err = json.Unmarshal([]byte(result.String()), &ssr)
	if err != nil {
		return nil, err
	}
	if data.Slots != nil {
		if ssr.Heads != nil {
			data.Slots.Slot("heads").Write(ssr.Heads)
		}
	}
	return []byte(ssr.HTML), nil
}

type SSR struct {
	HTML  string          `json:"html"`
	Heads json.RawMessage `json:"heads"`
}

// func (v *Viewer) Middleware(next http.Handler) http.Handler {
// 	return view.Middleware(v).Middleware(next)
// }

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

var ssrEntryTemplate = template.Must(template.New("ssr_entry.gotext").Parse(ssrEntry))

func (v *Viewer) ssrEntry(entry string) func(esbuild.OnLoadArgs) (esbuild.OnLoadResult, error) {
	return func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
		resolveDir := filepath.Dir(args.Path)
		relPath, err := filepath.Rel(resolveDir, entry)
		if err != nil {
			return result, err
		}
		props := map[string]interface{}{
			"Entry": relPath,
		}
		code := new(strings.Builder)
		if err := ssrEntryTemplate.Execute(code, props); err != nil {
			return result, err
		}
		contents := code.String()
		result.Contents = &contents
		result.ResolveDir = resolveDir
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
	ssrEntry := spliceExt(entry, ".ssr")
	options := esbuild.BuildOptions{
		AbsWorkingDir:   v.module.Directory(),
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
			esb.Env("bud/env.json", v.env),
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

// TODO: can we just check if it's not absolute and make this an implementation
// detail inside CompileSSR?
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

var domEntryTemplate = template.Must(template.New("dom_entry.gotext").Parse(domEntry))

func (v *Viewer) domEntry(entry string) func(esbuild.OnLoadArgs) (esbuild.OnLoadResult, error) {
	return func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
		resolveDir := filepath.Dir(args.Path)
		relPath, err := filepath.Rel(resolveDir, entry)
		if err != nil {
			return result, err
		}
		props := map[string]interface{}{
			"Entry": relPath,
		}
		code := new(strings.Builder)
		if err := domEntryTemplate.Execute(code, props); err != nil {
			return result, err
		}
		contents := code.String()
		result.Contents = &contents
		result.ResolveDir = filepath.Dir(args.Path)
		result.Loader = esbuild.LoaderTSX
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
			esb.Env("bud/env.json", v.env),
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
