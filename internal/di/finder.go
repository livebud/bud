package di

import (
	"errors"
	"fmt"

	"gitlab.com/mnm/bud/go/mod"
)

var ErrNoMatch = errors.New("no match")

// Finder finds a declaration that will instantiate the data type
type Finder interface {
	Find(modFile *mod.File, importPath, dataType string) (Declaration, error)
}

func (i *Injector) Find(modFile *mod.File, importPath, dataType string) (Declaration, error) {
	// If modfile is nil, we default to the project modfile
	if modFile == nil {
		modFile = i.modFile
	}
	// Resolve the absolute directory based on the import
	dir, err := modFile.ResolveDirectory(importPath)
	if err != nil {
		return nil, fmt.Errorf("di: unable to find dependency %q.%s: %w", importPath, dataType, err)
	}
	// Parse the package
	pkg, err := i.parser.Parse(dir)
	if err != nil {
		return nil, err
	}
	// Look through the functions
	for _, fn := range pkg.Functions() {
		decl, err := tryFunction(fn, importPath, dataType)
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
		decl, err := tryStruct(stct, dataType)
		if err != nil {
			if err == ErrNoMatch {
				continue
			}
			return nil, err
		}
		return decl, nil
	}
	return nil, fmt.Errorf("di: unclear how to provide %q.%s", importPath, dataType)
}
