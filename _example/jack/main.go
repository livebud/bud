package main

import (
	"fmt"
	"net/http"

	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/view"
	"github.com/livebud/bud/pkg/view/preact"
)

func main() {
	log := logs.Default()
	module := mod.MustFind()
	router := mux.New()
	preact := preact.New(module)
	// publicFS, err := fs.Sub(module.FS, "public")
	// if err != nil {
	// 	log.Error(err)
	// 	os.Exit(1)
	// }
	router.Get("/{path*}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		http.FileServer(http.Dir("public")).ServeHTTP(w, r)
	}))
	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := preact.Render(w, "view/index.tsx", &view.Data{}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))
	router.Get("/index.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		preact.RenderJS(w, "view/index.tsx", &view.Data{})
	}))
	log.Infof("Listening on http://localhost%s", ":8080")
	http.ListenAndServe(":8080", router)
}
