// Package enables or disables web dyno autoscaling
package autoscale

import (
	"context"
	"time"

	"example.com/heroku/command"
)

func Load(client *command.Client) *Command {
	return &Command{client}
}

type Command struct {
	hc *command.Client
}

type Enable struct {
	command.Global
	Max           int           `flag:"max" help:"maximum number of dynos"`
	Min           int           `flag:"min" help:"minimum number of dynos"`
	Notifications bool          `flag:"notifications" default:"true" help:"receive email notifications when the max dyno limit is reached"`
	P95           time.Duration `flag:"p95" default:"2s" help:"desired p95 response time"`
}

// Enable web dyno autoscaling
func (c *Command) Enable(ctx context.Context, in *Enable) error {
	return nil
}

type Disable struct {
	command.Global
}

// Disable web dyno autoscaling
func (c *Command) Disable(ctx context.Context, in *Disable) error {
	return nil
}
