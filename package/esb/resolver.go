package esb

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"
)

// type Resolver interface {
// 	Resolve(fsys fs.FS, from, to string) (string, error)
// }

type Resolver struct {
	Conditions []string
	MainFields []string
	Extensions []string
}

func (r *Resolver) Resolve(fsys fs.FS, from, to string) (string, error) {
	if !isPackagePath(to) {
		fpath := path.Clean(to)
		// Resolve paths relative to the importer
		if from != "" {
			fpath = path.Join(path.Dir(from), fpath)
		}
		return resolveRelative(fsys, fpath)
	}
	fmt.Println("resolving", to)
	return "", fmt.Errorf("Not implemented yet")
}

func resolveRelative(fsys fs.FS, fpath string) (string, error) {
	stat, err := fs.Stat(fsys, fpath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("unable to stat %q. %w", fpath, err)
		}
		// Read the directory, it might be an extension-less import.
		// e.g. "react-dom/server" where "react-dom/server.js" exists
		dir := path.Dir(fpath)
		des, err := fs.ReadDir(fsys, dir)
		if err != nil {
			return "", fmt.Errorf("unable to read dir %q. %w", dir, err)
		}
		baseAndDot := path.Base(fpath) + "."
		for _, de := range des {
			if de.IsDir() {
				continue
			}
			name := de.Name()
			if !strings.HasPrefix(name, baseAndDot) {
				continue
			}
			switch path.Ext(name) {
			case ".js", ".mjs", ".cjs", ".jsx", ".ts", ".tsx", ".json":
				fpath = path.Join(dir, name)
				return fpath, nil
			}
		}
		return "", fmt.Errorf("unable to resolve %q. %w", fpath, fs.ErrNotExist)
	}
	// Handle reading the index file from a main directory
	// e.g. "react" resolving to "react/index.js"
	if stat.IsDir() {
		des, err := fs.ReadDir(fsys, fpath)
		if err != nil {
			return "", fmt.Errorf("unable to read dir %q. %w", fpath, err)
		}
		for _, de := range des {
			if de.IsDir() {
				continue
			}
			name := de.Name()
			if strings.HasPrefix(name, "index.") {
				switch path.Ext(name) {
				case ".js", ".mjs", ".cjs", ".jsx", ".ts", ".tsx", ".json":
					fpath = path.Join(fpath, name)
					return fpath, nil
				}
			}
		}
		return "", fmt.Errorf("unable to resolve %q. %w", fpath, fs.ErrNotExist)
	}
	return fpath, nil
}

// Package paths are loaded from a "node_modules" directory. Non-package paths
// are relative or absolute paths.
func isPackagePath(path string) bool {
	return !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "./") &&
		!strings.HasPrefix(path, "../") && path != "." && path != ".."
}
