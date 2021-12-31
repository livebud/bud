package web

import (
	"io/fs"
	"path"

	"gitlab.com/mnm/bud/internal/scan"

	"github.com/matthewmueller/text"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/internal/parser"
	"gitlab.com/mnm/bud/ldflag"
	"gitlab.com/mnm/bud/vfs"
)

func Load(module *mod.Module, parser *parser.Parser) (*State, error) {
	exist := vfs.SomeExist(module,
		"bud/action/action.go",
		"bud/public/public.go",
		"bud/view/view.go",
	)
	if len(exist) == 0 {
		return nil, fs.ErrNotExist
	}
	loader := &loader{
		imports: imports.New(),
		module:  module,
		parser:  parser,
		exist:   exist,
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	imports *imports.Set
	module  *mod.Module
	parser  *parser.Parser
	exist   map[string]bool
}

// Load the command state
func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	// Add initial imports
	l.imports.AddStd("net", "net/http", "context")
	l.imports.AddNamed("hot", "gitlab.com/mnm/bud/hot")
	l.imports.AddNamed("middleware", "gitlab.com/mnm/bud/middleware")
	l.imports.AddNamed("web", "gitlab.com/mnm/bud/web")
	l.imports.AddNamed("router", "gitlab.com/mnm/bud/router")
	if l.exist["bud/public/public.go"] {
		state.HasPublic = true
		l.imports.AddNamed("public", l.module.Import("bud/public"))
	}
	if l.exist["bud/view/view.go"] {
		state.HasView = true
		l.imports.AddNamed("view", l.module.Import("bud/view"))
	}
	// Load the conditionals
	state.HasHot = ldflag.Hot()
	if l.exist["bud/action/action.go"] {
		l.imports.AddNamed("action", l.module.Import("bud/action"))
		state.Actions = l.loadControllerActions()
	}
	// state.Command = l.loadRoot("command")
	// Load the imports
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadControllerActions() (actions []*Action) {
	subfs, err := fs.Sub(l.module, "action")
	if err != nil {
		l.Bail(err)
	}
	scanner := scan.Actions(subfs)
	for scanner.Scan() {
		actions = append(actions, l.loadActions(scanner.Text())...)
	}
	if scanner.Err() != nil {
		l.Bail(err)
	}
	return actions
}

func (l *loader) loadActions(dir string) (actions []*Action) {
	pkg, err := l.parser.Parse(path.Join("action", dir))
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
		action.Route = l.loadActionRoute(basePath, actionName)
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
		return path.Join(basePath, text.Path(actionName))
	}
}

func (l *loader) loadActionCallName(basePath, actionName string) string {
	return text.Dot(text.Title(basePath + " " + actionName))
}
