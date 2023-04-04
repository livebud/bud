// Package ps lists dynos for your app
package ps

import (
	"context"
	"time"

	"example.com/heroku/command"
)

func New(client *command.Client) *Command {
	return &Command{hc: client}
}

type Command struct {
	hc *command.Client
}

type Ps struct {
	*command.Global
	Json bool `flag:"app" default:"false" help:"app to run command against"`
}

// Ps lists dynos for an app
func (c *Command) Ps(ctx context.Context, in *Ps) error {
	return nil
}

type Exec struct {
	*command.Global
	Dyno   string `flag:"dyno" help:"specify the dyno to connect to"`
	SSH    bool   `flag:"ssh" default:"false" help:"use native ssh"`
	Status bool   `flag:"status" default:"false" help:"lists the status of the SSH server in the dyno"`
}

// Exec creates an SSH session to a dyno
func (c *Command) Exec(ctx context.Context, in *Exec) error {
	return nil
}

type Scale struct {
	*command.Global
	Expr string `arg:"expr" default:"" help:"expression to scale by"`
}

// Scale dyno quantities up or down
func (c *Command) Scale(ctx context.Context, in *Scale) error {
	return nil
}

type Stop struct {
	*command.Global
	Dyno string `arg:"dyno" help:"stop a dyno"`
}

// Stop app dyno
func (c *Command) Stop(ctx context.Context, in *Stop) error {
	return nil
}

type Wait struct {
	*command.Global
	Dyno         string        `arg:"dyno" help:"stop a dyno"`
	WithRun      bool          `flag:"with-run" short:"R" default:"false" help:"wait for dyno to be running"`
	Type         string        `flag:"type" short:"t" help:"wait for one specific dyno type"`
	WaitInterval time.Duration `flag:"wait-interval" short:"w" default:"10s" help:"how frequently to poll in seconds"`
}

// Wait for all dynos to be running latest version after a release
func (c *Command) Wait(ctx context.Context, in *Wait) error {
	return nil
}
