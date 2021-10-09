package entrypoint

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/matthewmueller/text"
)

// List the views
func List(fsys fs.FS) ([]*View, error) {
	views, err := list(fsys, ".", &reserved{
		Error:  map[string]Path{},
		Layout: map[string]Path{},
		Frames: map[string][]Path{},
	})
	if err != nil {
		return nil, err
	}
	// Sort the views
	sort.Slice(views, func(i, j int) bool {
		return views[i].Page < views[j].Page
	})
	return views, nil
}

// reserved views
type reserved struct {
	Error  map[string]Path
	Layout map[string]Path
	Frames map[string][]Path
}

func list(fsys fs.FS, dir string, reserved *reserved) (views []*View, err error) {
	fis, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}
	var remaining []fs.DirEntry
	// var errorPaths []string
	for _, fi := range fis {
		name := fi.Name()
		ext := filepath.Ext(name)
		if fi.IsDir() {
			remaining = append(remaining, fi)
			continue
		} else if !validEntry(ext) {
			continue
		}
		path := filepath.Join(dir, name)
		switch extless(name) {
		case "error":
			// errorPaths = append(errorPaths, name)
			reserved.Error[ext] = Path(path)
			continue
		case "layout":
			reserved.Layout[ext] = Path(path)
			continue
		case "frame":
			reserved.Frames[ext] = append(reserved.Frames[ext], Path(path))
			continue
		}
		// Add to the remaining for additional inspection
		remaining = append(remaining, fi)
	}
	// Second, loop over the directory adding views as we go and traversing
	// directories
	for _, fi := range remaining {
		name := fi.Name()
		base := filepath.Base(name)
		// Handle directories
		if fi.IsDir() {
			if len(base) == 0 || base[0] == '_' {
				continue
			}
			subdir := filepath.Join(dir, name)
			subviews, err := list(fsys, subdir, reserved)
			if err != nil {
				return nil, err
			}
			views = append(views, subviews...)
			continue
		}
		// Handle files
		// Ignore capitalized files and files that start with an underscore
		if len(base) == 0 || base[0] == '_' || base[0] == '.' || unicode.IsUpper(firstRune(base)) {
			continue
		}
		path := filepath.Join(dir, name)
		ext := filepath.Ext(name)
		views = append(views, &View{
			Page:   Path(path),
			Client: client(path),
			Route:  route(dir, name),
			Frames: reserved.Frames[ext],
			Layout: Path(reserved.Layout[ext]),
			Type:   ext[1:],
			Hot:    true, // TODO: remove
		})
	}
	return views, nil
}

func validEntry(ext string) bool {
	switch ext {
	case ".svelte", ".jsx", ".tsx", ".gohtml":
		return true
	default:
		return false
	}
}

func firstRune(s string) rune {
	r, _ := utf8.DecodeRuneInString(s)
	return r
}

// Path is the route to the action
func route(dir, name string) string {
	dir = strings.TrimPrefix(dir, "view")
	if dir == "." {
		dir = ""
	}
	base := extless(name)
	switch base {
	case "show":
		return "/" + path.Join(dir, ":id")
	case "new":
		return "/" + path.Join(dir, "new")
	case "edit":
		return "/" + path.Join(dir, ":id", "edit")
	case "index":
		return "/" + dir
	default:
		return "/" + path.Join(dir, text.Path(base))
	}
}

// The client path's entrypoint
//
//   e.g. view/index.jsx => bud/view/_index.jsx
//
// We use the _ because this is a generated entrypoint
// that
func client(name string) string {
	dir, path := filepath.Split(name)
	return fmt.Sprintf("bud/%s_%s", dir, path)
}
