package action

import (
	"context"

	"gitlab.com/mnm/bud/example/hn/internal/hn"
)

type Controller struct {
	HN *hn.Client
}

func (c *Controller) Index(ctx context.Context) (*hn.News, error) {
	return c.HN.FrontPage(ctx)
}

func (c *Controller) Show(ctx context.Context, id string) (*ShowResponse, error) {
	story, err := c.HN.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	return &ShowResponse{story}, nil
}

type ShowResponse struct {
	Story *hn.Story `json:"story"`
}
