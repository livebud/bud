package controller

import (
	context "context"
	sessions1 "github.com/livebud/bud/example/zero/bud/pkg/sessions"
	view "github.com/livebud/bud/example/zero/bud/pkg/web/view"
	posts "github.com/livebud/bud/example/zero/controller/posts"
	sessions "github.com/livebud/bud/example/zero/controller/sessions"
	users "github.com/livebud/bud/example/zero/controller/users"
	router "github.com/livebud/bud/package/router"
	http "net/http"
)

func New(
	view *view.View,
	posts *posts.Controller,
	sessions *sessions.Controller,
	users *users.Controller,
) *Controller {
	return &Controller{
		&PostsController{
			&PostsIndexAction{posts, view},
		},
		&SessionsController{
			&SessionsNewAction{sessions, view},
		},
		&UsersController{
			&UsersIndexAction{users, view},
			&UsersNewAction{users, view},
		},
	}
}

type Controller struct {
	Posts    *PostsController
	Sessions *SessionsController
	Users    *UsersController
}

// TODO: use a router.Router interface
func (c *Controller) Mount(r *router.Router) error {
	c.Posts.Mount(r)
	c.Sessions.Mount(r)
	c.Users.Mount(r)
	return nil
}

type PostsController struct {
	Index *PostsIndexAction
}

func (c *PostsController) Mount(r *router.Router) error {
	r.Get("/posts", c.Index)
	return nil
}

type PostsIndexAction struct {
	controller *posts.Controller
	view       *view.View
}

func (a *PostsIndexAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	propMap := map[string]interface{}{}
	indexContext, err := loadPostsIndexContext(r.Context())
	if err != nil {
		html := a.view.RenderError(ctx, "posts/index", propMap, err)
		w.WriteHeader(500)
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(html))
		return
	}
	res, err := a.controller.Index(indexContext)
	if err != nil {
		html := a.view.RenderError(ctx, "posts/index", propMap, err)
		w.WriteHeader(500)
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(html))
		return
	}
	propMap["posts/index"] = res
	html, err := a.view.Render(ctx, "posts/index", propMap)
	if err != nil {
		html = a.view.RenderError(ctx, "posts/index", propMap, err)
	}
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(html))
}

type SessionsController struct {
	New *SessionsNewAction
}

func (c *SessionsController) Mount(r *router.Router) error {
	r.Get("/sessions/new", c.New)
	return nil
}

type SessionsNewAction struct {
	controller *sessions.Controller
	view       *view.View
}

func (a *SessionsNewAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := a.controller.New()
	w.Write([]byte(res))
}

type UsersController struct {
	Index *UsersIndexAction
	New   *UsersNewAction
}

func (c *UsersController) Mount(r *router.Router) error {
	r.Get("/users/new", c.New)
	r.Get("/users", c.Index)
	return nil
}

type UsersIndexAction struct {
	controller *users.Controller
	view       *view.View
}

func (a *UsersIndexAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := a.controller.Index()
	w.Write([]byte(res))
}

type UsersNewAction struct {
	controller *users.Controller
	view       *view.View
}

func (a *UsersNewAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := a.controller.New()
	w.Write([]byte(res))
}

func loadPostsIndexContext(contextContext context.Context) (*posts.IndexContext, error) {
	sessions1Session, err := sessions1.From(contextContext)
	if err != nil {
		return nil, err
	}
	postsIndexContext := &posts.IndexContext{Session: sessions1Session}
	return postsIndexContext, err
}
