package publicrt

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"path"
	"time"

	"github.com/livebud/bud/package/budhttp"
	"github.com/livebud/bud/package/middleware"

	"github.com/livebud/bud/package/budfs/mergefs"
)

type Server interface {
	Serve(fsys fs.FS) middleware.Middleware
}

func Live(client budhttp.Client) (*LiveServer, error) {
	return &LiveServer{client}, nil
}

type LiveServer struct {
	fsys fs.FS
}

func (l *LiveServer) Serve(fsys fs.FS) middleware.Middleware {
	fsys = mergefs.Merge(l.fsys, fsys)
	return serve(fsys, serveContent)
}

func Static() StaticServer {
	return StaticServer{}
}

type StaticServer struct{}

// TODO: serve gzipped content
func (StaticServer) Serve(fsys fs.FS) middleware.Middleware {
	return serve(fsys, serveContent)
}

func serve(fsys fs.FS, serveContent func(w http.ResponseWriter, req *http.Request, name string, modtime time.Time, content io.ReadSeeker)) middleware.Middleware {
	hfs := http.FS(fsys)
	return middleware.Function(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlPath := r.URL.Path
			if r.Method != http.MethodGet || path.Ext(urlPath) == "" {
				next.ServeHTTP(w, r)
				return
			}
			file, err := hfs.Open(path.Join("public", urlPath))
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					next.ServeHTTP(w, r)
					return
				}
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
				next.ServeHTTP(w, r)
				return
			}
			serveContent(w, r, urlPath, stat.ModTime(), file)
		})
	})
}

func serveContent(w http.ResponseWriter, req *http.Request, name string, modtime time.Time, content io.ReadSeeker) {
	http.ServeContent(w, req, name, modtime, content)
}
