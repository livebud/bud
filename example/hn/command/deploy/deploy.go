package deploy

import (
	"context"
	"fmt"

	"gitlab.com/mnm/bud/package/js"
)

type Command struct {
	VM        js.VM
	AccessKey string `flag:"access-key" help:"aws access key"`
	SecretKey string `flag:"secret-key" help:"aws secret key"`
}

func (c *Command) Run(ctx context.Context) error {
	fmt.Println("deploying", c.AccessKey, c.SecretKey)
	return nil
}
