package posts

import (
	"fmt"
	"net/http"

	"github.com/livebud/bud/example/docs/controller"
	"github.com/livebud/bud/framework/controller/controllerrt/request"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/viewer"
)

// func NewViewer(svelteViewer *svelte.Viewer) *Viewer {
// 	return &Viewer{}
// }

type SvelteViewer struct{}
type HTMLViewer struct{}

type Viewer struct {
}

func New(layoutFrame *controller.Controller, viewer *Viewer) *Controller {
	return &Controller{
		layoutFrame: layoutFrame,
		viewer:      viewer,
	}
}

type Controller struct {
	layoutFrame *controller.Controller
	viewer      *Viewer
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

	// This part should happen concurrently
	title, style := c.layoutFrame.Layout(&layout.Theme)
	layoutProps := viewer.Props{
		"title": title,
		"style": style,
	}
	categories := c.layoutFrame.Frame()
	frame0Props := viewer.Props{
		"categories": categories,
	}
	likes := c.Frame(postsFrame.ID)
	frame1Props := viewer.Props{
		"likes": likes,
	}
	post, err := c.show(postsPage.ID)
	if err != nil {
		console.Err(err, "unable to show post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pageProps := viewer.Props{
		"post": post,
	}
	// Once done, we come back together for rendering
	_ = layoutProps
	_ = frame0Props
	_ = frame1Props
	_ = pageProps
	fmt.Println("got post", post)
	// c.Post()
	// fmt.Println(layout, frame, postsFrame, postsPage)
	// _, _, _, _ = Layout, Frame, PostsFrame, PostsPage
	w.Write([]byte("Hello world!!!"))
}
