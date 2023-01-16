package cli

import (
	"context"
	_ "embed"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/livebud/bud/internal/gotemplate"

	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/virtual"
	"github.com/matthewmueller/gotext"
	"github.com/matthewmueller/text"
)

var newControllerViews = map[string]string{
	"index": "new_controller/view_index.gotext",
	"new":   "new_controller/view_new.gotext",
	"edit":  "new_controller/view_edit.gotext",
	"show":  "new_controller/view_show.gotext",
}

type NewController struct {
	Path    string
	Actions []string
}

func (c *CLI) NewController(ctx context.Context, in *NewController) error {
	var files []*virtual.File

	// Generate the controller
	controller, err := c.newController(in)
	if err != nil {
		return err
	}

	// Check for restrictions on nested resources
	if strings.Contains(controller.key, "/") {
		for _, action := range controller.Actions {
			if action.Name == "index" || action.Name == "new" {
				return fmt.Errorf(`new controller: scaffolding the "index" or "new" action of a nested resource like %q isn't supported yet, see https://github.com/livebud/bud/issues/209 for details`, in.Path)
			}
		}
	}

	file, err := controller.Generate()
	if err != nil {
		return err
	}
	files = append(files, file)

	// Generate the views
	for _, action := range controller.Actions {
		// Skip over generating views that don't have a known action
		if _, ok := newControllerViews[action.Name]; !ok {
			continue
		}
		view := &newControllerView{
			ActionName: action.Name,
			Controller: controller,
			Path:       filepath.Join("view", controller.key, action.Name+".svelte"),
			Title:      text.Title(controller.Struct),
			Variable:   text.Camel(action.Result),
			Singular:   text.Camel(controller.Singular),
			Plural:     text.Camel(controller.Plural),
		}
		file, err := view.Generate()
		if err != nil {
			return err
		}
		files = append(files, file)
	}

	// Find the module
	module, err := c.findModule()
	if err != nil {
		return err
	}

	// Write the files out
	for _, file := range files {
		if err := module.MkdirAll(filepath.Dir(file.Path), 0755); err != nil {
			return err
		}
		if err := module.WriteFile(file.Path, file.Data, 0644); err != nil {
			return err
		}
	}

	return nil
}

// Rank in the order in which the actions are typically called
var newControllerActionRank = map[string]int{
	"index":  7,
	"new":    6,
	"create": 5,
	"show":   4,
	"edit":   3,
	"update": 2,
	"delete": 1,
}

func splitKeyAndResource(rel string) (key string, resource string) {
	rel = strings.ToLower(rel)
	parts := strings.SplitN(rel, ":", 2)
	if len(parts) == 1 {
		key := filepath.ToSlash(strings.TrimPrefix(rel, "/"))
		resource := path.Base(key)
		if resource == "." {
			return key, "resource"
		}
		return key, resource
	}
	key = filepath.ToSlash(strings.TrimPrefix(parts[1], "/"))
	return key, parts[0]
}

func (c *CLI) newController(in *NewController) (*newController, error) {
	controller := new(newController)
	imports := imports.New()
	imports.AddStd("context")
	controller.Imports = imports.List()
	key, resource := splitKeyAndResource(in.Path)
	controller.key = key
	controller.path = filepath.Join("controller", key, "controller.go")
	controller.Struct = gotext.Pascal(text.Singular(resource))
	controller.Plural = text.Plural(resource)
	controller.Singular = text.Singular(resource)
	controller.Package = gotext.Snake(controller.Name())
	controller.Pascal = gotext.Pascal(controller.Name())

	// Load paths
	controller.ShowPath = path.Join(controller.IndexPath(), "${"+controller.Singular+".id || 0}")
	controller.NewPath = path.Join(controller.IndexPath(), "new")
	controller.EditPath = path.Join(controller.ShowPath, "edit")

	// Load the actions
	for _, action := range in.Actions {
		controllerAction, err := c.newControllerAction(controller, action)
		if err != nil {
			return nil, err
		}
		controller.Actions = append(controller.Actions, controllerAction)
	}

	// Consistent action names
	sort.Slice(controller.Actions, func(i, j int) bool {
		left := newControllerActionRank[controller.Actions[i].Name]
		right := newControllerActionRank[controller.Actions[j].Name]
		return left > right
	})
	return controller, nil
}

type newController struct {
	Imports  []*imports.Import
	key      string
	path     string
	name     string
	Package  string
	Pascal   string
	route    string
	Struct   string
	Plural   string
	Singular string
	Actions  []*newControllerAction

	// Paths
	indexPath string
	EditPath  string
	ShowPath  string
	NewPath   string
}

func (c *newController) Name() string {
	if c.name != "" {
		return c.name
	}
	name := path.Base(c.key)
	if name == "." {
		return "controller"
	}
	c.name = name
	return c.name
}

func (c *newController) Route() string {
	if c.route != "" {
		return c.route
	}
	segments := strings.Split(c.key, "/")
	path := new(strings.Builder)
	for i := 0; i < len(segments); i++ {
		if i%2 != 0 {
			path.WriteString("/")
			path.WriteString(":" + text.Slug(text.Singular(segments[i-1])) + "_id")
			path.WriteString("/")
		}
		path.WriteString(text.Slug(segments[i]))
	}
	c.route = strings.TrimSuffix("/"+path.String(), "/")
	return c.route
}

func (c *newController) IndexPath() string {
	if c.indexPath != "" {
		return c.indexPath
	}
	segments := strings.Split(c.key, "/")
	path := new(strings.Builder)

	for i := 0; i < len(segments); i++ {
		if i%2 != 0 {
			path.WriteString("/")
			path.WriteString("${" + c.Singular + "." + text.Slug(text.Singular(segments[i-1])) + "_id || 0}")
			path.WriteString("/")
		}
		path.WriteString(text.Slug(segments[i]))
	}
	c.indexPath = "/" + strings.TrimSuffix(path.String(), "/")
	return c.indexPath
}

func (c *newController) Generate() (*virtual.File, error) {
	code, err := fs.ReadFile(embedFS, "new_controller/controller.gotext")
	if err != nil {
		return nil, err
	}
	template, err := gotemplate.Parse("new_controller/controller.gotext", string(code))
	if err != nil {
		return nil, err
	}
	data, err := template.Generate(c)
	if err != nil {
		return nil, err
	}
	return &virtual.File{
		Path: c.path,
		Data: data,
	}, nil
}

func (c *CLI) newControllerAction(controller *newController, a string) (*newControllerAction, error) {
	action := new(newControllerAction)
	action.Name = strings.ToLower(a)
	switch action.Name {
	case "index":
		action.Index = true
		action.Route = controller.Route()
		action.Result = gotext.Camel(controller.Plural)
	case "new":
		action.New = true
		action.Route = path.Join(controller.Route(), "new")
	case "create":
		action.Create = true
		action.Route = controller.Route()
		action.Result = gotext.Camel(controller.Singular)
	case "show":
		action.Show = true
		action.Route = path.Join(controller.Route(), ":id")
		action.Result = gotext.Camel(controller.Singular)
	case "edit":
		action.Edit = true
		action.Route = path.Join(controller.Route(), ":id/edit")
		action.Result = gotext.Camel(controller.Singular)
	case "update":
		action.Update = true
		action.Route = path.Join(controller.Route(), ":id")
		action.Result = gotext.Camel(controller.Singular)
	case "delete":
		action.Delete = true
		action.Route = path.Join(controller.Route(), ":id")
		action.Result = gotext.Camel(controller.Singular)
	default:
		return nil, fmt.Errorf(`new controller: invalid action %q, expected "index", "new", "create", "show", "edit", "update" or "delete"`, a)
	}
	return action, nil
}

type newControllerAction struct {
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

type newControllerView struct {
	ActionName string
	Controller *newController
	Path       string
	Title      string
	Variable   string
	Singular   string
	Plural     string
}

func (c *newControllerView) Generate() (*virtual.File, error) {
	path, ok := newControllerViews[c.ActionName]
	if !ok {
		return nil, fmt.Errorf("unable to find view template for %q", c.ActionName)
	}
	code, err := fs.ReadFile(embedFS, path)
	if err != nil {
		return nil, err
	}
	template, err := gotemplate.Parse(path, string(code))
	if err != nil {
		return nil, err
	}
	data, err := template.Generate(c)
	if err != nil {
		return nil, err
	}
	return &virtual.File{
		Path: c.Path,
		Data: data,
	}, nil
}
