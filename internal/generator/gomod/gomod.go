package gomod

import (
	_ "embed"
	"errors"
	"io/fs"
	"os"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed gomod.gotext
var template string

var generator = gotemplate.MustParse("gomod", template)

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
	Module   *mod.Module
	Go       *Go
	Requires []*Require
	Replaces []*Replace
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	code, err := os.ReadFile(g.Module.Directory(file.Path()))
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return g.createFile(f, file)
	}
	return g.updateFile(f, file, code)
}

func (g *Generator) updateFile(f gen.F, file *gen.File, code []byte) error {
	modFinder := mod.New()
	module, err := modFinder.Parse(file.Path(), code)
	if err != nil {
		return err
	}
	modFile := module.File()
	// Add any additional requires and replaces if they don't exist already
	for _, require := range g.Requires {
		if err := modFile.AddRequire(require.Path, require.Version); err != nil {
			return err
		}
	}
	for _, replace := range g.Replaces {
		if err := modFile.AddReplace(replace.Old, "", replace.New, ""); err != nil {
			return err
		}
	}
	file.Write(modFile.Format())
	return nil
}

func (g *Generator) createFile(f gen.F, file *gen.File) error {
	code, err := generator.Generate(g)
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}