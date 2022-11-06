package welcome

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/livebud/bud/package/router"
	"github.com/livebud/bud/package/virtual"
)

// Files are built in https://github.com/livebud/welcome and manually copied
// over.

//go:embed build/index.html
var index []byte

//go:embed build/bud/view/_index.svelte.js
var clientJS []byte

type Handler struct {
}

func (h *Handler) Register(r *router.Router) {
	handle := handler(http.FS(&virtual.Map{
		".": &virtual.File{
			Path: "index.html",
			Data: index,
		},
		"bud/view/_index.svelte.js": &virtual.File{
			Path: "bud/view/_index.svelte.js",
			Data: clientJS,
		},
	}))
	r.Get("/", handle)
	r.Get("/bud/view/_index.svelte.js", handle)
}

func handler(fsys http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := fsys.Open(r.URL.Path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), 500)
			return
		}
		defer file.Close()
		stat, err := file.Stat()
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		if stat.IsDir() {
			http.Error(w, fmt.Sprintf("%q is a directory", r.URL.Path), 500)
			return
		}
		http.ServeContent(w, r, r.URL.Path, stat.ModTime(), file)
	})
}
