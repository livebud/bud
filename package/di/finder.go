package di

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
)

var ErrNoMatch = errors.New("no match")

// Finder finds a declaration that will instantiate the data type
type Finder interface {
	Find(module *gomod.Module, dep Dependency) (Declaration, error)
}

func (i *Injector) Find(currModule *gomod.Module, dep Dependency) (Declaration, error) {
	// If modfile is nil, we default to the project modfile
	if currModule == nil {
		currModule = i.module
	}
	// Use the passed in filesystem if we're in the application module
	// Otherwise use the module's filesystem
	var fsys fs.FS = currModule
	if currModule.Directory() == i.module.Directory() {
		fsys = i.fsys
	}
	// Find the module within the filesystem
	nextModule, err := currModule.FindIn(fsys, dep.ImportPath())
	if err != nil {
		return nil, fmt.Errorf("di: unable to find module for dependency %s . %w", dep.ID(), err)
	}
	// Check again with the newly found module
	if nextModule.Directory() != currModule.Directory() {
		fsys = nextModule
	}
	// Resolve the package directory from within the module
	dir, err := nextModule.ResolveDirectoryIn(fsys, dep.ImportPath())
	if err != nil {
		return nil, fmt.Errorf("di: unable to find directory for dependency %s . %w", dep.ID(), err)
	}
	rel, err := filepath.Rel(nextModule.Directory(), dir)
	if err != nil {
		return nil, err
	}
	pkg, err := parser.New(fsys, nextModule).Parse(rel)
	if err != nil {
		return nil, err
	}
	// Look through the functions
	for _, fn := range pkg.Functions() {
		decl, err := tryFunction(fn, dep.ImportPath(), dep.TypeName())
		if err != nil {
			if err == ErrNoMatch {
				continue
			}
			return nil, err
		}
		return decl, nil
	}
	// Look through the structs
	for _, stct := range pkg.Structs() {
		decl, err := tryStruct(stct, dep.TypeName())
		if err != nil {
			if err == ErrNoMatch {
				continue
			}
			return nil, err
		}
		return decl, nil
	}
	// TODO: add breadcrumbs to help with finding the root of this error
	return nil, fmt.Errorf("di: unclear how to provide %s.", dep.ID())
}
