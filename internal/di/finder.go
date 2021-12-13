package di

import (
	"errors"
	"fmt"
	"path/filepath"

	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/parser"
)

var ErrNoMatch = errors.New("no match")

// Finder finds a declaration that will instantiate the data type
type Finder interface {
	Find(module *mod.Module, dep Dependency) (Declaration, error)
}

func (i *Injector) Find(module *mod.Module, dep Dependency) (Declaration, error) {
	// If modfile is nil, we default to the project modfile
	if module == nil {
		module = i.module
	}
	// Check if we have a type alias
	if alias, ok := i.typeMap[dep.ID()]; ok {
		dep = alias
	}
	// Resolve the absolute directory based on the import
	dir, err := module.ResolveDirectory(dep.ImportPath())
	if err != nil {
		// This error shouldn't be wrapped because it can be an fs.ErrNotExist which
		// is ignored by fsync.Dir. If a dependency doesn't exist, di should error
		// out with it's own error type.
		return nil, fmt.Errorf("di: unable to find dependency %s > %s", dep.ID(), err)
	}
	rel, err := filepath.Rel(module.Directory(), dir)
	if err != nil {
		return nil, err
	}
	// Parse the package
	pkg, err := parser.New(module).Parse(rel)
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
	return nil, fmt.Errorf("di: unclear how to provide %s", dep.ID())
}
