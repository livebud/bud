package controller

type Controller struct {
}

func (c *Controller) Index() string {
	return "hello world."
}

func (c *Controller) Show(id string) string {
	return "shows/" + id
}
