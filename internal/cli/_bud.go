package cli

import "context"

// budCommand creates a *budCmd in a type-safe way.
func budCommand() *budCmd {
	return &budCmd{}
}

// budCmd command
type budCmd struct {
	Dir  string
	Args []string
}

// Run is triggered when calling `bud [...args]`
func (c *budCmd) Run(ctx context.Context) error {
	return nil
}
