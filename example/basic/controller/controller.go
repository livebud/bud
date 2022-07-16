package controller

type Controller struct {
}

type Post struct {
	ID   string
	HTML string `json:"html"`
}

func (c *Controller) Index() string {
	return "hello world"
}

func (c *Controller) Show(id string) (*Post, error) {
	return &Post{
		ID:   id,
		HTML: "<b>hello " + id + "</b></script><script>console.log('log', '!', '!')</script>",
	}, nil
}

func (c *Controller) Update(id, html string) (*Post, error) {
	return &Post{
		ID: id,
	}, nil
}
