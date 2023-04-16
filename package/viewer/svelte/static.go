package svelte

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/router"
	"github.com/livebud/bud/package/viewer"
	"github.com/livebud/bud/package/virtual"
	"github.com/livebud/js"
)

func Static(fsys fs.FS, js js.VM, log log.Log, pages viewer.Pages) *StaticViewer {
	return &StaticViewer{fsys, http.FS(fsys), js, log, pages}
}

type StaticViewer struct {
	fsys  fs.FS
	hfs   http.FileSystem
	js    js.VM
	log   log.Log
	pages viewer.Pages
}

var _ viewer.Viewer = (*StaticViewer)(nil)

func (v *StaticViewer) Mount(r *router.Router) error {
	for _, page := range v.pages {
		// Serve the entrypoints (for hydrating)
		r.Get(page.Client.Route, v.serveDOMEntry(page))
		// Serve the individual views themselves
		r.Get(page.View.Client.Route, v.serveDOMView(page.View))
	}
	return nil
}

func (v *StaticViewer) Render(ctx context.Context, key string, propMap viewer.PropMap) ([]byte, error) {
	page, ok := v.pages[key]
	if !ok {
		return nil, fmt.Errorf("svelte: unable to find page from key %q", key)
	}
	v.log.Info("svelte: rendering", page.Path)
	code, err := fs.ReadFile(v.fsys, page.Path)
	if err != nil {
		return nil, err
	}
	propBytes, err := json.Marshal(propMap)
	if err != nil {
		return nil, err
	}
	expr := fmt.Sprintf(`%s; bud.render(%s)`, string(code), propBytes)
	html, err := v.js.Evaluate(ctx, page.Path, expr)
	if err != nil {
		return nil, err
	}
	return []byte(html), nil
}

func (v *StaticViewer) RenderError(ctx context.Context, key string, propMap viewer.PropMap, originalError error) []byte {
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
	code, err := fs.ReadFile(v.fsys, errorPage.Path)
	if err != nil {
		return []byte(fmt.Sprintf("svelte: unable to read error page %q code to render error. %s. %s", errorPage.Path, err, originalError))
	}
	propMap[errorPage.Key] = viewer.Error(originalError)
	propBytes, err := json.Marshal(propMap)
	if err != nil {
		return []byte(fmt.Sprintf("svelte: unable to marshal props for %q to render error. %s. %s", errorPage.Path, err, originalError))
	}
	expr := fmt.Sprintf(`%s; bud.render(%s)`, code, propBytes)
	html, err := v.js.Evaluate(ctx, errorPage.Path, expr)
	if err != nil {
		return []byte(fmt.Sprintf("svelte: unable to evaluate javascript to render %q to render error. %s. %s", errorPage.Path, err, originalError))
	}
	return []byte(html)
}

// TODO: split Viewer into Bundler and Renderer interfaces.
func (v *StaticViewer) Bundle(ctx context.Context, embed virtual.Tree) error {
	return fmt.Errorf("can't bundle with the static viewer")
}

func (v *StaticViewer) Handler(page *viewer.Page) http.Handler {
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

func (v *StaticViewer) renderError(ctx context.Context, w http.ResponseWriter, page *viewer.Page, propMap map[string]interface{}, err error) {
	html := v.RenderError(ctx, page.Key, propMap, err)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(html)
}

func (v *StaticViewer) serveDOMEntry(page *viewer.Page) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.log.Info("svelte: serving dom entry", page.Client.Path)
		file, err := v.hfs.Open(page.Client.Path)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer file.Close()
		stat, err := file.Stat()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if stat.IsDir() {
			http.Error(w, fmt.Sprintf("%q is a directory", r.URL.Path), 500)
			return
		}
		http.ServeContent(w, r, page.Client.Path, stat.ModTime(), file)
	})
}

func (v *StaticViewer) serveDOMView(view *viewer.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.log.Info("svelte: serving dom view", view.Client.Path)
		file, err := v.hfs.Open(view.Client.Path)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer file.Close()
		stat, err := file.Stat()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if stat.IsDir() {
			http.Error(w, fmt.Sprintf("%q is a directory", r.URL.Path), 500)
			return
		}
		http.ServeContent(w, r, view.Client.Path, stat.ModTime(), file)
	})
}
