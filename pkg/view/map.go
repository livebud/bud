package view

import (
	"fmt"
	"io"
	"path"
)

type Map map[string]Viewer

func (m Map) Render(w io.Writer, viewPath string, data *Data) error {
	viewer, ok := m[path.Ext(viewPath)]
	if !ok {
		return fmt.Errorf("view: no viewer for path %q", viewPath)
	}
	return viewer.Render(w, viewPath, data)
}

// func (m Map) Routes(router mux.Router) {
// 	for _, viewer := range m {
// 		viewer.Routes(router)
// 	}
// }

// func (m Map) Middleware(next http.Handler) http.Handler {
// 	return Middleware(m).Middleware(next)
// }
