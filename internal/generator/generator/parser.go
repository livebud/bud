package generator

import (
	"context"
	"io/fs"

	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/gomod"
	goparse "gitlab.com/mnm/bud/package/parser"
)

type parser struct {
	bail.Struct
	fs      fs.FS
	module  *gomod.Module
	parser  *goparse.Parser
	imports *imports.Set
}

func (p *parser) Parse(ctx context.Context) (state *State, err error) {
	defer p.Recover2(&err, "generator: unable to parse")
	p.imports.AddStd("io/fs")
	p.imports.AddNamed("overlay", "gitlab.com/mnm/bud/package/overlay")
	p.imports.AddNamed("mainfile", "gitlab.com/mnm/bud/runtime/generator/mainfile")
	p.imports.AddNamed("program", "gitlab.com/mnm/bud/runtime/generator/program")
	p.imports.AddNamed("command", "gitlab.com/mnm/bud/runtime/generator/command")
	p.imports.AddNamed("web", "gitlab.com/mnm/bud/runtime/generator/web")
	p.imports.AddNamed("public", "gitlab.com/mnm/bud/runtime/generator/public")
	p.imports.AddNamed("action", "gitlab.com/mnm/bud/runtime/generator/action")
	p.imports.AddNamed("view", "gitlab.com/mnm/bud/runtime/generator/view")
	p.imports.AddNamed("transform", "gitlab.com/mnm/bud/runtime/generator/transform")
	state = new(State)
	state.Imports = p.imports.List()
	return state, nil
}
