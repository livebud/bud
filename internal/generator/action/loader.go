package action

import (
	"errors"
	"fmt"
	"io/fs"
	"path"

	"gitlab.com/mnm/bud/internal/valid"

	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/internal/parser"
)

func Load(injector *di.Injector, module *mod.Module, parser *parser.Parser) (*State, error) {
	loader := &loader{
		contexts: newContextSet(),
		imports:  imports.New(),
		injector: injector,
		module:   module,
		parser:   parser,
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	imports  *imports.Set
	contexts *contextSet
	injector *di.Injector
	module   *mod.Module
	parser   *parser.Parser
}

// Load the command state
func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	// TODO: move these
	l.imports.AddStd("net/http", "context")
	// l.imports.AddStd("fmt")
	l.imports.AddNamed("controller", "gitlab.com/mnm/bud/controller")
	l.imports.AddNamed("view", "gitlab.com/mnm/bud/view")
	l.imports.AddNamed("action", l.module.Import("action"))
	l.imports.AddNamed("users", l.module.Import("action", "users"))
	// state.Controllers = l.loadControllers("action", "")
	// state.Imports = l.imports.List()
	// return state, nil
	return &State{
		Imports: l.imports.List(),
		Controller: &Controller{
			Name:  "",
			Route: "/",
			Import: &imports.Import{
				Name: "action",
				Path: l.module.Import("action"),
			},
			// Context: &Context{},
			Actions: []*Action{
				&Action{
					Name:   "Index",
					Method: "GET",
					Route:  "/",
					View:   true,
					Params: ActionParams{
						&ActionParam{
							Name: "ctx",
							Type: "context.Context",
						},
					},
					Results: ActionResults{
						&ActionResult{
							Name: "",
							Type: "*hn.News",
						},
						&ActionResult{
							Name: "",
							Type: "error",
						},
					},
				},
				&Action{
					Name:   "Show",
					Method: "GET",
					Route:  "/:id",
					View:   true,
					Params: ActionParams{
						&ActionParam{
							Name: "ctx",
							Type: "context.Context",
						},
						&ActionParam{
							Name: "id",
							Type: "string",
						},
					},
					Results: ActionResults{
						&ActionResult{
							Name: "",
							Type: "*ShowResponse",
						},
						&ActionResult{
							Name: "",
							Type: "error",
						},
					},
				},
			},
			Controllers: []*Controller{
				&Controller{
					Name:  "users",
					Route: "/users",
					Import: &imports.Import{
						Name: "users",
						Path: l.module.Import("action/users"),
					},
					// Context: &Context{},
					Actions: []*Action{
						&Action{
							Name:   "Index",
							Method: "GET",
							Results: ActionResults{
								&ActionResult{
									Name: "",
									Type: "error",
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

func (l *loader) loadControllers(base, dir string) (controllers []*Controller) {
	des, err := fs.ReadDir(l.module, path.Join(base, dir))
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			l.Bail(err)
		}
		return controllers
	}
	hasActionFiles := false
	for _, de := range des {
		// Recurse the subdirectories
		if de.IsDir() && valid.Dir(de.Name()) {
			subdir := path.Join(dir, de.Name())
			controllers = append(controllers, l.loadControllers(base, subdir)...)
			continue
		}
		// Count action files
		if !de.IsDir() && valid.ActionFile(de.Name()) {
			hasActionFiles = true
			continue
		}
	}
	if hasActionFiles {
		controller := l.loadController(base, dir)
		if controller != nil {
			controllers = append(controllers, controller)
		}
	}
	return controllers
}

func (l *loader) loadController(base, dir string) *Controller {
	pkg, err := l.parser.Parse(path.Join(base, dir))
	if err != nil {
		l.Bail(err)
	}
	stct := pkg.Struct("Controller")
	if stct == nil || len(stct.PublicMethods()) == 0 {
		return nil
	}
	controller := new(Controller)
	controller.Name = pkg.Name()
	// controller.Path = "/" + dir
	controller.Actions = l.loadActions(stct)
	// controller.Context = l.contexts.List()
	fmt.Println("controller path", pkg.Name(), stct.Name())
	return controller
}

func (l *loader) loadActions(stct *parser.Struct) (actions []*Action) {
	return actions
}
