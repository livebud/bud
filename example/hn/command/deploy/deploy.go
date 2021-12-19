package deploy

import (
	"context"
	"fmt"

	v8 "gitlab.com/mnm/bud/js/v8"
)

type Command struct {
	VM        *v8.Pool
	AccessKey string `flag:"access-key" help:"aws access key"`
	SecretKey string `flag:"secret-key" help:"aws secret key"`
}

func (c *Command) Run(ctx context.Context) error {
	fmt.Println(c.VM, c.AccessKey, c.SecretKey)
	return nil
}
