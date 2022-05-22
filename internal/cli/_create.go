package cli

import "context"

func createCommand(bud *budCmd) *createCmd {
	return &createCmd{Bud: bud}
}

type createCmd struct {
	Bud *budCmd
	Dir string
}

func (c *createCmd) Run(ctx context.Context) error {
	return nil
}
