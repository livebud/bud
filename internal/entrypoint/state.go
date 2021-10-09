package entrypoint

import (
	"encoding/base64"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/go-duo/bud/internal/gotext"
)

// type State struct {
// 	Views []*View
// }

// func (s *State) Pages(suffixes ...string) (pages []string) {
// 	suffix := strings.Join(suffixes, "")
// 	for _, view := range s.Views {
// 		pages = append(pages, string(view.Page)+suffix)
// 	}
// 	return pages
// }

// func (s *State) Imports() (imports []Path, err error) {
// 	for _, view := range s.Views {
// 		imports = append(imports, view.Page)
// 	}
// 	return imports, nil
// }

// func (s *State) FindViewByPage(page string) (*View, error) {
// 	page = filepath.Clean(page)
// 	for _, view := range s.Views {
// 		if string(view.Page) == page {
// 			return view, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("state: no view with a page named %q", page)
// }

type View struct {
	Page   Path   // Path to the page
	Type   string // View extension
	Route  string
	Frames []Path
	Layout Path
	Error  Path
	Client string
	Hot    bool
}

func (v *View) ServerImports() (imports []Path) {
	imports = append(imports, v.Page)
	imports = append(imports, v.Frames...)
	if v.Layout != "" {
		imports = append(imports, v.Layout)
	}
	if v.Error != "" {
		imports = append(imports, v.Error)
	}
	return imports
}

func (v *View) BrowserImports() (imports []Path) {
	imports = append(imports, v.Page)
	imports = append(imports, v.Frames...)
	if v.Error != "" {
		imports = append(imports, v.Error)
	}
	return imports
}

// View data as a query string for the hot page
func (v *View) Query() string {
	values := url.Values{}
	values.Set("page", "/bud/"+string(v.Page))
	return values.Encode()
}

type Path string

func (p Path) Pascal() string {
	return gotext.Pascal(string(p))
}

func (p Path) Camel() string {
	return gotext.Camel(string(p))
}

func (p Path) Route() string {
	name := strings.TrimPrefix(filepath.ToSlash(extless(string(p))), "view/")
	dir, base := path.Split(name)
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
		return "/" + path.Join(dir, strings.ToLower(base))
	}
}

func (path Path) Ext() string {
	return filepath.Ext(string(path))[1:]
}

func (path Path) Type() string {
	switch extless(filepath.Base(string(path))) {
	case "layout":
		return "layout"
	case "frame":
		return "frame"
	case "error":
		return "error"
	default:
		return "page"
	}
}

func (path Path) Hash() (string, error) {
	contents, err := os.ReadFile(string(path))
	if err != nil {
		return "", err
	}
	hash := xxhash.New()
	hash.Write(contents)
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil)), nil
}

func (path Path) Layout() bool {
	return extless(filepath.Base(string(path))) == "layout"
}

func (path Path) Frame() bool {
	return extless(filepath.Base(string(path))) == "frame"
}

func (path Path) Error() bool {
	return extless(filepath.Base(string(path))) == "error"
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
