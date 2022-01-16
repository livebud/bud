package modcache

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
	"golang.org/x/mod/sumdb/dirhash"
	modzip "golang.org/x/mod/zip"
)

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

type Files = map[string]string
type Modules = map[string]Files

// Write modules directly into the cache directory in an acceptable format so
// that Go thinks these files are cached and doesn't try reading them from the
// network.
//
// This implementation is the minimal format needed to get `go mod tidy` to
// think the files are cached. This shouldn't be used outside of testing
// contexts.
//
// Based on: https://github.com/golang/go/blob/master/src/cmd/go/internal/modfetch/fetch.go
func (c *Cache) Write(modules Modules) error {
	eg := new(errgroup.Group)
	for modulePathVersion, files := range modules {
		modulePathVersion, files := modulePathVersion, files
		eg.Go(func() error { return c.writeModule(modulePathVersion, files) })
	}
	return eg.Wait()
}

func (c *Cache) writeModule(modulePathVersion string, files map[string]string) error {
	moduleParts := strings.SplitN(modulePathVersion, "@", 2)
	if len(moduleParts) != 2 {
		return fmt.Errorf("modcache: invalid module key")
	}
	modulePath, moduleVersion := moduleParts[0], moduleParts[1]
	goMod, ok := files["go.mod"]
	if !ok {
		goMod = `module ` + modulePath + "\n"
		// Write go.mod back into module to make cached files a valid go.mod
		files["go.mod"] = goMod
	}
	if modulePath == "" {
		return fmt.Errorf("modcache: missing module path in go.mod")
	}
	if modfile.ModulePath([]byte(goMod)) != modulePath {
		return fmt.Errorf("modcache: %q does not match module path in go.mod", modulePath)
	}
	downloadDir, err := c.downloadDir(modulePath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return err
	}
	escapedVersion, err := module.EscapeVersion(moduleVersion)
	if err != nil {
		return err
	}
	extlessPath := filepath.Join(downloadDir, escapedVersion)
	zipPath := extlessPath + ".zip"
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	module := module.Version{Path: modulePath, Version: moduleVersion}
	var zipFiles []modzip.File
	for path, data := range files {
		zipFiles = append(zipFiles, &zipEntry{path, data})
	}
	if err := modzip.Create(zipFile, module, zipFiles); err != nil {
		return err
	}
	if err := zipFile.Close(); err != nil {
		return err
	}
	hash, err := dirhash.HashZip(zipPath, dirhash.DefaultHash)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(extlessPath+".ziphash", []byte(hash), 0644); err != nil {
		return err
	}
	if err := ioutil.WriteFile(extlessPath+".mod", []byte(goMod), 0644); err != nil {
		return err
	}
	moduleDir, err := c.getModuleDirectory(modulePath, moduleVersion)
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

type zipEntry struct {
	path, data string
}

func (z *zipEntry) Path() string {
	return z.path
}

// Lstat returns information about the file. If the file is a symbolic link,
func (z *zipEntry) Lstat() (os.FileInfo, error) {
	return &fileInfo{
		name:    filepath.Base(z.path),
		data:    []byte(z.data),
		size:    int64(len(z.data)),
		mode:    fs.FileMode(0644),
		modTime: time.Now(),
	}, nil
}

// A fileInfo implements fs.FileInfo and fs.DirEntry for a given map file.
type fileInfo struct {
	name    string
	data    []byte
	size    int64
	mode    fs.FileMode
	modTime time.Time
	sys     interface{}
}

func (i *fileInfo) Name() string               { return i.name }
func (i *fileInfo) Mode() fs.FileMode          { return i.mode }
func (i *fileInfo) Type() fs.FileMode          { return i.mode.Type() }
func (i *fileInfo) ModTime() time.Time         { return i.modTime }
func (i *fileInfo) IsDir() bool                { return i.mode&fs.ModeDir != 0 }
func (i *fileInfo) Sys() interface{}           { return i.sys }
func (i *fileInfo) Info() (fs.FileInfo, error) { return i, nil }
func (i *fileInfo) Size() int64                { return i.size }

// Open provides access to the data within a regular file. Open may return
// an error if called on a directory or symbolic link.
func (z *zipEntry) Open() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBufferString(z.data)), nil
}

// func (c *Cache) downloadDir(modulePath string) (string, error) {
// 	module.EscapePath
// 	enc, err := module.(modulePath)
// 	if err != nil {
// 		return "", err
// 	}
// 	return filepath.Join(c.cacheDir, "cache/download", enc, "/@v"), nil
// }

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
	// partialPath, err := c.partialDownloadPath(modulePath, version, "partial")
	// if err != nil {
	// 	return dir, err
	// }
	// if _, err := os.Stat(partialPath); err == nil {
	// 	return dir, &downloadDirPartialError{dir, errors.New("not completely extracted")}
	// } else if !os.IsNotExist(err) {
	// 	return dir, err
	// }
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

// partialDownloadPath returns the partial download path
func (c *Cache) downloadDir(modulePath string) (string, error) {
	enc, err := module.EscapePath(modulePath)
	if err != nil {
		return "", err
	}
	return filepath.Join(c.cacheDir, "cache/download", enc, "/@v"), nil
}
