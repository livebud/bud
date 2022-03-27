package view

import (
	"context"
	_ "embed"
	"io/fs"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/gomod"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/runtime/bud"
)

//go:embed view.gotext
var template string

var generator = gotemplate.MustParse("view.gotext", template)

type Compiler struct {
	Flag   *bud.Flag
	FS     fs.FS
	Module *gomod.Module
	// DOM    *dom.Compiler
}

type State struct {
	Imports []*imports.Import
}

// Generate the view
func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func (c *Compiler) Parse(ctx context.Context) (*State, error) {
	return (&parser{
		FS:      c.FS,
		Module:  c.Module,
		Imports: imports.New(),
	}).Parse(ctx)
}

func (c *Compiler) Compile(ctx context.Context) ([]byte, error) {
	state, err := c.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return Generate(state)
}

func (c *Compiler) GenerateFile(ctx context.Context, _ overlay.F, file *overlay.File) error {
	code, err := c.Compile(ctx)
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
