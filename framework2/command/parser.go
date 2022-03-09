package command

import (
	"context"
	"io/fs"

	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/pkg/gomod"
	goparse "gitlab.com/mnm/bud/pkg/parser"
)

func New(fs fs.FS, module *gomod.Module, parser *goparse.Parser) *Parser {
	return &Parser{fs, module, parser}
}

type Parser struct {
	fs     fs.FS
	module *gomod.Module
	parser *goparse.Parser
}

func (p *Parser) Parse(ctx context.Context) (*State, error) {
	return (&parser{
		fs:      p.fs,
		module:  p.module,
		parser:  p.parser,
		imports: imports.New(),
	}).Parse(ctx)
}

type parser struct {
	bail.Struct
	fs      fs.FS
	module  *gomod.Module
	parser  *goparse.Parser
	imports *imports.Set
}

func (p *parser) Parse(ctx context.Context) (state *State, err error) {
	defer p.Recover(&err)
	state = new(State)
	return state, nil
}
