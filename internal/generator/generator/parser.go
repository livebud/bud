package generator

import (
	"context"
	"io/fs"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/gomod"
	goparse "github.com/livebud/bud/package/parser"
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
	p.imports.AddNamed("overlay", "github.com/livebud/bud/package/overlay")
	p.imports.AddNamed("mainfile", "github.com/livebud/bud/runtime/generator/mainfile")
	p.imports.AddNamed("program", "github.com/livebud/bud/runtime/generator/program")
	p.imports.AddNamed("command", "github.com/livebud/bud/runtime/generator/command")
	p.imports.AddNamed("web", "github.com/livebud/bud/runtime/generator/web")
	p.imports.AddNamed("public", "github.com/livebud/bud/runtime/generator/public")
	p.imports.AddNamed("controller", "github.com/livebud/bud/runtime/generator/controller")
	p.imports.AddNamed("view", "github.com/livebud/bud/runtime/generator/view")
	state = new(State)
	state.Imports = p.imports.List()
	return state, nil
}
