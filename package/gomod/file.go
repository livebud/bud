package gomod

import (
	"path"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

type Version = module.Version
type Require = modfile.Require
type Replace = modfile.Replace

type File struct {
	file *modfile.File
}

// Import returns the module's import path (e.g. gitlab.com/mnm/bud)
func (f *File) Import(subpaths ...string) string {
	return path.Join(append([]string{f.file.Module.Mod.Path}, subpaths...)...)
}

func (f *File) AddRequire(importPath, version string) error {
	return f.file.AddRequire(importPath, version)
}

func (f *File) Replace(oldPath, newPath string) error {
	return f.AddReplace(oldPath, "", newPath, "")
}

func (f *File) AddReplace(oldPath, oldVers, newPath, newVers string) error {
	return f.file.AddReplace(oldPath, oldVers, newPath, newVers)
}

// Return a list of replaces
func (f *File) Replaces() (reps []*Replace) {
	reps = make([]*Replace, len(f.file.Replace))
	for i, rep := range f.file.Replace {
		reps[i] = rep
	}
	return reps
}

// Return a list of requires
func (f *File) Requires() (reqs []*Require) {
	reqs = make([]*Require, len(f.file.Require))
	for i, req := range f.file.Require {
		reqs[i] = req
	}
	return reqs
}

func (f *File) Format() []byte {
	return modfile.Format(f.file.Syntax)
}
