package users

type Controller struct {
}

func (c *Controller) Index() (string, error) {
	return "hello user", nil
}
