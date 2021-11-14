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

// Cache for faster subsequent requests
var cacheDir string

// Directory returns the module cache directory
func Directory() string {
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

// WriteModule writes a module to the cache directory in the proper format for
// the proxy to pick it up
func WriteModule(cacheDir, version string, files map[string][]byte) error {
	goMod, ok := files["go.mod"]
	if !ok {
		return fmt.Errorf("modcache: missing go.mod in files map")
	}
	modulePath := modfile.ModulePath(goMod)
	if modulePath == "" {
		return fmt.Errorf("modcache: missing module path in go.mod")
	}
	cacheDir, err := getCacheDir(cacheDir)
	if err != nil {
		return err
	}
	moduleDir, err := getModuleDir(cacheDir, modulePath, version)
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
			return ioutil.WriteFile(path, data, 0644)
		})
	}
	return eg.Wait()
}

// ResolveDirectory returns the directory to which m should have been
// downloaded. An error will be returned if the module path or version cannot be
// escaped. An error satisfying errors.Is(err, os.ErrNotExist) will be returned
// along with the directory if the directory does not exist or if the directory
// is not completely populated.
func ResolveDirectory(cacheDir, modulePath, version string) (string, error) {
	cacheDir, err := getCacheDir(cacheDir)
	if err != nil {
		return "", err
	}
	dir, err := getModuleDir(cacheDir, modulePath, version)
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
	partialPath, err := partialDownloadPath(cacheDir, modulePath, version, "partial")
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

func getModuleDir(cacheDir, modulePath, version string) (string, error) {
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
	dir := filepath.Join(cacheDir, enc+"@"+encVer)
	return dir, nil
}

func getCacheDir(cacheDir string) (string, error) {
	if cacheDir != "" {
		return cacheDir, nil
	}
	cacheDir = Directory()
	if cacheDir != "" {
		return cacheDir, nil
	}
	return "", fmt.Errorf("internal error: GOMODCACHE not set")
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
func partialDownloadPath(cacheDir, modulePath, version, suffix string) (string, error) {
	enc, err := module.EscapePath(modulePath)
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cacheDir, "cache/download", enc, "/@v")
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
