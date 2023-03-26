package users

import (
	"github.com/livebud/bud/example/zero/env"
)

type Controller struct {
	Env *env.Env
}

func (c *Controller) Index() string {
	return "user index"
}

func (c *Controller) New() string {
	return "new user form"
}
