package deploy

import (
	"fmt"

	v8 "gitlab.com/mnm/bud/js/v8"
)

type Command struct {
	VM *v8.Pool
}

type Input struct {
	AccessKey string `flag:"access-key" help:"aws access key" default:""`
	SecretKey string `flag:"secret-key" help:"aws secret key" default:""`
}

func (c *Command) Deploy(in *Input) error {
	fmt.Println(c.VM, in.AccessKey, in.SecretKey)
	return nil
}
