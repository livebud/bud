package controller

type Controller struct {
}

func (c *Controller) Layout(theme *string) (title, style string) {
	if theme == nil {
		*theme = "light"
	}
	return "This title is from the controller!", *theme
}

func (c *Controller) Frame() (categories []string) {
	return []string{
		"Finance",
		"Sports",
		"Politics",
		"Entertainment",
		"Technology",
	}
}
