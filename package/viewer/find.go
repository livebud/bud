package viewer

import (
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/package/valid"
	"github.com/matthewmueller/text"

	"github.com/livebud/bud/internal/gitignore"
)

// Find pages
func Find(fsys FS) (pages map[Key]*Page, err error) {
	ignore := gitignore.FromFS(fsys)
	pages = make(map[Key]*Page)
	inherited := &inherited{
		Layout: make(map[ext]*View),
		Frames: make(map[ext][]*View),
		Error:  make(map[ext]*View),
	}
	if err := find(fsys, ignore, pages, inherited, "."); err != nil {
		return nil, err
	}
	return pages, nil
}

type ext = string

type inherited struct {
	Layout map[ext]*View
	Frames map[ext][]*View
	Error  map[ext]*View
}

func find(fsys FS, ignore func(path string) bool, pages map[Key]*Page, inherited *inherited, dir string) error {
	des, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return err
	}

	// First pass: look for layouts, frames and errors
	for _, de := range des {
		if de.IsDir() {
			continue
		}
		fpath := path.Join(dir, de.Name())
		if ignore(fpath) {
			continue
		}
		ext := filepath.Ext(de.Name())
		extless := extless(de.Name())
		key := path.Join(dir, extless)
		switch extless {
		case "layout":
			inherited.Layout[ext] = &View{
				Path: fpath,
				Key:  key,
				Ext:  ext,
			}
		case "frame":
			inherited.Frames[ext] = append([]*View{&View{
				Path:   fpath,
				Key:    key,
				Ext:    ext,
				Client: viewClient(fpath),
			}}, inherited.Frames[ext]...)
		case "error":
			inherited.Error[ext] = &View{
				Path:   fpath,
				Key:    key,
				Ext:    ext,
				Client: viewClient(fpath),
			}
		}
	}

	// Second pass: go through pages
	for _, de := range des {
		if de.IsDir() {
			continue
		}
		ext := filepath.Ext(de.Name())
		if !valid.View(de.Name()) {
			continue
		}
		extless := extless(de.Name())
		switch extless {
		case "layout", "frame":
			continue
		// Errors are treated just like regular pages with frames and layouts
		case "error":
			key := path.Join(dir, extless)
			fpath := path.Join(dir, de.Name())
			pages[key] = &Page{
				View: &View{
					Path:   fpath,
					Key:    key,
					Ext:    ext,
					Client: viewClient(fpath),
				},
				Layout: inherited.Layout[ext],
				Frames: inherited.Frames[ext],
				Error:  nil, // Error pages can't have their own error page
				Route:  route(dir, extless),
				Client: entryClient(fpath),
			}
		default:
			key := path.Join(dir, extless)
			fpath := path.Join(dir, de.Name())
			pages[key] = &Page{
				View: &View{
					Path:   fpath,
					Key:    key,
					Ext:    ext,
					Client: viewClient(fpath),
				},
				Layout: inherited.Layout[ext],
				Frames: inherited.Frames[ext],
				Error:  inherited.Error[ext],
				Route:  route(dir, extless),
				Client: entryClient(fpath),
			}
		}
	}

	// Third pass: go through directories
	for _, de := range des {
		if !de.IsDir() {
			continue
		}
		fpath := path.Join(dir, de.Name())
		if ignore(fpath) {
			continue
		}
		if err := find(fsys, ignore, pages, inherited, fpath); err != nil {
			return err
		}
	}

	return nil
}

// Generate the IDs for a nested route
// TODO: consolidate with the function in internal/generator/action/loader.go.
func routeDir(dir string) string {
	segments := strings.Split(dir, "/")
	path := new(strings.Builder)
	for i := 0; i < len(segments); i++ {
		if i%2 != 0 {
			path.WriteString("/")
			path.WriteString(":" + text.Snake(text.Singular(segments[i-1])) + "_id")
			path.WriteString("/")
		}
		path.WriteString(text.Snake(segments[i]))
	}
	if path.Len() == 0 {
		return ""
	}
	return path.String()
}

// Path is the route to the action
func route(dir, extless string) string {
	dir = strings.TrimPrefix(strings.TrimPrefix(dir, "view"), "/")
	if dir == "." {
		dir = ""
	}
	dir = routeDir(dir)
	switch extless {
	case "show":
		return "/" + path.Join(dir, ":id")
	case "new":
		return "/" + path.Join(dir, "new")
	case "edit":
		return "/" + path.Join(dir, ":id", "edit")
	case "index":
		return "/" + dir
	default:
		return "/" + path.Join(dir, text.Lower(text.Snake(extless)))
	}
}

// Recursively trim file extensions until there aren't any left
func extless(path string) string {
	ext := filepath.Ext(path)
	for ext != "" {
		path = strings.TrimSuffix(path, ext)
		ext = filepath.Ext(path)
	}
	return path
}

func viewClient(fpath string) *Client {
	viewPath := path.Clean(fpath) + ".js"
	return &Client{
		Path:  viewPath,
		Route: "/view/" + viewPath,
	}
}

func entryClient(fpath string) *Client {
	entryPath := path.Clean(fpath) + ".entry.js"
	return &Client{
		Path:  entryPath,
		Route: "/view/" + entryPath,
	}
}
