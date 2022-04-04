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
	"gitlab.com/mnm/bud/runtime/transform"
)

//go:embed view.gotext
var template string

var generator = gotemplate.MustParse("view.gotext", template)

type Compiler struct {
	Flag      *bud.Flag
	Module    *gomod.Module
	Transform *transform.Map
}

// Generate the view
func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func (c *Compiler) Parse(fsys fs.FS, ctx context.Context) (*State, error) {
	return (&parser{
		Flag:      c.Flag,
		Module:    c.Module,
		Imports:   imports.New(),
		Transform: c.Transform,
	}).Parse(fsys, ctx)
}

func (c *Compiler) Compile(fsys fs.FS, ctx context.Context) ([]byte, error) {
	state, err := c.Parse(fsys, ctx)
	if err != nil {
		return nil, err
	}
	return Generate(state)
}

func (c *Compiler) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	code, err := c.Compile(fsys, ctx)
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
