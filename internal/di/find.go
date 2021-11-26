package di

import (
	"errors"
	"fmt"
	"os"
)

var ErrNoMatch = errors.New("no match")

// Find a dependency using the searcher
func (i *Injector) Find(dep *Dependency) (Declaration, error) {
	searchPaths := i.Searcher(dep.Import)
	return i.find(searchPaths, []string{}, dep)
}

// Looks for the dependency by search path in order. The search paths are
// import paths to packages.
//
// TODO: consider parallelizing find. The order should still matter in the end
// but we can look for dependencies across search paths simultaneously.
func (i *Injector) find(searchPaths []string, foundPaths []string, dep *Dependency) (Declaration, error) {
	// No more search paths, we're unable to find this dependency
	if len(searchPaths) == 0 {
		if len(foundPaths) == 0 {
			return nil, fmt.Errorf("di: unable to find dependency %q.%s", dep.Import, dep.Type)
		} else {
			return nil, fmt.Errorf("di: unclear how to provide %q.%s", dep.Import, dep.Type)
		}
	}
	// Resolve the absolute directory based on the import
	dir, err := dep.ModFile.ResolveDirectory(searchPaths[0])
	if err != nil {
		// If the directory doesn't exist, search the next search path
		if errors.Is(err, os.ErrNotExist) {
			return i.find(searchPaths[1:], foundPaths, dep)
		}
		return nil, err
	}
	// Parse the package
	ast, err := i.Parser.Parse(dir)
	if err != nil {
		return nil, err
	}
	// Look through the functions
	for _, fn := range ast.Functions() {
		decl, err := tryFunction(fn, dep)
		if err != nil {
			if err == ErrNoMatch {
				continue
			}
			return nil, err
		}
		return decl, nil
	}
	// Look through the structs
	for _, stct := range ast.Structs() {
		decl, err := tryStruct(stct, dep)
		if err != nil {
			if err == ErrNoMatch {
				continue
			}
			return nil, err
		}
		return decl, nil
	}
	// Add to foundPaths and try again
	foundPaths = append(foundPaths, dir)
	// Search the next search path
	decl, err := i.find(searchPaths[1:], foundPaths, dep)
	if err != nil {
		return nil, err
	}
	return decl, nil
}
