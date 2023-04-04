package addons

import (
	"context"

	"example.com/heroku/command"
)

type Command struct {
	// Dependencies
}

type Addons struct {
	*command.Global
	All  bool `flag:"all" default:"false" help:"show add-ons and attachments for all accessible apps"`
	Json bool `flag:"app" default:"false" help:"return add-ons in json format"`
}

// Addons lists your add-ons and attachments
func (c *Command) Addons(ctx context.Context, in *Addons) error {
	return nil
}

type Attach struct {
	*command.Global
	Name string  `arg:"addon_name" help:"name of the addon"`
	As   *string `flag:"as" help:"name for add-on attachment"`
}

// Attach an existing add-on resource to an app
func (c *Command) Attach(ctx context.Context, in *Attach) error {
	return nil
}

type Create struct {
	*command.Global
	ServicePlan string  `arg:"service_plan" help:"addon service plan"`
	As          *string `flag:"as" help:"name for add-on attachment"`
	Name        *string `flag:"name" help:"name for the add-on resource"`
	Wait        bool    `flag:"wait" default:"false" help:"watch add-on creation status and exit when complete"`
}

func (c *Command) Create(ctx context.Context, in *Create) error {
	return nil
}

type Destroy struct {
	*command.Global
	Force bool `flag:"force" default:"false" help:"force destroy"`
}

// Destroy permanently destroys an add-on resource
func (c *Command) Destroy(ctx context.Context, in *Destroy) error {
	return nil
}

type Info struct {
	*command.Global
}

// Info shows detailed add-on resource and attachment information
func (c *Command) Info(ctx context.Context, in *Info) error {
	return nil
}

type Services struct {
	Json bool `flag:"app" default:"false" help:"app to run command against"`
}

func (c *Command) Services(ctx context.Context, in *Services) error {
	return nil
}
