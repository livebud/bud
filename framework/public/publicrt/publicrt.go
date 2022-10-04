package publicrt

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"time"
)

type FS = fs.FS

func NewHandler(fsys FS) *Handler {
	return &Handler{http.FS(fsys)}
}

type Handler struct {
	fsys http.FileSystem
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, err := h.fsys.Open(path.Join("public", r.URL.Path))
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
	serveContent(w, r, r.URL.Path, stat.ModTime(), file)
}

func serveContent(w http.ResponseWriter, req *http.Request, name string, modtime time.Time, content io.ReadSeeker) {
	http.ServeContent(w, req, name, modtime, content)
}
