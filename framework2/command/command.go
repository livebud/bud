package command

import (
	"context"
	"io/fs"

	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/gomod"
	goparse "gitlab.com/mnm/bud/pkg/parser"
)

func New(fs fs.FS, module *gomod.Module, parser *goparse.Parser) *Command {
	return &Command{fs, module, parser}
}

type Command struct {
	fs     fs.FS
	module *gomod.Module
	parser *goparse.Parser
}

func (c *Command) Parse(ctx context.Context) (*State, error) {
	return (&parser{
		fs:      c.fs,
		module:  c.module,
		parser:  c.parser,
		imports: imports.New(),
	}).Parse(ctx)
}

func (c *Command) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	return nil
}
