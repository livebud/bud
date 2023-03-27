package controller

import (
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/gotemplate"
	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/runtime/generator"
)

func New(injector *di.Injector, module *gomod.Module) *Generator {
	return &Generator{injector, module}
}

type Generator struct {
	injector *di.Injector
	module   *gomod.Module
}

func (g *Generator) Extend(gen generator.FileSystem) {
	gen.GenerateFile("bud/pkg/web/controller/controller.go", g.generateFile)
}

const template = `package controller

{{- if $.Imports }}

import (
	{{- range $import := $.Imports }}
	{{ $import.Name }} "{{ $import.Path }}"
	{{- end }}
)
{{- end }}

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
	Posts *PostsController
	Sessions *SessionsController
	Users *UsersController
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
	view *view.View
}

func (a *PostsIndexAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	propMap := map[string]interface{}{}
	indexContext, err := {{ $.PostIndexProvider.Name}}(
		{{- if $.PostIndexProvider.Variable "context.Context" }}r.Context(),{{ end }}
		{{- if $.PostIndexProvider.Variable "net/http.*Request" }}r,{{ end }}
		{{- if $.PostIndexProvider.Variable "net/http.ResponseWriter" }}w,{{ end }}
	)
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
	view *view.View
}

func (a *SessionsNewAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := a.controller.New()
	w.Write([]byte(res))
}

type UsersController struct {
	Index *UsersIndexAction
	New *UsersNewAction
}

func (c *UsersController) Mount(r *router.Router) error {
	r.Get("/users/new", c.New)
	r.Get("/users", c.Index)
	return nil
}

type UsersIndexAction struct {
	controller *users.Controller
	view *view.View
}

func (a *UsersIndexAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := a.controller.Index()
	w.Write([]byte(res))
}

type UsersNewAction struct {
	controller *users.Controller
	view *view.View
}

func (a *UsersNewAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := a.controller.New()
	w.Write([]byte(res))
}

{{ $.PostIndexProvider.Function }}
`

var gen = gotemplate.MustParse("controller.gotext", template)

type State struct {
	Imports           []*imports.Import
	PostIndexProvider *di.Provider
}

func (g *Generator) generateFile(fsys generator.FS, file *generator.File) error {
	imset := imports.New()
	imset.AddStd("net/http")
	imset.AddNamed("router", "github.com/livebud/bud/package/router")
	imset.AddNamed("posts", g.module.Import("controller/posts"))
	imset.AddNamed("users", g.module.Import("controller/users"))
	imset.AddNamed("sessions", g.module.Import("controller/sessions"))
	imset.AddNamed("view", g.module.Import("bud/pkg/web/view"))

	provider, err := g.injector.Wire(&di.Function{
		Name:    "loadPostsIndexContext",
		Hoist:   true,
		Imports: imset,
		Target:  g.module.Import("bud/pkg/web/controller"),
		Params: []*di.Param{
			{
				Import: "net/http",
				Type:   "*Request",
			},
			{
				Import: "net/http",
				Type:   "ResponseWriter",
			},
			{
				Import: "context",
				Type:   "Context",
			},
		},
		Aliases: di.Aliases{
			di.ToType(g.module.Import("env"), "*Env"):         di.ToType(g.module.Import("bud/internal/env"), "*Env"),
			di.ToType(g.module.Import("session"), "*Session"): di.ToType(g.module.Import("bud/pkg/sessions"), "*Session"),
		},
		Results: []di.Dependency{
			di.ToType(g.module.Import("controller/posts"), "*IndexContext"),
			&di.Error{},
		},
	})
	if err != nil {
		return err
	}
	// fmt.Println(provider.Function())

	code, err := gen.Generate(&State{
		Imports:           imset.List(),
		PostIndexProvider: provider,
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
