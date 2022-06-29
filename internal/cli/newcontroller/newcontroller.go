package newcontroller

import (
	"context"
	_ "embed"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/internal/scaffold"
	"github.com/matthewmueller/gotext"
	"github.com/matthewmueller/text"
)

func New(bud *bud.Command) *Command {
	return &Command{bud: bud}
}

type Command struct {
	bud     *bud.Command
	Path    string
	Actions []string

	// Private
	bail bail.Struct
}

//go:embed controller.gotext
var controller string

//go:embed view_index.gotext
var indexView string

//go:embed view_show.gotext
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
	template   string
	Controller *Controller
	Path       string
	Title      string
	Variable   string
	Singular   string
	Plural     string
}

func (c *Command) Run(ctx context.Context) (err error) {
	state, err := c.Load()
	if err != nil {
		return err
	}
	return c.Scaffold(state)
}

func (c *Command) Load() (state *State, err error) {
	defer c.bail.Recover2(&err, "new controller")
	state = new(State)
	state.Controller = c.controller()
	state.Views = c.views(state.Controller)
	return state, nil
}

// Scaffold the files from state
func (c *Command) Scaffold(state *State) error {
	module, err := bud.Module(c.bud.Dir)
	if err != nil {
		return err
	}
	scaffolds := []scaffold.Scaffolding{
		scaffold.Template(state.Controller.path, controller, state.Controller),
	}
	for _, view := range state.Views {
		scaffolds = append(scaffolds, scaffold.Template(view.Path, view.template, view))
	}
	fsys := scaffold.MapFS{}
	if err := scaffold.Scaffold(fsys, scaffolds...); err != nil {
		return err
	}
	if err := scaffold.Write(fsys, module.Directory()); err != nil {
		return err
	}
	return nil
}

func (c *Command) controller() *Controller {
	controller := new(Controller)
	imports := imports.New()
	imports.AddStd("context")
	controller.Imports = imports.List()
	key, resource := splitKeyAndResource(c.Path)
	controller.key = key
	controller.path = controllerPath(key)
	controller.Name = controllerName(key)
	controller.Struct = gotext.Pascal(text.Singular(resource))
	controller.Plural = text.Plural(resource)
	controller.Singular = text.Singular(resource)
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
		c.bail.Bail(fmt.Errorf("invalid path:resource %q", a))
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
		template = ""
	}
	return &View{
		template:   template,
		Controller: controller,
		Path:       filepath.Join("view", controller.key, action.Name+".svelte"),
		Title:      text.Title(controller.Struct),
		Variable:   text.Camel(action.Result),
		Singular:   text.Camel(text.Singular(action.Result)),
		Plural:     text.Camel(text.Plural(action.Result)),
	}
}

func splitKeyAndResource(rel string) (p string, r string) {
	rel = strings.ToLower(rel)
	parts := strings.SplitN(rel, ":", 2)
	if len(parts) == 1 {
		key := controllerKey(rel)
		resource := path.Base(key)
		if resource == "." {
			return key, "resource"
		}
		return key, resource
	}
	key := controllerKey(parts[1])
	return key, parts[0]
}

func controllerKey(path string) string {
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		return ""
	}
	path = filepath.ToSlash(path)
	return path
}

func controllerPath(controllerKey string) string {
	// TODO: change this back to controller
	return filepath.Join("controller", controllerKey, "controller.go")
}

func controllerName(controllerKey string) string {
	name := path.Base(controllerKey)
	if name == "." {
		return "controller"
	}
	return name
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
	return strings.TrimSuffix("/"+path.String(), "/")
}

func viewPath(controllerKey, path string) string {
	return filepath.Join("view", controllerKey, text.Snake(strings.ToLower(path))+".svelte")
}
