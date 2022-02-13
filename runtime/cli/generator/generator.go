package generator

import (
	"context"

	"gitlab.com/mnm/bud/pkg/buddy"

	"gitlab.com/mnm/bud/pkg/gen"
)

type Map map[string]gen.Generator

// Load the generators
func Load(kit buddy.Kit, generators Map) (*Generator, error) {
	// Add the core generator
	if err := kit.Generators(generators); err != nil {
		return nil, err
	}
	// Add the user generators
	if err := kit.Generators(generators); err != nil {
		return nil, err
	}
	return &Generator{kit}, nil
}

type Generator struct {
	kit buddy.Kit
}

func (g *Generator) Generate(ctx context.Context) error {
	return g.kit.Sync("bud/.app", "bud/.app")
}
