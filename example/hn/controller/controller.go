package controller

import (
	"context"

	"github.com/matthewmueller/hackernews"
)

type Controller struct {
	HN *hackernews.Client
}

func (c *Controller) Index(ctx context.Context) (stories []*hackernews.Story, err error) {
	return c.HN.FrontPage(ctx)
}

// Show a comment
func (c *Controller) Show(ctx context.Context, id int) (story *hackernews.Story, err error) {
	return c.HN.Find(ctx, id)
}
