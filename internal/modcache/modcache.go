package modcache

import (
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

// Default loads a module cache from the default location
func Default() *Cache {
	return New(getCacheDir())
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

type Files = map[string]string
type Versions = map[string]Files

// Write modules to the cache directory in the proper format for
// the proxy to pick it up. Mostly used for testing.
func (c *Cache) Write(versions Versions) error {
	eg := new(errgroup.Group)
	for version, files := range versions {
		version, files := version, files
		eg.Go(func() error { return c.writeModule(version, files) })
	}
	return eg.Wait()
}

func (c *Cache) writeModule(version string, files map[string]string) error {
	goMod, ok := files["go.mod"]
	if !ok {
		return fmt.Errorf("modcache: missing go.mod in files map")
	}
	modulePath := modfile.ModulePath([]byte(goMod))
	if modulePath == "" {
		return fmt.Errorf("modcache: missing module path in go.mod")
	}
	moduleDir, err := c.getModuleDirectory(modulePath, version)
	if err != nil {
		return err
	}
	eg := new(errgroup.Group)
	for path, data := range files {
		path := filepath.Join(moduleDir, path)
		dir := filepath.Dir(path)
		data := data
		eg.Go(func() error {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
			return ioutil.WriteFile(path, []byte(data), 0644)
		})
	}
	return eg.Wait()
}

// ResolveDirectory returns the directory to which m should have been
// downloaded. An error will be returned if the module path or version cannot be
// escaped. An error satisfying errors.Is(err, os.ErrNotExist) will be returned
// along with the directory if the directory does not exist or if the directory
// is not completely populated.
func (c *Cache) ResolveDirectory(modulePath, version string) (string, error) {
	dir, err := c.getModuleDirectory(modulePath, version)
	if err != nil {
		return "", err
	}
	if fi, err := os.Stat(dir); os.IsNotExist(err) {
		return dir, err
	} else if err != nil {
		return dir, &downloadDirPartialError{dir, err}
	} else if !fi.IsDir() {
		return dir, &downloadDirPartialError{dir, errors.New("not a directory")}
	}
	partialPath, err := c.partialDownloadPath(modulePath, version, "partial")
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

// Cache for faster subsequent requests
var cacheDir string

// getCacheDir returns the module cache directory
func getCacheDir() string {
	if cacheDir != "" {
		return cacheDir
	}
	env := os.Getenv("GOMODCACHE")
	if env != "" {
		cacheDir = env
		return env
	}
	cacheDir = filepath.Join(build.Default.GOPATH, "pkg", "mod")
	return cacheDir
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

// partialDownloadPath returns the partial download path
func (c *Cache) partialDownloadPath(modulePath, version, suffix string) (string, error) {
	enc, err := module.EscapePath(modulePath)
	if err != nil {
		return "", err
	}
	dir := filepath.Join(c.cacheDir, "cache/download", enc, "/@v")
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
	return filepath.Join(dir, encVer+"."+suffix), nil
}
