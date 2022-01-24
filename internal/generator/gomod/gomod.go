package gomod

import (
	_ "embed"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/mod"
)

//// go:embed gomod.gotext
// var template string

// var generator = gotemplate.MustParse("gomod", template)

// type Go struct {
// 	Version string
// }

// type Require struct {
// 	Path     string
// 	Version  string
// 	Indirect bool
// }

type Replace struct {
	Old string
	New string
}

type Generator struct {
	Module *mod.Module
}

func (g *Generator) GenerateFile(_ gen.F, file *gen.File) error {
	file.Write(g.Module.File().Format())
	return nil
}
