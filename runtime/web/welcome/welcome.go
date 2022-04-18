package welcome

import (
	"bytes"
	_ "embed"
	"net/http"
	"time"

	"github.com/livebud/bud/package/middleware"
	"github.com/livebud/bud/package/router"
)

//go:generate go run github.com/evanw/esbuild/cmd/esbuild --bundle --log-level=error --outfile=index.out.css index.css
//go:generate go run github.com/evanw/esbuild/cmd/esbuild --bundle --log-level=error --outfile=index.out.js index.js

// Compute the modTime once when loaded
// TODO: can we do better here?
var modTime = time.Now()

func New() Middleware {
	router := router.New()
	router.Get("/", http.HandlerFunc(serveHTML))
	router.Get("/index.css", http.HandlerFunc(serveCSS))
	router.Get("/index.js", http.HandlerFunc(serveJS))
	return router
}

type Middleware = middleware.Middleware

//go:embed index.html
var indexHtml []byte

func serveHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, "index.html", modTime, bytes.NewReader(indexHtml))
}

//go:embed index.out.css
var indexCSS []byte

func serveCSS(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, "index.out.css", modTime, bytes.NewReader(indexCSS))
}

//go:embed index.out.js
var indexJS []byte

func serveJS(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, "index.out.js", modTime, bytes.NewReader(indexJS))
}
