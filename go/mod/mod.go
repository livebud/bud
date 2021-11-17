package mod

import (
	"errors"
	"fmt"
	"go/build"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/mnm/bud/internal/modcache"
	"golang.org/x/mod/modfile"
)

// ErrCantInfer occurs when you can't infer the module path from the $GOPATH.
var ErrCantInfer = errors.New("mod: unable to infer the module path")

// ErrFileNotFound occurs when no go.mod can be found
var ErrFileNotFound = fmt.Errorf("unable to find go.mod: %w", fs.ErrNotExist)

func Find() (*File, error) {
	return FindIn(modcache.Default(), ".")
}

// Find first tries finding an explicit module file (go.mod). If no go.mod is
// found, then Find will try inferring a virtual module file from $GOPATH.
func FindIn(cache *modcache.Cache, appDir string) (*File, error) {
	absdir, err := filepath.Abs(appDir)
	if err != nil {
		return nil, err
	}
	// First search for go.mod
	modfile, err := findModFile(cache, absdir)
	if nil == err {
		return modfile, nil
	} else if !errors.Is(err, ErrFileNotFound) {
		return nil, err
	}
	// If that fails, try inferring from the $GOPATH
	return Infer(cache, absdir)
}

// Parse a modfile from it's data
func Parse(cache *modcache.Cache, path string, data []byte) (*File, error) {
	modfile, err := modfile.Parse(path, data, nil)
	if err != nil {
		return nil, err
	}
	return &File{
		file:  modfile,
		cache: cache,
		dir:   filepath.Dir(path),
	}, nil
}

// Infer the module path from the $GOPATH. This only works if you work inside
// $GOPATH.
func Infer(cache *modcache.Cache, appDir string) (*File, error) {
	modulePath := modulePathFromGoPath(appDir)
	if modulePath == "" {
		return nil, fmt.Errorf("%w for %q, run `go mod init` to fix", ErrCantInfer, appDir)
	}
	virtualPath := filepath.Join(appDir, "go.mod")
	return Parse(cache, virtualPath, []byte("module "+modulePath))
}

// Find the go.mod file from anywhere in your project.
func findModFile(cache *modcache.Cache, path string) (*File, error) {
	moduleDir, err := findModPath(path)
	if err != nil {
		return nil, fmt.Errorf("%w in %q", ErrFileNotFound, path)
	}
	modulePath := filepath.Join(moduleDir, "go.mod")
	moduleData, err := ioutil.ReadFile(modulePath)
	if err != nil {
		return nil, err
	}
	return Parse(cache, modulePath, moduleData)
}

// findModPath traverses up the filesystem until it finds a directory containing
// go.mod or returns an error trying.
func findModPath(dir string) (abs string, err error) {
	path := filepath.Join(dir, "go.mod")
	// Check if this path exists, otherwise recursively traverse towards root
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) && dir != string(filepath.Separator) {
			return findModPath(filepath.Dir(dir))
		}
		return "", ErrFileNotFound
	}
	return dir, nil
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
