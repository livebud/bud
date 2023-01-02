package toolbs

import (
	"context"

	"github.com/livebud/bud/internal/config"
)

func New(provide config.Provide) *Command {
	return &Command{provide: provide}
}

type Command struct {
	provide config.Provide
}

func (c *Command) Run(ctx context.Context) error {
	log, err := c.provide.Logger()
	if err != nil {
		return err
	}
	server, err := c.provide.BudServer()
	if err != nil {
		return err
	}
	defer server.Close()
	log.Info("Listening on http://127.0.0.1:35729")
	return server.Wait()
}
