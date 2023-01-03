package newcontroller

import (
	"context"
	_ "embed"
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/config"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/internal/scaffold"
	"github.com/matthewmueller/gotext"
	"github.com/matthewmueller/text"
)

func New(provide config.Provide) *Command {
	return &Command{provide: provide}
}

type Command struct {
	provide config.Provide

	Path    string
	Actions []string

	// Private
	bail bail.Struct
}

//go:embed controller.gotext
var controller string

//go:embed view_index.gotext
var indexView string

//go:embed view_new.gotext
var newView string

//go:embed view_edit.gotext
var editView string

//go:embed view_show.gotext
var showView string

var views = map[string]string{
	"index": indexView,
	"new":   newView,
	"edit":  editView,
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

	// Paths
	IndexPath string
	EditPath  string
	ShowPath  string
	NewPath   string
}

type Action struct {
	Name   string
	Route  string
	Result string

	Index  bool
	Create bool
	New    bool
	Show   bool
	Edit   bool
	Update bool
	Delete bool
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
	state.Controller = c.loadController()
	state.Views = c.loadViews(state.Controller)
	return state, nil
}

// Scaffold the files from state
func (c *Command) Scaffold(state *State) error {
	module, err := c.provide.Module()
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

func (c *Command) loadController() *Controller {
	controller := new(Controller)
	imports := imports.New()
	imports.AddStd("context")
	controller.Imports = imports.List()
	key, resource := splitKeyAndResource(c.Path)
	// TODO: remove this constraint
	if strings.Contains(key, "/") && hasOneOrMore(c.Actions, "index", "new") {
		c.bail.Bail(fmt.Errorf(`scaffolding the "index" or "new" action of a nested resource like %q isn't supported yet, see https://github.com/livebud/bud/issues/209 for details`, c.Path))
	}
	controller.key = key
	controller.path = controllerPath(key)
	controller.Name = controllerName(key)
	controller.Struct = gotext.Pascal(text.Singular(resource))
	controller.Plural = text.Plural(resource)
	controller.Singular = text.Singular(resource)
	controller.Route = controllerRoute(key)
	controller.Package = gotext.Snake(controller.Name)
	controller.Pascal = gotext.Pascal(controller.Name)
	// Load paths
	controller.IndexPath = controllerIndexPath(controller.key, controller.Singular)
	controller.ShowPath = controllerShowPath(controller.IndexPath, controller.Singular)
	controller.NewPath = controllerNewPath(controller.IndexPath, controller.Singular)
	controller.EditPath = controllerEditPath(controller.ShowPath, controller.Singular)
	// Load the actions
	for _, action := range c.Actions {
		controller.Actions = append(controller.Actions, c.loadControllerAction(controller, action))
	}
	// Consistent action names
	sort.Slice(controller.Actions, func(i, j int) bool {
		left := actionRank[controller.Actions[i].Name]
		right := actionRank[controller.Actions[j].Name]
		return left > right
	})
	return controller
}

// Rank in the order in which the actions are typically called
var actionRank = map[string]int{
	"index":  7,
	"new":    6,
	"create": 5,
	"show":   4,
	"edit":   3,
	"update": 2,
	"delete": 1,
}

func (c *Command) loadControllerAction(controller *Controller, a string) *Action {
	action := new(Action)
	action.Name = strings.ToLower(a)
	switch action.Name {
	case "index":
		action.Index = true
		action.Route = controller.Route
		action.Result = gotext.Camel(controller.Plural)
	case "new":
		action.New = true
		action.Route = path.Join(controller.Route, "new")
	case "create":
		action.Create = true
		action.Route = controller.Route
		action.Result = gotext.Camel(controller.Singular)
	case "show":
		action.Show = true
		action.Route = path.Join(controller.Route, ":id")
		action.Result = gotext.Camel(controller.Singular)
	case "edit":
		action.Edit = true
		action.Route = path.Join(controller.Route, ":id/edit")
		action.Result = gotext.Camel(controller.Singular)
	case "update":
		action.Update = true
		action.Route = path.Join(controller.Route, ":id")
		action.Result = gotext.Camel(controller.Singular)
	case "delete":
		action.Delete = true
		action.Route = path.Join(controller.Route, ":id")
		action.Result = gotext.Camel(controller.Singular)
	default:
		c.bail.Bail(fmt.Errorf(`invalid action %q, expected "index", "new", "create", "show", "edit", "update" or "delete"`, a))
	}
	return action
}

func (c *Command) loadViews(controller *Controller) (views []*View) {
	for _, action := range controller.Actions {
		view := c.loadView(controller, action)
		if view == nil {
			continue
		}
		views = append(views, view)
	}
	return views
}

func (c *Command) loadView(controller *Controller, action *Action) *View {
	template, ok := views[action.Name]
	if !ok {
		return nil
	}
	return &View{
		template:   template,
		Controller: controller,
		Path:       filepath.Join("view", controller.key, action.Name+".svelte"),
		Title:      text.Title(controller.Struct),
		Variable:   text.Camel(action.Result),
		Singular:   text.Camel(controller.Singular),
		Plural:     text.Camel(controller.Plural),
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

func controllerIndexPath(controllerKey, propVar string) string {
	segments := strings.Split(controllerKey, "/")
	path := new(strings.Builder)
	for i := 0; i < len(segments); i++ {
		if i%2 != 0 {
			path.WriteString("/")
			path.WriteString("${" + propVar + "." + text.Slug(text.Singular(segments[i-1])) + "_id || 0}")
			path.WriteString("/")
		}
		path.WriteString(text.Slug(segments[i]))
	}
	return "/" + strings.TrimSuffix(path.String(), "/")
}

func controllerNewPath(controllerIndexPath, propVar string) string {
	return path.Join(controllerIndexPath, "new")
}

func controllerShowPath(controllerIndexPath, propVar string) string {
	return path.Join(controllerIndexPath, "${"+propVar+".id || 0}")
}

func controllerEditPath(controllerShowPath, propVar string) string {
	return path.Join(controllerShowPath, "edit")
}

func hasOneOrMore(actions []string, matches ...string) bool {
	for _, a := range actions {
		for _, m := range matches {
			if a == m {
				return true
			}
		}
	}
	return false
}
