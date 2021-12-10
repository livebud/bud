package new

import "fmt"

type Command struct {
}

type ViewInput struct {
	Name     string `arg:"name" help:"name of the view"`
	WithTest bool   `flag:"with-test" help:"include a view test" default:"true"`
}

func (c *Command) View(in *ViewInput) error {
	fmt.Println("creating new view", in.Name, in.WithTest)
	return nil
}
