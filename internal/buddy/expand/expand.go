package expand

import (
	"context"

	"gitlab.com/mnm/bud/internal/gobin"

	"gitlab.com/mnm/bud/internal/dsync"

	"gitlab.com/mnm/bud/generator/cli/mainfile"

	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gen"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

func New(
	genFS *gen.FileSystem,
	injector *di.Injector,
	module *gomod.Module,
	parser *parser.Parser,
) *Command {
	return &Command{genFS, injector, module, parser}
}

type Command struct {
	genFS    *gen.FileSystem
	injector *di.Injector
	module   *gomod.Module
	parser   *parser.Parser
}

type Input struct {
	Hot    bool
	Minify bool
	Embed  bool
}

func (c *Command) Expand(ctx context.Context, in *Input) error {
	c.genFS.Add(map[string]gen.Generator{
		"bud/.cli/main.go": gen.FileGenerator(mainfile.New(c.genFS, c.module)),
		// "bud/.cli/program/program.go":     gen.FileGenerator(program.New(c.genFS, c.injector, c.module)),
		// "bud/.cli/command/command.go":     gen.FileGenerator(command.New(c.genFS, c.module, c.parser)),
		// "bud/.cli/generator/generator.go": gen.FileGenerator(generator.New(c.genFS, c.module, c.parser)),
	})
	if err := dsync.Dir(c.genFS, "bud/.cli", c.module.DirFS("bud/.cli"), "."); err != nil {
		return err
	}
	return gobin.Build(ctx, c.module.Directory(), "bud/.cli/main.go", "bud/cli")
}
