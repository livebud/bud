package command

import (
	_ "embed"
	"fmt"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/gomod"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("framework/app/app.gotext", template)

func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func New(module *gomod.Module, flag *framework.Flag) *Generator {
	return &Generator{flag, module}
}

type Generator struct {
	flag   *framework.Flag
	module *gomod.Module
}

func (g *Generator) GenerateFile(fsys budfs.FS, file *budfs.File) error {
	state, err := Load(fsys, g.module)
	if err != nil {
		return err
	}
	code, err := Generate(state)
	if err != nil {
		return err
	}
	fmt.Println(string(code))
	file.Data = code
	return nil
}
