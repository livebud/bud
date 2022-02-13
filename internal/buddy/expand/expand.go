package expand

import (
	"context"

	"gitlab.com/mnm/bud/generator/cli/command"
	"gitlab.com/mnm/bud/generator/cli/generator"
	"gitlab.com/mnm/bud/generator/cli/mainfile"
	"gitlab.com/mnm/bud/generator/cli/program"
	"gitlab.com/mnm/bud/pkg/buddy"
	"gitlab.com/mnm/bud/pkg/gen"
)

func New(kit buddy.Kit) *Command {
	return &Command{kit}
}

type Command struct {
	kit buddy.Kit
}

type Input struct {
	Hot    bool
	Minify bool
	Embed  bool
}

func (c *Command) Expand(ctx context.Context, in *Input) error {
	c.kit.Generators(map[string]buddy.Generator{
		"bud/.cli/main.go":                gen.FileGenerator(mainfile.New(c.kit)),
		"bud/.cli/program/program.go":     gen.FileGenerator(program.New(c.kit)),
		"bud/.cli/command/command.go":     gen.FileGenerator(command.New(c.kit)),
		"bud/.cli/generator/generator.go": gen.FileGenerator(generator.New(c.kit)),
	})
	if err := c.kit.Sync("bud/.cli", "bud/.cli"); err != nil {
		return err
	}
	return c.kit.Go().Build(ctx, "bud/.cli/main.go", "bud/cli")
}
