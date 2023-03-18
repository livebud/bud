package controller

import (
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/gotemplate"
	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/runtime/generator"
)

func New(module *gomod.Module) *Generator {
	return &Generator{module}
}

type Generator struct {
	module *gomod.Module
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
	posts *posts.Controller,
	sessions *sessions.Controller,
	users *users.Controller,
) *Controller {
	return &Controller{
		&PostsController{
			&PostsIndexAction{posts},
		},
		&SessionsController{
			&SessionsNewAction{sessions},
		},
		&UsersController{
			&UsersIndexAction{users},
			&UsersNewAction{users},
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
}

func (a *PostsIndexAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := a.controller.Index()
	w.Write([]byte(res))
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
}

func (a *UsersIndexAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := a.controller.Index()
	w.Write([]byte(res))
}

type UsersNewAction struct {
	controller *users.Controller
}

func (a *UsersNewAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := a.controller.New()
	w.Write([]byte(res))
}
`

var gen = gotemplate.MustParse("controller.gotext", template)

type State struct {
	Imports []*imports.Import
}

func (g *Generator) generateFile(fsys generator.FS, file *generator.File) error {
	imset := imports.New()
	imset.AddStd("net/http")
	imset.AddNamed("router", "github.com/livebud/bud/package/router")
	imset.AddNamed("posts", g.module.Import("controller/posts"))
	imset.AddNamed("users", g.module.Import("controller/users"))
	imset.AddNamed("sessions", g.module.Import("controller/sessions"))
	code, err := gen.Generate(&State{
		Imports: imset.List(),
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
