package command

type Global struct {
	App    string  `flag:"app" short:"a" help:"app to run command against"`
	Remote *string `flag:"remote" short:"r" help:"git remote of app to use"`
}

func New() *Client {
	return &Client{"api.heroku.com"}
}

type Client struct {
	url string
}
