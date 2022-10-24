// Package migrate migrates the database
package migrate

import (
	"context"
	"fmt"

	"github.com/livebud/bud/package/router"
)

type Command struct {
	*router.Router
}

type UpInput struct {
	Table string `flag:"table" help:"table name"`
	N     int    `arg:"n" help:"number of migrations to run up"`
}

func (c *Command) Up(ctx context.Context, in *UpInput) error {
	fmt.Println("running up!", c.Router, in.N, in.Table)
	return nil
}

type DownInput struct {
	Table string `flag:"table" help:"table name"`
	N     int    `arg:"n" help:"number of migrations to run down"`
}

func (c *Command) Down(ctx context.Context, in *DownInput) error {
	fmt.Println("running down!", c.Router, in.N, in.Table)
	return nil
}

type InfoInput struct {
	Table string `flag:"table" help:"table name"`
}

func (c *Command) Info(ctx context.Context, in *InfoInput) error {
	fmt.Println("running input!", c.Router)
	return nil
}
