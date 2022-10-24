package imhash

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/livebud/bud/internal/versions"

	"github.com/cespare/xxhash"
	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
)

func find(module *gomod.Module, mainDir string) (*fileSet, error) {
	if exists, err := exists(module, mainDir); err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf("imhash: %q %w", mainDir, fs.ErrNotExist)
	}
	fset := newFileSet()
	// Add the following if they exist
	if err := addIfExist(module, fset, "go.mod", "package.json", "package-lock.json"); err != nil {
		return nil, err
	}
	imset := newImportSet()
	if err := findDeps(fset, imset, module, mainDir); err != nil {
		return nil, err
	}
	return fset, nil
}

// Hash traverse the imports in mainDir and generates a hash. This hash will
// change if the contents of any imported packages change.
func Hash(module *gomod.Module, mainDir string) (string, error) {
	hsh, err := hash(module, mainDir)
	if err != nil {
		return "", fmt.Errorf("imhash: unable to hash %q. %w", mainDir, err)
	}
	return hsh, err
}

func hash(module *gomod.Module, mainDir string) (string, error) {
	fset, err := find(module, mainDir)
	if err != nil {
		return "", err
	}
	return fset.Hash(module)
}

func Debug(module *gomod.Module, mainDir string, w io.Writer) error {
	fset, err := find(module, mainDir)
	if err != nil {
		return err
	}
	return fset.Debug(module, w)
}

func exists(fsys fs.FS, path string) (bool, error) {
	if _, err := fs.Stat(fsys, path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func addIfExist(module *gomod.Module, fset *fileSet, paths ...string) error {
	eg := new(errgroup.Group)
	for _, path := range paths {
		path := path
		eg.Go(func() error {
			// Add go.mod if it exists
			if exists, err := exists(module, path); err != nil {
				return err
			} else if exists {
				fset.Add(path)
			}
			return nil
		})
	}
	return eg.Wait()
}

func hashFile(fsys fs.FS, filePath string) ([]byte, error) {
	f, err := fsys.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := xxhash.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func newFileSet() *fileSet {
	return &fileSet{
		m: map[string]struct{}{},
	}
}

type fileSet struct {
	mu sync.RWMutex
	m  map[string]struct{}
}

func (s *fileSet) Has(path string) bool {
	s.mu.RLock()
	_, ok := s.m[path]
	s.mu.RUnlock()
	return ok
}

func (s *fileSet) Add(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.m[path]; !ok {
		s.m[path] = struct{}{}
	}
}

func (s *fileSet) List() (list []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for path := range s.m {
		list = append(list, path)
	}
	sort.Strings(list)
	return list
}

func (s *fileSet) Hash(fsys fs.FS) (string, error) {
	h := xxhash.New()
	for _, file := range s.List() {
		hash, err := hashFile(fsys, file)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(h, "%x %s\n", hash, file)
	}
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil)), nil
}

func (s *fileSet) Debug(fsys fs.FS, w io.Writer) error {
	for _, file := range s.List() {
		hash, err := hashFile(fsys, file)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%x %s\n", hash, file)
	}
	return nil
}

func newImportSet() *importSet {
	return &importSet{
		m: map[string]struct{}{},
	}
}

type importSet struct {
	mu sync.RWMutex
	m  map[string]struct{}
}

func (i *importSet) Has(path string) bool {
	i.mu.RLock()
	_, ok := i.m[path]
	i.mu.RUnlock()
	return ok
}

func (i *importSet) Add(path string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if _, ok := i.m[path]; !ok {
		i.m[path] = struct{}{}
	}
}

func shouldWalk(module *gomod.Module, importPath string) bool {
	if module.IsLocal(importPath) {
		return true
	}
	// Search livebud if we're in development, otherwise skip it
	if versions.Bud == "latest" && strings.HasPrefix(importPath, "github.com/livebud/bud") {
		return true
	}
	return false
}

func findDeps(fset *fileSet, imset *importSet, module *gomod.Module, dir string) (err error) {
	imported, err := parser.Import(module, dir)
	if err != nil {
		return err
	}
	// Add all the Go files
	for _, path := range imported.GoFiles {
		fset.Add(filepath.Join(dir, path))
	}
	// Add all the embeds
	for _, path := range imported.EmbedPatterns {
		files, err := filepath.Glob(module.Directory(dir, path))
		if err != nil {
			return err
		}

		for _, file := range files {
			relPath, err := filepath.Rel(module.Directory(), file)
			if err != nil {
				return err
			}

			fset.Add(relPath)
		}
	}
	// Traverse imports and compute a hash
	eg := new(errgroup.Group)
	for _, importPath := range imported.Imports {
		importPath := importPath
		if imset.Has(importPath) {
			continue
		}
		imset.Add(importPath)
		if !shouldWalk(module, importPath) {
			continue
		}
		eg.Go(func() error {
			return findImport(fset, imset, module, dir, importPath)
		})
	}
	return eg.Wait()
}

func findImport(fset *fileSet, imset *importSet, module *gomod.Module, from, importPath string) error {
	dir, err := module.ResolveDirectory(importPath)
	if err != nil {
		return fmt.Errorf("imhash: error finding import %q from %q. %w", importPath, from, err)
	}
	relPath, err := filepath.Rel(module.Directory(), dir)
	if err != nil {
		return err
	}
	if err := findDeps(fset, imset, module, relPath); err != nil {
		return err
	}
	return nil
}
