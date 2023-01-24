package gomod

import (
	"path"
	"sort"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

type Version = module.Version
type Require = modfile.Require
type Replace = modfile.Replace

type File struct {
	file *modfile.File
}

// Import returns the module's import path (e.g. github.com/livebud/bud)
func (f *File) Import(subPaths ...string) string {
	modulePath := f.file.Module.Mod.Path
	subPath := path.Join(subPaths...)
	if modulePath == "std" {
		return subPath
	}
	return path.Join(modulePath, subPath)
}

func (f *File) AddRequire(importPath, version string) error {
	return f.file.AddRequire(importPath, version)
}

// Replace finds a replaced package within go.mod or returns nil if not found.
func (f *File) Replace(path string) *module.Version {
	for _, rep := range f.file.Replace {
		if rep.Old.Path == path {
			return &rep.Old
		}
	}
	return nil
}

func (f *File) AddReplace(oldPath, oldVers, newPath, newVers string) error {
	return f.file.AddReplace(oldPath, oldVers, newPath, newVers)
}

// Return a list of replaces
func (f *File) Replaces() (reps []*Replace) {
	reps = make([]*Replace, len(f.file.Replace))
	copy(reps, f.file.Replace)
	// Consistent ordering regardless of modfile formatting
	sort.Slice(reps, func(i, j int) bool {
		return reps[i].Old.Path < reps[j].Old.Path
	})
	return reps
}

// Return a list of requires
func (f *File) Requires() (reqs []*Require) {
	reqs = make([]*Require, len(f.file.Require))
	copy(reqs, f.file.Require)
	// Consistent ordering regardless of modfile formatting
	sort.Slice(reqs, func(i, j int) bool {
		switch {
		case reqs[i].Indirect && !reqs[j].Indirect:
			return false
		case !reqs[i].Indirect && reqs[j].Indirect:
			return true
		default:
			return reqs[i].Mod.Path < reqs[j].Mod.Path
		}
	})
	return reqs
}

// Require finds a required package within go.mod
func (f *File) Require(path string) *module.Version {
	for _, req := range f.file.Require {
		if req.Mod.Path == path {
			return &req.Mod
		}
	}
	return nil
}

func (f *File) Format() []byte {
	return modfile.Format(f.file.Syntax)
}
