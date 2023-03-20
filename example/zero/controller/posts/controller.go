package posts

type Controller struct {
}

func (c *Controller) Index() (string, error) {
	return "post index", nil
}
