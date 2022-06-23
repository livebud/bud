package command

import (
	"context"
	"io/fs"

	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/overlay"
	goparse "github.com/livebud/bud/package/parser"
)

func New(injector *di.Injector, module *gomod.Module, parser *goparse.Parser) *Command {
	return &Command{injector, module, parser}
}

type Command struct {
	injector *di.Injector
	module   *gomod.Module
	parser   *goparse.Parser
}

func (c *Command) Parse(ctx context.Context, fsys fs.FS) (*State, error) {
	return (&parser{
		fs:       fsys,
		imports:  imports.New(),
		injector: c.injector,
		module:   c.module,
		parser:   c.parser,
	}).Parse(ctx)
}

func (c *Command) Compile(ctx context.Context, fsys fs.FS) ([]byte, error) {
	// Parse project commands into state
	state, err := c.Parse(ctx, fsys)
	if err != nil {
		return nil, err
	}
	// TODO: Add in the core commands or a default command

	// Generate code from the state
	return Generate(state)
}

func (c *Command) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	code, err := c.Compile(ctx, fsys)
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
