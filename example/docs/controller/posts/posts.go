package posts

import (
	"fmt"
	"net/http"

	"github.com/livebud/bud/framework/controller/controllerrt/request"
	"github.com/livebud/bud/package/log/console"
)

type Controller struct {
}

type Post struct {
	Title string
	Body  string
}

func (c *Controller) Frame(id int) (likes int) {
	// Another use case (related posts)
	return 10
}

func (c *Controller) show(id int) (*Post, error) {
	return &Post{
		Title: "Hello world!",
		Body:  "This is a post.",
	}, nil
}

func (c *Controller) ShowJSON(id int) (*Post, error) {
	return c.show(id)
}

// type layoutRequest struct {
// 	Theme string `json:"theme,omitempty"`
// }

// type layoutResponse struct {
// 	Title string `json:"title,omitempty"`
// 	Style string `json:"style,omitempty"`
// }

// type layoutView struct {
// 	Controller *controller.Controller
// 	Request    layoutRequest
// 	Response   layoutResponse
// }

// func (l *layoutView) Run(ctx context.Context) error {

// }

// Generated (by hand). Do not Edit.
func (c *Controller) Show(w http.ResponseWriter, r *http.Request) {
	var layout struct {
		Theme string `json:"theme,omitempty"`
	}
	var frame struct {
	}
	var postsFrame struct {
		ID int `json:"id,omitempty"`
	}
	var postsPage struct {
		ID int `json:"id,omitempty"`
	}
	// TODO: consider a request umarshaller that takes multiple structs
	if err := request.Unmarshal(r, &layout); err != nil {
		console.Err(err, "unable to marshal")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := request.Unmarshal(r, &frame); err != nil {
		console.Err(err, "unable to frame")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := request.Unmarshal(r, &postsFrame); err != nil {
		console.Err(err, "unable to posts/frame")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := request.Unmarshal(r, &postsPage); err != nil {
		console.Err(err, "unable to posts/page")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(layout, frame, postsFrame, postsPage)
	// _, _, _, _ = Layout, Frame, PostsFrame, PostsPage
	w.Write([]byte("Hello world!!!"))
}
