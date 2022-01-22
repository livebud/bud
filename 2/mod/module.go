package mod

import (
	"fmt"
	"go/build"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/mnm/bud/2/virtual"
	"gitlab.com/mnm/bud/go/is"
)

type Module struct {
	opt  *option
	file *File
	dir  string
}

// Directory returns the module directory (e.g. /Users/$USER/...)
func (m *Module) Directory(subpaths ...string) string {
	return filepath.Join(append([]string{m.dir}, subpaths...)...)
}

// ModCache returns the module cache directory
func (m *Module) ModCache() string {
	return m.opt.modCache.Directory()
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
	dir, err := m.resolveDirectory(importPath)
	if err != nil {
		return nil, err
	}
	// If it's a local path, that means it's inside the module, just return
	// the existing module
	if !filepath.IsAbs(dir) {
		return m, nil
	}
	return find(m.opt, dir)
}

// Open a file within the module
func (m *Module) Open(name string) (fs.File, error) {
	if m.opt.fileCache == nil {
		return os.Open(filepath.Join(m.dir, name))
	}
	return m.cachedOpen(name)
}

func (m *Module) cachedOpen(name string) (fs.File, error) {
	fcache := m.opt.fileCache
	if fcache.Has(name) {
		return fcache.Open(name)
	}
	file, err := os.Open(filepath.Join(m.dir, name))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	vfile, err := virtual.From(file)
	if err != nil {
		return nil, err
	}
	fcache.Set(name, vfile)
	return fcache.Open(name)
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
		// Check if the package path exists
		absdir := filepath.Join(m.dir, rel)
		if _, err := os.Stat(absdir); err != nil {
			return "", fmt.Errorf("mod: unable to resolve directory for package path %q: %w", importPath, err)
		}
		return absdir, nil
	}
	// Handle replace
	for _, rep := range m.file.Replaces() {
		if contains(rep.Old.Path, importPath) {
			relPath := strings.TrimPrefix(importPath, rep.Old.Path)
			newPath := filepath.Join(rep.New.Path, relPath)
			absdir, err := resolvePath(m.dir, newPath)
			if err != nil {
				return "", err
			}
			// Ensure the resolved directory exists.
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
			dir, err := m.opt.modCache.ResolveDirectory(req.Mod.Path, req.Mod.Version)
			if err != nil {
				return "", err
			}
			absdir := filepath.Join(dir, relPath)
			// Ensure the resolved directory exists.
			if _, err := os.Stat(absdir); err != nil {
				return "", fmt.Errorf("mod: unable to resolve directory for required import path %q: %w", importPath, err)
			}
			return absdir, nil
		}
	}
	return "", fmt.Errorf("mod: unable to resolve directory for import path %q: %w", importPath, fs.ErrNotExist)
}

// Resolve allows `path` to be replaced by an absolute path in `rest`
func resolvePath(path string, rest ...string) (string, error) {
	result := path
	for _, p := range rest {
		if filepath.IsAbs(p) {
			result = p
			continue
		}
		result = filepath.Join(result, p)
	}
	return filepath.Abs(result)
}

func contains(basePath, importPath string) bool {
	return basePath == importPath || strings.HasPrefix(importPath, basePath+"/")
}
