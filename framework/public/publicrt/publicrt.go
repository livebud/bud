package publicrt

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"path"
	"time"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/middleware"
	"github.com/livebud/bud/package/pluginmod"

	"github.com/livebud/bud/package/budfs/mergefs"
)

// Load the app and plugins into a single filesystem
func LoadFS(module *gomod.Module) (fsys fs.FS, err error) {
	publics, err := pluginmod.Glob(module, "public")
	if err != nil {
		return nil, err
	}
	// We'll still want to render public for the default favicon if
	// we have views.
	// TODO: this will go away when we scaffold a default favicon and then remove
	// auto-generating a default favicon when no favicon exists.
	if len(publics) == 0 {
		views, err := pluginmod.Glob(module, "view")
		if err != nil {
			return nil, err
		}
		if len(views) == 0 {
			return nil, fs.ErrNotExist
		}
	}
	// Merge the public modules into a single fs
	fileSystems := make([]fs.FS, len(publics))
	for i, public := range publics {
		fileSystems[i] = public
	}
	fsys = mergefs.Merge(fileSystems...)
	// Scope to the public/ directory
	return fs.Sub(fsys, "public")
}

type Server interface {
	Serve(fsys fs.FS) middleware.Middleware
}

func Live(module *gomod.Module) (*LiveServer, error) {
	fsys, err := LoadFS(module)
	if err != nil {
		return nil, err
	}
	return &LiveServer{fsys}, nil
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
			file, err := hfs.Open(urlPath)
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
