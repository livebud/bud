package cli

import "context"

// buildCommand creates a *buildCmd in a type-safe way
func buildCommand(bud *budCmd) *buildCmd {
	return &buildCmd{Bud: bud}
}

type buildCmd struct {
	Bud    *budCmd
	Embed  bool
	Hot    string
	Minify bool
}

func (c *buildCmd) Run(ctx context.Context) error {
	return nil
}
