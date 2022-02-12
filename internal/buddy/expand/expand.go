package expand

import (
	"context"

	"gitlab.com/mnm/bud/pkg/buddy"
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
	// c.kit.Generator("bud/.cli/main.go", gen.FileGenerator(mainfile.New(c.genFS, c.module)))
	// c.kit.Generator(map[string]gen.Generator{
	// 	"bud/.cli/main.go":                gen.FileGenerator(mainfile.New(c.genFS, c.module)),
	// 	"bud/.cli/program/program.go":     gen.FileGenerator(program.New(c.genFS, c.injector, c.module)),
	// 	"bud/.cli/command/command.go":     gen.FileGenerator(command.New(c.genFS, c.module, c.parser)),
	// 	"bud/.cli/generator/generator.go": gen.FileGenerator(generator.New(c.genFS, c.module, c.parser)),
	// })
	// if err := dsync.Dir(c.genFS, "bud/.cli", c.module.DirFS("bud/.cli"), "."); err != nil {
	// 	return err
	// }
	// return gobin.Build(ctx, c.module.Directory(), "bud/.cli/main.go", "bud/cli")
	return nil
}
