package imhash

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/cespare/xxhash"
	"gitlab.com/mnm/bud/package/parser"

	"gitlab.com/mnm/bud/package/gomod"
)

func Hash(module *gomod.Module, mainDir string) (string, error) {
	parser := parser.New(module, module)
	fset := newFileSet()
	if err := findDeps(fset, module, parser, mainDir); err != nil {
		return "", err
	}
	return fset.Hash(module)
}

func Debug(module *gomod.Module, mainDir string, w io.Writer) error {
	parser := parser.New(module, module)
	fset := newFileSet()
	if err := findDeps(fset, module, parser, mainDir); err != nil {
		return err
	}
	return fset.Debug(module, w)
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
	m map[string]struct{}
	d []string
}

func (s *fileSet) Add(path string) {
	if _, ok := s.m[path]; !ok {
		s.m[path] = struct{}{}
		s.d = append(s.d, path)
	}
}

func (s *fileSet) List() []string {
	return s.d
}

func (s *fileSet) Hash(fsys fs.FS) (string, error) {
	h := xxhash.New()
	for _, file := range s.d {
		hash, err := hashFile(fsys, file)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(h, "%x %s\n", hash, file)
	}
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil)), nil
}

func (s *fileSet) Debug(fsys fs.FS, w io.Writer) error {
	for _, file := range s.d {
		hash, err := hashFile(fsys, file)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%x %s\n", hash, file)
	}
	return nil
}

func shouldWalk(module *gomod.Module, importPath string) bool {
	return module.IsLocal(importPath) ||
		// TODO: consider removing once the API settles a bit
		strings.HasPrefix(importPath, "gitlab.com/mnm/bud/runtime")
}

func findDeps(fset *fileSet, module *gomod.Module, parser *parser.Parser, dir string) (err error) {
	imported, err := parser.Import(dir)
	if err != nil {
		return err
	}
	// Add all the Go files
	for _, path := range imported.GoFiles {
		fset.Add(filepath.Join(dir, path))
	}
	// Add all the embeds
	// TODO: resolve patterns
	for _, path := range imported.EmbedPatterns {
		fset.Add(filepath.Join(dir, path))
	}
	// Traverse imports and compute a hash
	for _, importPath := range imported.Imports {
		if !shouldWalk(module, importPath) {
			continue
		}
		dir, err := module.ResolveDirectory(importPath)
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(module.Directory(), dir)
		if err != nil {
			return err
		}
		if err := findDeps(fset, module, parser, relPath); err != nil {
			return err
		}
	}
	return nil
}
