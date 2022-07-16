package entrypoint

import (
	"encoding/base64"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/matthewmueller/gotext"
)

type View struct {
	Page   Path   // Path to the page
	Type   string // View extension
	Route  string
	Frames []Path
	Layout Path
	Error  Path
	Client string
	Hot    string
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
