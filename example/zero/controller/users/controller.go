package users

type Controller struct {
}

func (c *Controller) Index() string {
	return "user index"
}

func (c *Controller) New() string {
	return "new user form"
}
