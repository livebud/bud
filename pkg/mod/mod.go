package mod

import (
	"errors"
	"fmt"
	"go/build"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

// ErrFileNotFound occurs when no go.mod can be found
var ErrFileNotFound = fmt.Errorf(`mod: unable to find "go.mod": %w`, fs.ErrNotExist)

// Build with:
//
//	go build -trimpath -ldflags " -X github.com/livebud/mod.path=$(go list -m) -X github.com/livebud/mod.dir=$(go list -m -f {{.Dir}})" ./
var (
	path       string
	dir        string
	isEmbedded = path != "" && dir != ""
)

// New module
func New(dir string) *Module {
	modulePath := modulePathFromGoPath(dir)
	if modulePath == "" {
		modulePath = "change.me"
	}
	module, err := Parse(filepath.Join(dir, "go.mod"), []byte(`module `+modulePath))
	if err != nil {
		panic("mod: invalid module data: " + err.Error())
	}
	return module
}

// Find the first go.mod file in one of the directories below or return an
// error. Find will also search parent directories for a go.mod file.
func Find(dirs ...string) (*Module, error) {
	if len(dirs) == 0 {
		return find(".")
	}
	for _, dir := range dirs {
		module, err := find(dir)
		if err != nil {
			if !errors.Is(err, ErrFileNotFound) {
				return nil, err
			}
			continue
		}
		return module, nil
	}
	return nil, ErrFileNotFound
}

func find(dir string) (*Module, error) {
	if isEmbedded {
		return Parse(filepath.Join(dir, "go.mod"), []byte(`module `+path))
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	modPath, err := lookup(abs)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(modPath)
	if err != nil {
		return nil, fmt.Errorf(`mod: unable to read "go.mod": %w`, err)
	}
	return Parse(modPath, data)
}

// MustFind a go.mod file in this directory or any parent directory. If no
// go.mod file is found, this will panic.
func MustFind(dirs ...string) *Module {
	module, err := Find(dirs...)
	if err != nil {
		panic(err)
	}
	return module
}

// Parse a go.mod file
func Parse(path string, data []byte) (*Module, error) {
	modfile, err := modfile.Parse(path, data, nil)
	if err != nil {
		return nil, err
	}
	if modfile.Module == nil {
		modFile, err := modfile.Format()
		if err != nil {
			return nil, fmt.Errorf("mod: missing module statement in %q and got an error while formatting %s", path, err)
		}
		return nil, fmt.Errorf("mod: missing module statement in %q, received %q", path, string(modFile))
	}
	dir := filepath.Dir(path)
	return &Module{dir, modfile, os.DirFS(dir)}, nil
}

// Lookup finds the absolute path of the go.mod file in the given directory
func Lookup(directory string) (path string, err error) {
	if dir != "" {
		return filepath.Join(dir, "go.mod"), nil
	}
	path, err = lookup(directory)
	if err != nil {
		return "", err
	}
	return filepath.Abs(path)
}

func lookup(directory string) (path string, err error) {
	path = filepath.Join(directory, "go.mod")
	// Check if this path exists, otherwise recursively traverse towards root
	if _, err = os.Stat(path); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		nextDir := filepath.Dir(directory)
		if nextDir == directory {
			return "", ErrFileNotFound
		}
		return lookup(nextDir)
	}
	return filepath.EvalSymlinks(path)
}

// modulePathFromGoPath tries inferring the module path of directory. This only
// works if you're in working within the $GOPATH
func modulePathFromGoPath(path string) string {
	src := filepath.Join(build.Default.GOPATH, "src") + "/"
	if !strings.HasPrefix(path, src) {
		return ""
	}
	modulePath := strings.TrimPrefix(path, src)
	return modulePath
}
