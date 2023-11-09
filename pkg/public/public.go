package public

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
)

func New(fsys fs.FS) *FileServer {
	return &FileServer{
		hfs: http.FS(fsys),
	}
}

type FileServer struct {
	hfs http.FileSystem
}

// Middleware serves files. If a file doesn't exist, we move on to the next
// handler.
func (f *FileServer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		file, err := f.hfs.Open(r.URL.Path)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		defer file.Close()
		stat, err := file.Stat()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// TODO: consider optionally serving directories
		if stat.IsDir() {
			http.Error(w, fmt.Sprintf("%s is a directory", r.URL.Path), http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
	})
}

func (f *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.Middleware(http.NotFoundHandler()).ServeHTTP(w, r)
}
