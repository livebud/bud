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

	"github.com/go-duo/bud/go/is"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

// ErrFileNotFound occurs when no go.mod can be found
var ErrFileNotFound = fmt.Errorf("unable to find go.mod: %w", fs.ErrNotExist)

// Load a modfile or fail trying
func Load(path string) (*Module, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return findModFile(abs)
}

// Find the go.mod file from anywhere in your project. If it's unable to find a
// go.mod file, it will also try inferring the module name by it's $GOPATH. This
// will only work if your project is inside $GOPATH.
func findModFile(path string) (*Module, error) {
	moduleDir, err := findModPath(path)
	if err != nil {
		return nil, fmt.Errorf("%w in %q", ErrFileNotFound, path)
	}
	modulePath := filepath.Join(moduleDir, "go.mod")
	moduleData, err := ioutil.ReadFile(modulePath)
	if err != nil {
		return nil, err
	}
	return Parse(modulePath, moduleData)
}

// Parse a modfile from it's data
func Parse(path string, data []byte) (*Module, error) {
	modfile, err := modfile.Parse(path, data, nil)
	if err != nil {
		return nil, err
	}
	return &Module{
		file: modfile,
		dir:  filepath.Dir(path),
	}, nil
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

// https://golang.org/src/cmd/go/internal/modload/init.go
//
// func findModuleRoot(dir string) (root string) {
// 	if dir == "" {
// 		panic("dir not set")
// 	}
// 	dir = filepath.Clean(dir)

// 	// Look for enclosing go.mod.
// 	for {
// 		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
// 			return dir
// 		}
// 		d := filepath.Dir(dir)
// 		if d == dir {
// 			break
// 		}
// 		dir = d
// 	}
// 	return ""
// }

// Module struct
type Module struct {
	file *modfile.File
	dir  string
}

// Directory returns the module directory
// e.g. /Users/$USER/...
func (m *Module) Directory() string {
	return m.dir
}

// ModulePath returns the module path
// e.g. github.com/go-duo/duoc
func (m *Module) ModulePath() string {
	return m.file.Module.Mod.Path
}

func (m *Module) ResolveDirectory(importPath string) (directory string, err error) {
	dir, err := m.resolveDirectory(importPath)
	if err != nil {
		return "", err
	}
	// Ensure the resolved directory exists
	if _, err := os.Stat(dir); err != nil {
		return "", fmt.Errorf("unable to find directory for import path %q: %w", importPath, err)
	}
	return dir, nil
}

// ResolveDirectory resolves an import to an absolute path
func (m *Module) resolveDirectory(importPath string) (directory string, err error) {
	if is.StdLib(importPath) {
		return filepath.Join(stdDir, importPath), nil
	}
	if strings.HasPrefix(importPath, m.file.Module.Mod.Path) {
		directory = filepath.Join(m.dir, strings.TrimPrefix(importPath, m.file.Module.Mod.Path))
		return directory, nil
	}
	// loop over replaces
	for _, mod := range m.file.Replace {
		if strings.HasPrefix(importPath, mod.Old.Path) {
			relPath := strings.TrimPrefix(importPath, mod.Old.Path)
			newPath := filepath.Join(mod.New.Path, relPath)
			resolved := resolvePath(m.dir, newPath)
			return resolved, nil
		}
	}
	// loop over requires
	for _, mod := range m.file.Require {
		if strings.HasPrefix(importPath, mod.Mod.Path) {
			relPath := strings.TrimPrefix(importPath, mod.Mod.Path)
			dir, err := downloadDir(mod.Mod)
			if err != nil {
				return "", err
			}
			return filepath.Join(dir, relPath), nil
		}
	}
	return "", fmt.Errorf("unable to find directory for import path %q", importPath)
}

// ResolveImport returns an import path from a directory
func (m *Module) ResolveImport(directory string) (importPath string, err error) {
	if !strings.HasPrefix(directory, m.dir) {
		return "", fmt.Errorf("%q can't be outside the module directory %q", directory, m.dir)
	}
	importPath, err = resolveImport(m, directory)
	if err != nil {
		return "", err
	}
	return importPath, nil
}

func resolvePath(path string, rest ...string) (result string) {
	result = path
	for _, p := range rest {
		if filepath.IsAbs(p) {
			result = p
			continue
		}
		result = filepath.Join(result, p)
	}
	return result
}

// Version will return the version or an empty string if not found
func (m *Module) Version(importPath string) string {
	for _, require := range m.file.Require {
		if require.Mod.Path == importPath {
			return require.Mod.Version
		}
	}
	return ""
}

// Replacement will return the replaced value or an empty string if not found
func (m *Module) Replacement(importPath string) string {
	for _, replace := range m.file.Replace {
		if replace.Old.Path == importPath {
			newPath := replace.New.Path
			if newPath[0] != '.' {
				return newPath
			}
			return filepath.Join(m.Directory(), newPath)
		}
	}
	return ""
}

func (m *Module) AddRequire(importPath, version string) (err error) {
	return m.file.AddRequire(importPath, version)
}

// GOMODCACHE returns the cache directory
var GOMODCACHE = func() string {
	env := os.Getenv("GOMODCACHE")
	if env != "" {
		return env
	}
	return filepath.Join(build.Default.GOPATH, "pkg", "mod")
}()

// downloadDir returns the directory to which m should have been downloaded.
// An error will be returned if the module path or version cannot be escaped.
// An error satisfying errors.Is(err, os.ErrNotExist) will be returned
// along with the directory if the directory does not exist or if the directory
// is not completely populated.
func downloadDir(m module.Version) (string, error) {
	if GOMODCACHE == "" {
		// modload.Init exits if GOPATH[0] is empty, and GOMODCACHE
		// is set to GOPATH[0]/pkg/mod if GOMODCACHE is empty, so this should never happen.
		return "", fmt.Errorf("internal error: GOMODCACHE not set")
	}
	enc, err := module.EscapePath(m.Path)
	if err != nil {
		return "", err
	}
	if !semver.IsValid(m.Version) {
		return "", fmt.Errorf("non-semver module version %q", m.Version)
	}
	if module.CanonicalVersion(m.Version) != m.Version {
		return "", fmt.Errorf("non-canonical module version %q", m.Version)
	}
	encVer, err := module.EscapeVersion(m.Version)
	if err != nil {
		return "", err
	}
	dir := filepath.Join(GOMODCACHE, enc+"@"+encVer)
	if fi, err := os.Stat(dir); os.IsNotExist(err) {
		return dir, err
	} else if err != nil {
		return dir, &downloadDirPartialError{dir, err}
	} else if !fi.IsDir() {
		return dir, &downloadDirPartialError{dir, errors.New("not a directory")}
	}
	partialPath, err := cachePath(m, "partial")
	if err != nil {
		return dir, err
	}
	if _, err := os.Stat(partialPath); err == nil {
		return dir, &downloadDirPartialError{dir, errors.New("not completely extracted")}
	} else if !os.IsNotExist(err) {
		return dir, err
	}
	return dir, nil
}

// downloadDirPartialError is returned by DownloadDir if a module directory
// exists but was not completely populated.
//
// downloadDirPartialError is equivalent to os.ErrNotExist.
type downloadDirPartialError struct {
	Dir string
	Err error
}

// Error fn
func (e *downloadDirPartialError) Error() string { return fmt.Sprintf("%s: %v", e.Dir, e.Err) }

// Is fn
func (e *downloadDirPartialError) Is(err error) bool { return err == os.ErrNotExist }

// cachePath returns the cache path
func cachePath(m module.Version, suffix string) (string, error) {
	dir, err := cacheDir(m.Path)
	if err != nil {
		return "", err
	}
	if !semver.IsValid(m.Version) {
		return "", fmt.Errorf("non-semver module version %q", m.Version)
	}
	if module.CanonicalVersion(m.Version) != m.Version {
		return "", fmt.Errorf("non-canonical module version %q", m.Version)
	}
	encVer, err := module.EscapeVersion(m.Version)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, encVer+"."+suffix), nil
}

func cacheDir(path string) (string, error) {
	if GOMODCACHE == "" {
		// modload.Init exits if GOPATH[0] is empty, and GOMODCACHE
		// is set to GOPATH[0]/pkg/mod if GOMODCACHE is empty, so this should never happen.
		return "", fmt.Errorf("internal error: GOMODCACHE not set")
	}
	enc, err := module.EscapePath(path)
	if err != nil {
		return "", err
	}
	return filepath.Join(GOMODCACHE, "cache/download", enc, "/@v"), nil
}
