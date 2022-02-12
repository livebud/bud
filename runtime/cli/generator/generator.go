package generator

import (
	"context"

	"gitlab.com/mnm/bud/pkg/buddy"

	"gitlab.com/mnm/bud/pkg/gen"
)

type Handler struct {
	Path      string
	Generator gen.Generator
}

func Load(kit buddy.Kit, handlers ...Handler) (*Generator, error) {
	for _, handler := range handlers {
		if err := kit.Generator(handler.Path, handler.Generator); err != nil {
			return nil, err
		}
	}
	return &Generator{kit}, nil
}

type Generator struct {
	kit buddy.Kit
}

func (g *Generator) Generate(ctx context.Context) error {
	return g.kit.Sync("bud/.app", "bud/.app")
}
