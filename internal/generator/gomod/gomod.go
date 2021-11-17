package gomod

import (
	_ "embed"
	"errors"
	"io/fs"
	"os"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/modcache"
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
	Modfile  *mod.File
	Go       *Go
	Requires []*Require
	Replaces []*Replace
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	code, err := os.ReadFile(g.Modfile.Directory(file.Path()))
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return g.createFile(f, file)
	}
	return g.updateFile(f, file, code)
}

func (g *Generator) updateFile(f gen.F, file *gen.File, code []byte) error {
	modfile, err := mod.Parse(modcache.Default(), file.Path(), code)
	if err != nil {
		return err
	}
	// Add any additional requires and replaces if they don't exist already
	for _, require := range g.Requires {
		if err := modfile.AddRequire(require.Path, require.Version); err != nil {
			return err
		}
	}
	for _, replace := range g.Replaces {
		if err := modfile.AddReplace(replace.Old, "", replace.New, ""); err != nil {
			return err
		}
	}
	file.Write(modfile.Format())
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
