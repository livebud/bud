package sessions

type Controller struct {
}

func (c *Controller) New() string {
	return "new session"
}
