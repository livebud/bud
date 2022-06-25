package generator

import (
	"context"
	_ "embed"
	"io/fs"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/overlay"
	goparse "github.com/livebud/bud/package/parser"
)

//go:embed generator.gotext
var template string

var generator = gotemplate.MustParse("generator.gotext", template)

func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func New(fs fs.FS, module *gomod.Module, parser *goparse.Parser) *Generator {
	return &Generator{fs, module, parser}
}

type Generator struct {
	fs     fs.FS
	module *gomod.Module
	parser *goparse.Parser
}

func (c *Generator) Parse(ctx context.Context) (*State, error) {
	return (&parser{
		fs:      c.fs,
		module:  c.module,
		parser:  c.parser,
		imports: imports.New(),
	}).Parse(ctx)
}

func (c *Generator) Compile(ctx context.Context) ([]byte, error) {
	// Parse project commands into state
	state, err := c.Parse(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: Add in the core commands or a default command

	// Generate code from the state
	return Generate(state)
}

func (c *Generator) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	code, err := c.Compile(ctx)
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
