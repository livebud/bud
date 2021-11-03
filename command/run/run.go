package run

import (
	"context"
	"fmt"

	"gitlab.com/mnm/bud/bfs"
	"gitlab.com/mnm/bud/generator/command/run"
	"gitlab.com/mnm/bud/generator/controller"
	"gitlab.com/mnm/bud/generator/env"
	"gitlab.com/mnm/bud/generator/public"
	"gitlab.com/mnm/bud/generator/transform"
	"gitlab.com/mnm/bud/generator/web"
)

func New(bfs bfs.BFS) *Command {
	return &Command{
		bfs:   bfs,
		URL:   "localhost:3000",
		Hot:   false,
		Embed: true,
	}
}

type Command struct {
	bfs   bfs.BFS
	URL   string
	Hot   bool
	Embed bool
}

func (c *Command) Run(ctx context.Context, generators map[string]bfs.Generator) error {
	fmt.Println("running code!", c.URL, c.Hot, c.Embed)
	// 1. Run the generators
	c.bfs.Add(map[string]bfs.Generator{
		"bud/command/main.go":          run.Generator(),
		"bud/controller/controller.go": controller.Generator(),
		"bud/env/env.go":               env.Generator(),
		"bud/public/public.go":         public.Generator(),
		"bud/transform/transform.go":   transform.Generator(),
		"bud/web/web.go":               web.Generator(),
	})
	// Add the user-defined generators
	c.bfs.Add(generators)
	// fsync.Dir(c.bfs,)
	// 2. go run bud/main.go
	// 3. Wait for changes
	return nil
}
