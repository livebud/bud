package gomod

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/livebud/bud/internal/gois"
	"github.com/livebud/bud/package/virtual"
)

type Module struct {
	opt  *option
	file *File
	dir  string
	fsys virtual.FS
}

var _ virtual.FS = (*Module)(nil)

// Directory returns the module directory (e.g. /Users/$USER/...)
func (m *Module) Directory(subpaths ...string) string {
	return filepath.Join(append([]string{m.dir}, subpaths...)...)
}

// ModCache returns the module cache directory
func (m *Module) ModCache() string {
	return m.opt.modCache.Directory()
}

// Import returns the module's import path (e.g. github.com/livebud/bud)
func (m *Module) Import(subpaths ...string) string {
	return m.file.Import(subpaths...)
}

// Get go.mod
func (m *Module) File() *File {
	return m.file
}

// Find a dependency from an import path
func (m *Module) Find(importPath string) (*Module, error) {
	return m.FindIn(os.DirFS(m.dir), importPath)
}

// Find a dependency from an import path within fsys
// Note: go.mod itself needs to really be in the filesystem
func (m *Module) FindIn(fsys fs.FS, importPath string) (*Module, error) {
	dir, err := m.ResolveDirectoryIn(fsys, importPath)
	if err != nil {
		return nil, err
	}
	return find(m.opt, dir)
}

func (m *Module) FindBy(match func(req *Require) bool) (modules []*Module, err error) {
	for _, req := range m.file.Requires() {
		if match(req) {
			module, err := m.Find(req.Mod.Path)
			if err != nil {
				// Ignore imports that don't have a go.mod (e.g. legacy modules)
				if errors.Is(err, fs.ErrNotExist) {
					continue
				}
				return nil, err
			}
			modules = append(modules, module)
		}
	}
	return modules, nil
}

// Open a file within the module
func (m *Module) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(m.dir, name))
}

var _ fs.StatFS = (*Module)(nil)
var _ fs.ReadDirFS = (*Module)(nil)

func (m *Module) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(filepath.Join(m.dir, name))
}

func (m *Module) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(filepath.Join(m.dir, name))
}

// ResolveImport returns an import path from a local directory.
func (m *Module) ResolveImport(dir string) (importPath string, err error) {
	return m.resolveImport(dir, true)
}

func (m *Module) resolveImport(dir string, evalSymlinks bool) (string, error) {
	relPath, err := filepath.Rel(m.dir, dir)
	if err != nil {
		return "", err
	} else if strings.HasPrefix(relPath, "..") {
		if !evalSymlinks {
			return "", fmt.Errorf("module: unable to resolve import. %q can't be outside the module directory %q", dir, m.dir)
		}
		// Maybe the directory is a symlink, resolve that symlink and try again.
		if dir, err = filepath.EvalSymlinks(dir); err != nil {
			return "", fmt.Errorf("module: unable to resolve import for %q. %w", dir, err)
		}
		return m.resolveImport(dir, false)
	}
	return m.Import(relPath), nil
}

// dir containing the standard libraries
var stdDir = filepath.Join(findGoRoot(), "src")

// ResolveDirectory resolves an import to an absolute path
func (m *Module) ResolveDirectory(importPath string) (directory string, err error) {
	return m.ResolveDirectoryIn(os.DirFS(m.dir), importPath)
}

// IsLocal returns true if the import path is within the module
func (m *Module) IsLocal(importPath string) bool {
	return contains(m.Import(), importPath)
}

// ResolveDirectory resolves an import to an absolute path.
// LocalFS maybe used if we're resolving an import path from within the current
// modules filesystem.
func (m *Module) ResolveDirectoryIn(localFS fs.FS, importPath string) (directory string, err error) {
	// Handle standard library
	if gois.StdLib(importPath) {
		return filepath.Join(stdDir, importPath), nil
	}
	// Handle local packages
	modulePath := m.Import()
	if contains(modulePath, importPath) {
		// Ensure the resolved relative dir exists
		rel, err := filepath.Rel(modulePath, importPath)
		if err != nil {
			return "", err
		}
		// Check if the package path exists
		if _, err := fs.Stat(localFS, rel); err != nil {
			return "", fmt.Errorf("mod: unable to resolve directory for package path %q. %w", importPath, err)
		}
		absdir := filepath.Join(m.dir, rel)
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
				return "", fmt.Errorf("mod: unable to resolve directory for replaced import path %q. %w", importPath, err)
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
				return "", fmt.Errorf("mod: unable to resolve directory for required import path %q. %w", importPath, err)
			}
			return absdir, nil
		}
	}
	return "", fmt.Errorf("mod: unable to resolve directory for import path %q. %w", importPath, fs.ErrNotExist)
}

// Hash the module
func (m *Module) Hash() []byte {
	code := m.File().Format()
	h := xxhash.New()
	h.Write(code)
	return h.Sum(nil)
}

// MkdirAll creates dir within the module dir. Used to implement virtual.FS
func (m *Module) MkdirAll(path string, perm fs.FileMode) error {
	return m.fsys.MkdirAll(path, perm)
}

// WriteFile writes a file within the module dir. Used to implement virtual.FS.
func (m *Module) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return m.fsys.WriteFile(name, data, perm)
}

// RemoveAll removes a file within the module dir. Used to implement virtual.FS.
func (m *Module) RemoveAll(path string) error {
	return m.fsys.RemoveAll(path)
}

// RemoveAll removes a file within the module dir. Used to implement virtual.FS.
func (m *Module) Sub(dir string) (virtual.FS, error) {
	return m.fsys.Sub(dir)
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
