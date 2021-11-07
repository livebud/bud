// Package mod provides an API for working with go modules.
package mod

import (
	"errors"
	"go/build"
	"path/filepath"
	"strings"
)

// dir containing the standard libraries
var stdDir = filepath.Join(build.Default.GOROOT, "src")

type Plugin struct {
	Import string
	Name   string
	Dir    string
}

// File is an interface for working with go modules.
type File interface {
	// Returns the dir of `go.mod`.
	Directory() string
	// Returns the name of the module in `go.mod`.
	ModulePath() string
	// Resolve an absolute dir to an import path.
	ResolveImport(dir string) (importPath string, err error)
	// Resolve an import path to an absolute dir.
	ResolveDirectory(importPath string) (dir string, err error)
	// Plugins returns all the bud plugins
	Plugins() ([]*Plugin, error)
}

// Directory is a convenience function for finding the directory containing
// go.mod.
func Directory(dir string) (string, error) {
	modfile, err := Find(dir)
	if err != nil {
		return "", err
	}
	return modfile.Directory(), nil
}

// Find first tries finding an explicit module file (go.mod). If no go.mod is
// found, then Find will try inferring a virtual module file from $GOPATH.
func Find(dir string) (File, error) {
	// First search for go.mod
	modfile, err := Load(dir)
	if nil == err {
		return modfile, nil
	} else if !errors.Is(err, ErrFileNotFound) {
		return nil, err
	}
	// Try inferring from the $GOPATH
	absdir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	modPath, err := Infer(absdir)
	if nil == err {
		return Virtual(modPath, absdir), nil
	} else if !errors.Is(err, ErrCantInfer) {
		return nil, err
	}
	return nil, ErrFileNotFound
}

func resolveImport(f File, dir string) (importPath string, err error) {
	abspath, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	// Handle stdlib paths
	if strings.HasPrefix(dir, stdDir) {
		return filepath.Rel(stdDir, dir)
	}
	relPath, err := filepath.Rel(f.Directory(), abspath)
	if err != nil {
		return "", err
	}
	return filepath.Join(f.ModulePath(), relPath), nil
}
