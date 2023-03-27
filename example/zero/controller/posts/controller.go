package posts

import (
	"fmt"

	"github.com/livebud/bud/example/zero/session"
)

type Controller struct {
}

type IndexContext struct {
	Session *session.Session
}

func (c *Controller) Index(ctx *IndexContext) (string, error) {
	fmt.Println("got session", ctx.Session)
	return "Welcome to my blog!", nil
}
