package gomod

import (
	"errors"
	"fmt"
	"go/build"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/internal/current"
	"github.com/livebud/bud/package/modcache"
	"github.com/livebud/bud/package/virtual"
	"golang.org/x/mod/modfile"
)

// ErrFileNotFound occurs when no go.mod can be found
var ErrFileNotFound = fmt.Errorf("unable to find go.mod: %w", fs.ErrNotExist)

type Option = func(o *option)

type option struct {
	modCache *modcache.Cache
}

// WithModCache uses a custom mod cache instead of the default
func WithModCache(cache *modcache.Cache) func(o *option) {
	return func(opt *option) {
		opt.modCache = cache
	}
}

func Find(dir string, options ...Option) (*Module, error) {
	opt := &option{
		modCache: modcache.Default(),
	}
	for _, option := range options {
		option(opt)
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	return find(opt, abs)
}

func find(opt *option, dir string) (*Module, error) {
	moduleDir, err := Absolute(dir)
	if err != nil {
		return nil, fmt.Errorf("module: %q %w", dir, ErrFileNotFound)
	}
	modulePath := filepath.Join(moduleDir, "go.mod")
	moduleData, err := os.ReadFile(modulePath)
	if err != nil {
		return nil, err
	}
	return parse(opt, modulePath, moduleData)
}

// MustFind is like Find but panics if an error occurs
func MustFind(dir string, options ...Option) *Module {
	module, err := Find(dir, options...)
	if err != nil {
		panic(err)
	}
	return module
}

// Infer the module path from the $GOPATH. This only works if you work inside
// $GOPATH.
func Infer(dir string) string {
	return modulePathFromGoPath(dir)
}

// Parse a modfile from it's data
func Parse(path string, data []byte, options ...Option) (*Module, error) {
	opt := &option{
		modCache: modcache.Default(),
	}
	for _, option := range options {
		option(opt)
	}
	return parse(opt, path, data)
}

func New(dir string, options ...Option) *Module {
	opt := &option{
		modCache: modcache.Default(),
	}
	for _, option := range options {
		option(opt)
	}
	modulePath := modulePathFromGoPath(dir)
	if modulePath == "" {
		modulePath = "change.me"
	}
	module, err := parse(opt, filepath.Join(dir, "go.mod"), []byte(`module `+modulePath))
	if err != nil {
		panic("mod: invalid module data: " + err.Error())
	}
	return module
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

func parse(opt *option, path string, data []byte) (*Module, error) {
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
	return &Module{opt, &File{modfile}, dir, virtual.OS(dir)}, nil
}

// Absolute traverses up the filesystem until it finds a directory
// containing go.mod or returns an error trying.
func Absolute(dir string) (abs string, err error) {
	dir, err = absolute(dir)
	if err != nil {
		return "", err
	}
	return filepath.Abs(dir)
}

// FindBudModule finds the go.mod for bud itself
func FindBudModule() (*Module, error) {
	dirname, err := current.Directory()
	if err != nil {
		return nil, err
	}
	return Find(dirname)
}

func absolute(dir string) (abs string, err error) {
	path := filepath.Join(dir, "go.mod")
	// Check if this path exists, otherwise recursively traverse towards root
	if _, err = os.Stat(path); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		nextDir := filepath.Dir(dir)
		if nextDir == dir {
			return "", ErrFileNotFound
		}
		return absolute(filepath.Dir(dir))
	}
	return filepath.EvalSymlinks(dir)
}
