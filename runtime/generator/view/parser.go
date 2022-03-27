package view

import (
	"context"
	"io/fs"

	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/gomod"
	"gitlab.com/mnm/bud/package/vfs"
)

type parser struct {
	FS      fs.FS
	Module  *gomod.Module
	Imports *imports.Set
}

func (p *parser) Parse(ctx context.Context) (*State, error) {
	exist, err := vfs.SomeExist(p.FS, "view")
	if err != nil {
		return nil, err
	} else if len(exist) == 0 {
		return nil, fs.ErrNotExist
	}
	p.Imports.AddNamed("transform", p.Module.Import("bud/.app/transform"))
	p.Imports.AddNamed("overlay", "gitlab.com/mnm/bud/package/overlay")
	p.Imports.AddNamed("mod", "gitlab.com/mnm/bud/package/gomod")
	p.Imports.AddNamed("js", "gitlab.com/mnm/bud/package/js")
	p.Imports.AddNamed("view", "gitlab.com/mnm/bud/runtime/view")
	return &State{
		Imports: p.Imports.List(),
	}, nil
}
