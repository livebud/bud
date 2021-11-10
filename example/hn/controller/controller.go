package controller

import (
	"context"
	"fmt"
)

type Controller struct {
}

func (c *Controller) Index(ctx context.Context) ([]interface{}, error) {
	return nil, fmt.Errorf("/ is not implemented yet")
}

func (c *Controller) Show(ctx context.Context, id int) error {
	return fmt.Errorf("GET /%d is not implemented yet", id)
}
