package gomod

import (
	_ "embed"
	"io/fs"
	"os"

	"gitlab.com/mnm/bud/go/mod"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed gomod.gotext
var template string

var generator = gotemplate.MustParse("gomod", template)

type State struct {
	Modfile  mod.File
	Go       *Go
	Requires []*Require
	Replaces []*Replace
}

type Go struct {
	Version string
}

type Require struct {
	Path     string
	Version  string
	Indirect bool
}

type Replace struct {
	Old string
	New string
}

type Generator struct {
	Modfile  mod.File
	Go       *Go
	Requires []*Require
	Replaces []*Replace
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	// Don't create this file if it already exists
	dirfs := os.DirFS(g.Modfile.Directory())
	if _, err := fs.Stat(dirfs, file.Path()); nil == err {
		return nil
	}
	code, err := generator.Generate(g)
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}
