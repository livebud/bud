package svelte

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/livebud/bud/package/js"
	"github.com/livebud/bud/package/router"
)

func New(fsys fs.FS, vm js.VM) *Viewer {
	return &Viewer{fsys, vm}
}

type Viewer struct {
	fsys fs.FS
	vm   js.VM
}

func (v *Viewer) Handler(route string, props interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.Render(w, route, props)
	})
}

func (v *Viewer) Render(w http.ResponseWriter, viewPath string, props interface{}) {
	script, err := fs.ReadFile(v.fsys, "bud/view/_ssr.js")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	propBytes, err := json.Marshal(props)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := v.vm.Eval("_ssr.js", fmt.Sprintf(`%s; bud.render(%q, %s)`, script, viewPath, propBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("got result", result)
	return
}

func (v *Viewer) Serve(r *router.Router) {
	// Serving files!
}
