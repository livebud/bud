package resource

import (
	"context"
	_ "embed"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"gitlab.com/mnm/bud/package/vfs"

	"github.com/matthewmueller/gotext"

	"github.com/matthewmueller/text"
	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/gomod"
	"gitlab.com/mnm/bud/package/scaffold"
)

func New(module *gomod.Module) *Command {
	return &Command{
		module:   module,
		template: scaffold.Template{},
	}
}

type Command struct {
	bail.Struct
	module   *gomod.Module
	template scaffold.Template
	Path     string
	Actions  []string
}

//go:embed controller.gotext
var controller string

//go:embed view/default.gotext
var defaultView string

//go:embed view/index.gotext
var indexView string

//go:embed view/show.gotext
var showView string

var views = map[string]string{
	"index": indexView,
	"show":  showView,
}

type State struct {
	Controller *Controller
	Views      []*View
}

type Controller struct {
	Imports  []*imports.Import
	key      string
	path     string
	Package  string
	Name     string
	Pascal   string
	Struct   string
	Route    string
	Plural   string
	Singular string
	Actions  []*Action
}

type Action struct {
	Name   string
	Route  string
	Result string

	Index bool
	Show  bool
}

type View struct {
	template string
	Path     string
	Title    string
	Variable string
	Singular string
}

func (c *Command) Run(ctx context.Context) (err error) {
	state, err := c.Load()
	if err != nil {
		return err
	}
	return Generate(c.module.DirFS(), state)
}

func (c *Command) Load() (state *State, err error) {
	defer c.Recover2(&err, "new resource error")
	state = new(State)
	state.Controller = c.controller()
	state.Views = c.views(state.Controller)
	return state, nil
}

func (c *Command) controller() *Controller {
	controller := new(Controller)
	imports := imports.New()
	imports.AddStd("context")
	controller.Imports = imports.List()
	controller.key = controllerKey(c.Path)
	controller.path = controllerPath(controller.key)
	controller.Name = path.Base(controller.key)
	if controller.Name == "." {
		controller.Name = "controller"
		controller.Struct = "Resource"
		controller.Plural = "resources"
		controller.Singular = "resource"
	} else {
		controller.Struct = gotext.Pascal(text.Singular(controller.Name))
		controller.Plural = text.Plural(controller.Name)
		controller.Singular = text.Singular(controller.Name)
	}
	controller.Route = controllerRoute(controller.key)
	controller.Package = gotext.Snake(controller.Name)
	controller.Pascal = gotext.Pascal(controller.Name)
	for _, action := range c.Actions {
		controller.Actions = append(controller.Actions, c.controllerAction(controller, action))
	}
	return controller
}

func (c *Command) controllerAction(controller *Controller, a string) *Action {
	action := new(Action)
	action.Name = strings.ToLower(a)
	switch action.Name {
	case "index":
		action.Index = true
		action.Route = controller.Route
		// action.View = &View{
		// 	Path:     filepath.Join("view", controller.key, "index.svelte"),
		// 	Template: index,
		// }
		action.Result = gotext.Camel(controller.Plural)
	case "show":
		action.Show = true
		action.Route = path.Join(controller.Route, "/:id")
		action.Result = gotext.Camel(controller.Singular)
	default:
		c.Bail(fmt.Errorf("new resource: action not implemented yet %q", a))
	}
	return action
}

func (c *Command) views(controller *Controller) (views []*View) {
	for _, action := range controller.Actions {
		views = append(views, c.view(controller, action))
	}
	return views
}

func (c *Command) view(controller *Controller, action *Action) *View {
	template, ok := views[action.Name]
	if !ok {
		template = defaultView
	}
	return &View{
		template: template,
		Path:     filepath.Join("view", controller.key, action.Name+".svelte"),
		Title:    text.Title(controller.Struct),
		Variable: text.Camel(action.Result),
		Singular: text.Camel(text.Singular(action.Result)),
	}
}

// Generate
func Generate(fsys vfs.ReadWritable, state *State) error {
	var templates scaffold.Templates
	templates = append(templates, &scaffold.Template{
		Path:  state.Controller.path,
		Code:  controller,
		State: state.Controller,
	})
	for _, view := range state.Views {
		templates = append(templates, &scaffold.Template{
			Path:  view.Path,
			Code:  view.template,
			State: view,
		})
	}
	return templates.Write(fsys)
}

func controllerKey(path string) string {
	path = strings.ToLower(path)
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		return ""
	}
	path = filepath.ToSlash(path)
	return path
}

func controllerPath(controllerKey string) string {
	// TODO: change this back to controller
	return filepath.Join("action", controllerKey, "controller.go")
}

// TODO: dedupe with controllerRoute in runtime/generator/web
func controllerRoute(controllerKey string) string {
	segments := strings.Split(controllerKey, "/")
	path := new(strings.Builder)
	for i := 0; i < len(segments); i++ {
		if i%2 != 0 {
			path.WriteString("/")
			path.WriteString(":" + text.Slug(text.Singular(segments[i-1])) + "_id")
			path.WriteString("/")
		}
		path.WriteString(text.Slug(segments[i]))
	}
	return "/" + path.String()
}

func viewPath(controllerKey, path string) string {
	return filepath.Join("view", controllerKey, text.Snake(strings.ToLower(path))+".svelte")
}
