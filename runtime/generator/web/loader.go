package web

import (
	"io/fs"
	"path"
	"strings"

	"github.com/livebud/bud/internal/scan"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/vfs"
	"github.com/matthewmueller/text"
)

func Load(fsys fs.FS, module *gomod.Module, parser *parser.Parser) (*State, error) {
	loader := &loader{
		imports: imports.New(),
		fsys:    fsys,
		module:  module,
		parser:  parser,
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	imports *imports.Set
	fsys    fs.FS
	module  *gomod.Module
	parser  *parser.Parser
}

// Load the command state
func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	// Ensure the web files exist
	exist, err := vfs.SomeExist(l.fsys,
		"bud/.app/controller/controller.go",
		"bud/.app/public/public.go",
		"bud/.app/view/view.go",
	)
	if err != nil {
		return nil, err
	}
	// Add initial imports
	l.imports.AddStd("net/http", "context")
	l.imports.AddNamed("middleware", "github.com/livebud/bud/package/middleware")
	l.imports.AddNamed("web", "github.com/livebud/bud/runtime/web")
	l.imports.AddNamed("router", "github.com/livebud/bud/package/router")
	// Show the welcome page if we don't have controllers, views or public files
	if len(exist) == 0 {
		l.imports.AddNamed("welcome", "github.com/livebud/bud/runtime/web/welcome")
		state.ShowWelcome = true
		state.Imports = l.imports.List()
		return state, nil
	}
	// Turn on parts of the web server, based on what's generated
	if exist["bud/.app/public/public.go"] {
		state.HasPublic = true
		l.imports.AddNamed("public", l.module.Import("bud/.app/public"))
	}
	if exist["bud/.app/view/view.go"] {
		state.HasView = true
		l.imports.AddNamed("view", l.module.Import("bud/.app/view"))
	}
	// Load the controllers
	if exist["bud/.app/controller/controller.go"] {
		l.imports.AddNamed("controller", l.module.Import("bud/.app/controller"))
		state.Actions = l.loadControllerActions()
	}
	// state.Command = l.loadRoot("command")
	// Load the imports
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadControllerActions() (actions []*Action) {
	subfs, err := fs.Sub(l.fsys, "controller")
	if err != nil {
		l.Bail(err)
	}
	scanner := scan.Controllers(subfs)
	for scanner.Scan() {
		actions = append(actions, l.loadActions(scanner.Text())...)
	}
	if scanner.Err() != nil {
		l.Bail(err)
	}
	return actions
}

func (l *loader) loadActions(dir string) (actions []*Action) {
	pkg, err := l.parser.Parse(path.Join("controller", dir))
	if err != nil {
		l.Bail(err)
	}
	stct := pkg.Struct("Controller")
	if stct == nil {
		return nil
	}
	basePath := toBasePath(dir)
	for _, method := range stct.PublicMethods() {
		action := new(Action)
		actionName := method.Name()
		action.Method = l.loadActionMethod(actionName)
		action.Route = l.loadActionRoute(l.loadControllerRoute(basePath), actionName)
		action.CallName = l.loadActionCallName(basePath, actionName)
		actions = append(actions, action)
	}
	return actions
}

func toBasePath(dir string) string {
	if dir == "." {
		return "/"
	}
	return "/" + dir
}

// Method is the HTTP method for this controller
func (l *loader) loadActionMethod(name string) string {
	switch name {
	case "Create":
		return "Post"
	case "Update":
		return "Patch"
	case "Delete":
		return "Delete"
	default:
		return "Get"
	}
}

func (l *loader) loadControllerRoute(controllerPath string) string {
	segments := strings.Split(text.Path(controllerPath), "/")
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

// Route to the action
func (l *loader) loadActionRoute(basePath, actionName string) string {
	switch actionName {
	case "Show", "Update", "Delete":
		return path.Join(basePath, ":id")
	case "New":
		return path.Join(basePath, "new")
	case "Edit":
		return path.Join(basePath, ":id", "edit")
	case "Index", "Create":
		return basePath
	default:
		return path.Join(basePath, text.Lower(text.Path(actionName)))
	}
}

func (l *loader) loadActionCallName(basePath, actionName string) string {
	return text.Dot(text.Title(basePath + " " + actionName))
}
