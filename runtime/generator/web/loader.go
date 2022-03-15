package web

import (
	"io/fs"
	"path"
	"strings"

	"gitlab.com/mnm/bud/internal/scan"

	"github.com/matthewmueller/text"
	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/internal/ldflag"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
	"gitlab.com/mnm/bud/pkg/vfs"
)

func Load(fsys fs.FS, module *gomod.Module, parser *parser.Parser) (*State, error) {
	exist := vfs.SomeExist(fsys,
		"bud/.app/action/action.go",
		"bud/.app/public/public.go",
		"bud/.app/view/view.go",
	)
	if len(exist) == 0 {
		return nil, fs.ErrNotExist
	}
	loader := &loader{
		imports: imports.New(),
		fsys:    fsys,
		module:  module,
		parser:  parser,
		exist:   exist,
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	imports *imports.Set
	fsys    fs.FS
	module  *gomod.Module
	parser  *parser.Parser
	exist   map[string]bool
}

// Load the command state
func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	// Add initial imports
	l.imports.AddStd("net", "net/http", "context")
	l.imports.AddNamed("hot", "gitlab.com/mnm/bud/runtime/hot")
	l.imports.AddNamed("middleware", "gitlab.com/mnm/bud/pkg/middleware")
	l.imports.AddNamed("web", "gitlab.com/mnm/bud/runtime/web")
	l.imports.AddNamed("router", "gitlab.com/mnm/bud/pkg/router")
	if l.exist["bud/.app/public/public.go"] {
		state.HasPublic = true
		l.imports.AddNamed("public", l.module.Import("bud/.app/public"))
	}
	if l.exist["bud/.app/view/view.go"] {
		state.HasView = true
		l.imports.AddNamed("view", l.module.Import("bud/.app/view"))
	}
	// Load the conditionals
	state.HasHot = ldflag.Hot()
	if l.exist["bud/.app/action/action.go"] {
		l.imports.AddNamed("action", l.module.Import("bud/.app/action"))
		state.Actions = l.loadControllerActions()
	}
	// state.Command = l.loadRoot("command")
	// Load the imports
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadControllerActions() (actions []*Action) {
	subfs, err := fs.Sub(l.fsys, "action")
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
		return path.Join(basePath, text.Path(actionName))
	}
}

func (l *loader) loadActionCallName(basePath, actionName string) string {
	return text.Dot(text.Title(basePath + " " + actionName))
}
