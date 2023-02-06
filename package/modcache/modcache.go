package modcache

import (
	"errors"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

// ErrModuleNotFound occurs when a go.mod module hasn't been downloaded yet
var ErrModuleNotFound = errors.New("module not found")

// Default loads a module cache from the default location
func Default() *Cache {
	return New(getModDir())
}

// New module cache relative to the cache directory
func New(cacheDir string) *Cache {
	return &Cache{cacheDir}
}

type Cache struct {
	cacheDir string
}

// Directory returns the cache directory joined with optional subpaths
func (c *Cache) Directory(subpaths ...string) string {
	return filepath.Join(append([]string{c.cacheDir}, subpaths...)...)
}

// SplitPathVersion splits a path@version into path & version
func SplitPathVersion(pathVersion string) (path, version string, err error) {
	parts := strings.SplitN(pathVersion, "@", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("modcache: invalid module key for %q", pathVersion)
	}
	path, version = parts[0], parts[1]
	if path == "" {
		return "", "", fmt.Errorf("modcache: missing module path in %q", pathVersion)
	}
	if version == "" {
		return "", "", fmt.Errorf("modcache: missing module version in %q", pathVersion)
	}
	return path, version, nil
}

// ResolveDirectory returns the directory to which the module should have been
// downloaded. An error will be returned if the module path or version cannot be
// escaped. An error satisfying errors.Is(err, ErrModuleNotFound) will be
// returned along with the directory if the directory does not exist or if the
// directory is not completely populated.
func (c *Cache) ResolveDirectory(modulePath, version string) (string, error) {
	dir, err := c.getModuleDirectory(modulePath, version)
	if err != nil {
		return "", err
	}
	if fi, err := os.Stat(dir); os.IsNotExist(err) {
		return "", fmt.Errorf("modcache: %s@%s %w", modulePath, version, ErrModuleNotFound)
	} else if err != nil {
		return "", &downloadDirPartialError{dir, err}
	} else if !fi.IsDir() {
		return "", &downloadDirPartialError{dir, errors.New("not a directory")}
	}
	return dir, nil
}

// Cache for faster subsequent requests
var modDir string

// getModDir returns the module cache directory
func getModDir() string {
	if modDir != "" {
		return modDir
	}
	env := os.Getenv("GOMODCACHE")
	if env != "" {
		modDir = env
		return env
	}
	modDir = filepath.Join(build.Default.GOPATH, "pkg", "mod")
	return modDir
}

// getModuleDirectory returns an absolute path to the required module.
func (c *Cache) getModuleDirectory(modulePath, version string) (string, error) {
	enc, err := module.EscapePath(modulePath)
	if err != nil {
		return "", err
	}
	if !semver.IsValid(version) {
		return "", fmt.Errorf("non-semver module version %q", version)
	}
	if module.CanonicalVersion(version) != version {
		return "", fmt.Errorf("non-canonical module version %q", version)
	}
	encVer, err := module.EscapeVersion(version)
	if err != nil {
		return "", err
	}
	dir := filepath.Join(c.cacheDir, enc+"@"+encVer)
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
