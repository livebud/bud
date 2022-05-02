package welcome

import (
	"embed"
	_ "embed"
	"io/fs"
	"net/http"

	"github.com/livebud/bud/package/middleware"
)

// Files are built in https://github.com/livebud/welcome and manually copied
// over.
//go:embed build/index.html build/default.css build/bud/view/_index.svelte.js
var embeds embed.FS

func Load() (Middleware, error) {
	fsys, err := fs.Sub(embeds, "build")
	if err != nil {
		return nil, err
	}
	server := http.FileServer(http.FS(fsys))
	return middleware.Function(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			server.ServeHTTP(w, r)
		})
	}), nil
}

type Middleware = middleware.Middleware
