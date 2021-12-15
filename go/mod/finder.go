package mod

import (
	"errors"
	"fmt"
	"go/build"
	"io/fs"
	"path/filepath"
	"strings"

	"gitlab.com/mnm/bud/internal/modcache"
	"golang.org/x/mod/modfile"
)

// Finder struct
type Finder struct {
	fsys  fs.FS
	cache *modcache.Cache
}

// Find first tries finding an explicit module file (go.mod). If no go.mod is
// found, then Find will try inferring a virtual module file from $GOPATH.
// Fullpath should be an absolute path if you're using the default filesystem.
func (f *Finder) Find(fullpath string) (*Module, error) {
	// First search for go.mod
	modfile, err := f.findModFile(fullpath)
	if nil == err {
		return modfile, nil
	} else if !errors.Is(err, ErrFileNotFound) {
		return nil, err
	}
	// If that fails, try inferring from the $GOPATH
	return f.Infer(fullpath)
}

// Infer the module path from the $GOPATH. This only works if you work inside
// $GOPATH.
func (f *Finder) Infer(dir string) (*Module, error) {
	modulePath := modulePathFromGoPath(dir)
	if modulePath == "" {
		return nil, fmt.Errorf("%w for %q, run `go mod init` to fix", ErrCantInfer, dir)
	}
	virtualPath := filepath.Join(dir, "go.mod")
	return f.parse(virtualPath, []byte("module "+modulePath))
}

// Parse a modfile from it's data
func (f *Finder) Parse(path string, data []byte) (*Module, error) {
	return f.parse(path, data)
}

func (f *Finder) parse(path string, data []byte) (*Module, error) {
	modfile, err := modfile.Parse(path, data, nil)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(path)
	fsys, err := subFS(f.fsys, dir)
	if err != nil {
		return nil, err
	}
	return &Module{
		file:  &File{modfile},
		cache: f.cache,
		fsys:  fsys,
		dir:   dir,
	}, nil
}

// Find the go.mod file from anywhere in your project.
func (f *Finder) findModFile(path string) (*Module, error) {
	moduleDir, err := findDirectory(f.fsys, path)
	if err != nil {
		return nil, fmt.Errorf("%w in %q", ErrFileNotFound, path)
	}
	modulePath := filepath.Join(moduleDir, "go.mod")
	moduleData, err := fs.ReadFile(f.fsys, modulePath)
	if err != nil {
		return nil, err
	}
	return f.parse(modulePath, moduleData)
}

// gopathToModulePath tries inferring the module path of directory. This only
// works if you're in working within the $GOPATH
func modulePathFromGoPath(path string) string {
	src := filepath.Join(build.Default.GOPATH, "src") + "/"
	if !strings.HasPrefix(path, src) {
		return ""
	}
	modulePath := strings.TrimPrefix(path, src)
	return modulePath
}
