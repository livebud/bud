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
	}
	return finder.findModFile(absdir)
}

// Open a file within the module
func (m *Module) Open(name string) (fs.File, error) {
	return os.DirFS(m.dir).Open(name)
}

// ResolveDirectory resolves an import to an absolute path
func (m *Module) ResolveDirectory(importPath string) (directory string, err error) {
	absdir, err := m.resolveDirectory(importPath)
	if err != nil {
		return "", err
	}
	// Ensure the resolved directory exists
	if _, err := os.Stat(absdir); err != nil {
		return "", fmt.Errorf("mod: unable to resolve directory for import path %q: %w", importPath, err)
	}
	return absdir, nil
}

// ResolveImport returns an import path from a local directory.
func (m *Module) ResolveImport(directory string) (importPath string, err error) {
	if !strings.HasPrefix(directory, m.dir) {
		return "", fmt.Errorf("%q can't be outside the module directory %q", directory, m.dir)
	}
	return m.resolveImport(directory)
}

// dir containing the standard libraries
var stdDir = filepath.Join(build.Default.GOROOT, "src")

func (m *Module) resolveImport(dir string) (importPath string, err error) {
	abspath, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	// Handle stdlib paths
	if strings.HasPrefix(dir, stdDir) {
		return filepath.Rel(stdDir, dir)
	}
	relPath, err := filepath.Rel(m.dir, abspath)
	if err != nil {
		return "", err
	}
	return m.Import(relPath), nil
}

// ResolveDirectory resolves an import to an absolute path
func (m *Module) resolveDirectory(importPath string) (directory string, err error) {
	if is.StdLib(importPath) {
		return filepath.Join(stdDir, importPath), nil
	}
	modulePath := m.Import()
	if contains(modulePath, importPath) {
		directory = filepath.Join(m.dir, strings.TrimPrefix(importPath, modulePath))
		return directory, nil
	}
	// loop over replaces
	for _, rep := range m.file.Replaces() {
		if contains(rep.Old.Path, importPath) {
			relPath := strings.TrimPrefix(importPath, rep.Old.Path)
			newPath := filepath.Join(rep.New.Path, relPath)
			resolved := resolvePath(m.dir, newPath)
			return resolved, nil
		}
	}
	// loop over requires
	for _, req := range m.file.Requires() {
		if contains(req.Mod.Path, importPath) {
			relPath := strings.TrimPrefix(importPath, req.Mod.Path)
			dir, err := m.cache.ResolveDirectory(req.Mod.Path, req.Mod.Version)
			if err != nil {
				return "", err
			}
			return filepath.Join(dir, relPath), nil
		}
	}
	return "", fmt.Errorf("mod: unable to resolve directory for import path %q", importPath)
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
