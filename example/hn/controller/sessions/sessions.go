package sessions

import "fmt"

type Controller struct {
}

func (c *Controller) New() {
}

func (c *Controller) Create(email, password string) error {
	return fmt.Errorf("not done yet")
}
