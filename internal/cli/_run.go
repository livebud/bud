package cli

import "context"

// runCommand creates a *runCmd in a type-safe way
func runCommand(bud *budCmd) *runCmd {
	return &runCmd{Bud: bud}
}

// runCmd command
type runCmd struct {
	Bud    *budCmd
	Embed  bool
	Hot    string
	Listen string
	Minify bool
}

// Run is triggered when calling `bud run [...args]`.
func (c *runCmd) Run(ctx context.Context) error {
	return nil
}
