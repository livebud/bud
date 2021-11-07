package run

import (
	"context"
	"fmt"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/internal/generator/controller"
	"gitlab.com/mnm/bud/internal/generator/web"
)

func New(gen gen.FS) *Command {
	return &Command{
		gen:   gen,
		URL:   "localhost:3000",
		Hot:   false,
		Embed: true,
	}
}

type Command struct {
	gen   gen.FS
	URL   string
	Hot   bool
	Embed bool
}

func (c *Command) Run(ctx context.Context, generators map[string]gen.Generator) error {
	fmt.Println("running code!", c.URL, c.Hot, c.Embed)
	// 1. Run the generators
	c.gen.Add(map[string]gen.Generator{
		// "bud/command/main.go":          run.Generator(),
		"bud/controller/controller.go": gen.FileGenerator(&controller.Generator{}),
		"bud/web/web.go":               gen.FileGenerator(&web.Generator{}),
		// "bud/env/env.go":               env.Generator(),
		// "bud/public/public.go":         public.Generator(),
		// "bud/transform/transform.go":   transform.Generator(),
	})
	// Add the user-defined generators
	c.gen.Add(generators)
	// fsync.Dir(c.bfs,)
	// 2. go run bud/main.go
	// 3. Wait for changes
	return nil
}
