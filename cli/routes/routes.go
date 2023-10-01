package routes

import (
	"context"
	"fmt"

	"github.com/livebud/bud/mux"
)

func New() *Command {
	return &Command{}
}

type Command struct {
}

func (c *Command) Run(ctx context.Context, router *mux.Router) error {
	for _, route := range router.List() {
		fmt.Println(route.String())
	}
	return nil
}
