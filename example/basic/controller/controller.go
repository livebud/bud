package controller

type Controller struct {
}

type Post struct {
	HTML string `json:"html"`
}

func (c *Controller) Index() string {
	return "hello world"
}

func (c *Controller) Show(id string) (*Post, error) {
	return &Post{
		HTML: "<b>hello " + id + "</b></script><script>console.log('log', '!', '!')</script>",
	}, nil
}
