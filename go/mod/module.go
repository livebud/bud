package mod

import (
	"fmt"
	"go/build"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/mnm/bud/go/is"
	"gitlab.com/mnm/bud/internal/modcache"
)

type Module struct {
	cache *modcache.Cache
	file  *File
	dir   string
	fsys  fs.FS
}

// Directory returns the module directory (e.g. /Users/$USER/...)
func (m *Module) Directory(subpaths ...string) string {
	return filepath.Join(append([]string{m.dir}, subpaths...)...)
}

// Import returns the module's import path (e.g. gitlab.com/mnm/bud)
func (m *Module) Import(subpaths ...string) string {
	return m.file.Import(subpaths...)
}

// Get go.mod
func (m *Module) File() *File {
	return m.file
}

// Find a dependency from an import path
func (m *Module) Find(importPath string) (*Module, error) {
	absdir, err := m.resolveDirectory(importPath)
	if err != nil {
		return nil, err
	}
	finder := &Finder{
		cache: m.cache,
		fsys:  osfs{},
	}
	module, err := finder.findModFile(absdir)
	if err != nil {
		return nil, err
	}
	return module, nil
}

// Open a file within the module
func (m *Module) Open(name string) (fs.File, error) {
	return m.fsys.Open(name)
}

// ResolveDirectory resolves an import to an absolute path
func (m *Module) ResolveDirectory(importPath string) (directory string, err error) {
	return m.resolveDirectory(importPath)
}

// ResolveImport returns an import path from a local directory.
func (m *Module) ResolveImport(directory string) (importPath string, err error) {
	relPath, err := filepath.Rel(m.dir, filepath.Clean(directory))
	if err != nil {
		return "", err
	} else if strings.HasPrefix(relPath, "..") {
		return "", fmt.Errorf("%q can't be outside the module directory %q", directory, m.dir)
	}
	return m.Import(relPath), nil
}

// dir containing the standard libraries
var stdDir = filepath.Join(build.Default.GOROOT, "src")

// ResolveDirectory resolves an import to an absolute path
func (m *Module) resolveDirectory(importPath string) (directory string, err error) {
	// Handle standard library
	if is.StdLib(importPath) {
		return filepath.Join(stdDir, importPath), nil
	}
	// Handle local packages within fsys
	modulePath := m.Import()
	if contains(modulePath, importPath) {
		// Ensure the resolved relative dir exists
		rel, err := filepath.Rel(modulePath, importPath)
		if err != nil {
			return "", err
		}
		if _, err := fs.Stat(m.fsys, rel); err != nil {
			return "", fmt.Errorf("mod: unable to resolve directory for package path %q: %w", importPath, err)
		}
		// But return the absolute dir
		absdir := filepath.Join(m.dir, rel)
		return absdir, nil
	}
	// Handle replace
	for _, rep := range m.file.Replaces() {
		if contains(rep.Old.Path, importPath) {
			relPath := strings.TrimPrefix(importPath, rep.Old.Path)
			newPath := filepath.Join(rep.New.Path, relPath)
			absdir := resolvePath(m.dir, newPath)
			// Ensure the resolved directory exists. Use os because we're outside of
			// outside of fsys.
			if _, err := os.Stat(absdir); err != nil {
				return "", fmt.Errorf("mod: unable to resolve directory for replaced import path %q: %w", importPath, err)
			}
			return absdir, nil
		}
	}
	// Handle require
	for _, req := range m.file.Requires() {
		if contains(req.Mod.Path, importPath) {
			relPath := strings.TrimPrefix(importPath, req.Mod.Path)
			dir, err := m.cache.ResolveDirectory(req.Mod.Path, req.Mod.Version)
			if err != nil {
				return "", err
			}
			absdir := filepath.Join(dir, relPath)
			// Ensure the resolved directory exists. Use os because we're outside of
			// outside of fsys.
			if _, err := os.Stat(absdir); err != nil {
				return "", fmt.Errorf("mod: unable to resolve directory for required import path %q: %w", importPath, err)
			}
			return absdir, nil
		}
	}
	return "", fmt.Errorf("mod: unable to resolve directory for import path %q: %w", importPath, fs.ErrNotExist)
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

func contains(basePath, importPath string) bool {
	return basePath == importPath || strings.HasPrefix(importPath, basePath+"/")
}
