package controller

import (
	"context"

	"github.com/matthewmueller/hackernews"
)

func New(hn *hackernews.Client) *Controller {
	return &Controller{hn}
}

type Controller struct {
	hn *hackernews.Client
}

func (c *Controller) Index(ctx context.Context) (stories []*hackernews.Story, err error) {
	return c.hn.FrontPage(ctx)
}

// Show a comment
func (c *Controller) Show(ctx context.Context, id int) (story *hackernews.Story, err error) {
	return c.hn.Find(ctx, id)
}
