package action

import (
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/matthewmueller/gotext"
	"github.com/matthewmueller/text"
	"gitlab.com/mnm/bud/go/is"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/internal/parser"
	"gitlab.com/mnm/bud/vfs"
)

func Load(injector *di.Injector, module *mod.Module, parser *parser.Parser) (*State, error) {
	exist := vfs.SomeExist(module, "action")
	if len(exist) == 0 {
		return nil, fs.ErrNotExist
	}
	loader := &loader{
		contexts: newContextSet(),
		imports:  imports.New(),
		injector: injector,
		module:   module,
		parser:   parser,
		exist:    exist,
	}
	return loader.Load()
}

// loader struct
type loader struct {
	bail.Struct
	injector *di.Injector
	imports  *imports.Set
	contexts *contextSet
	module   *mod.Module
	parser   *parser.Parser
	exist    map[string]bool

	// cached
	controllerPath string
}

// load fn
func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	state.Controller = l.loadController()
	state.Contexts = l.contexts.List()
	state.Imports = l.imports.List()
	return state, err
}

func (l *loader) loadController() *Controller {
	pkg, err := l.parser.Parse("./action")
	if err != nil {
		l.Bail(err)
	}
	l.controllerPath = l.loadControllerPath(pkg)
	controller := new(Controller)
	controller.Name = pkg.Name()
	controller.Path = l.loadControllerPath(pkg)
	controller.Actions = l.loadActions(pkg)
	return controller
}

func (l *loader) loadControllerPath(pkg *parser.Package) string {
	controllerDir := filepath.Join(l.module.Directory(), "controller")
	rel := strings.TrimPrefix(pkg.Directory(), controllerDir)
	return text.Path(rel)
}

func (l *loader) loadActions(pkg *parser.Package) (actions []*Action) {
	for _, fn := range pkg.PublicFunctions() {
		actions = append(actions, l.loadAction(fn))
	}
	// Add the imports if we have more than one action
	if len(actions) > 0 {
		importPath, err := pkg.Import()
		if err != nil {
			l.Bail(err)
		}
		l.imports.Add(importPath)
		l.imports.Add("gitlab.com/mnm/duo/response")
		l.imports.Add("net/http")
	}
	return actions
}

func (l *loader) loadAction(fn *parser.Function) *Action {
	action := new(Action)
	action.Name = fn.Name()
	action.Pascal = gotext.Pascal(action.Name)
	action.Camel = gotext.Camel(action.Name)
	action.Short = gotext.Short(action.Name)
	// action.View = l.loadView(action.Name)
	action.Key = text.Path(action.Name)
	action.Path = l.loadActionPath(action.Name)
	action.Method = l.loadActionMethod(action.Name)
	action.Inputs = l.loadActionInputs(fn)
	action.Outputs = l.loadActionOutputs(fn)
	action.ResponseJSON = len(action.Outputs) > 0
	action.Context = l.loadContext(fn)
	return action
}

// View returns the view
// func (l *loader) loadView(actionName string) *View {
// 	if l.view == nil {
// 		return nil
// 	}
// 	actionPath := filepath.Join(l.controllerPath, gotext.Slug(actionName))
// 	viewDir := filepath.Join(l.module.Directory(), "view")
// 	viewPath := filepath.Join(viewDir, actionPath)
// 	// Lookup the generated views individual views
// 	for _, prereq := range l.view.Prerequisites() {
// 		if !strings.HasPrefix(prereq.Path(), viewPath) {
// 			continue
// 		}
// 		l.imports.Add(path.Join(l.module.Import(), "generated", "view"))
// 		rel, err := filepath.Rel(viewDir, prereq.Path())
// 		if err != nil {
// 			l.Bail(fmt.Errorf("plugin/controller: unable to make the view path relative: %w", err))
// 			return nil
// 		}
// 		return &View{
// 			Path: rel,
// 		}
// 	}
// 	l.Bail(fmt.Errorf("plugin/controller: view is not a prerequisite to %s: %q", actionName, l.view.Path()))
// 	return nil
// }

// Path is the route to the action
func (l *loader) loadActionPath(actionName string) string {
	basePath := l.controllerPath
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

// Method is the HTTP method for this controller
func (l *loader) loadActionMethod(actionName string) string {
	switch actionName {
	case "Create":
		return "POST"
	case "Update":
		return "PATCH"
	case "Delete":
		return "DELETE"
	default:
		return "GET"
	}
}

func (l *loader) loadActionInputs(fn *parser.Function) (inputs []*ActionInput) {
	for order, param := range fn.Params() {
		inputs = append(inputs, l.loadActionInput(order, param))
	}
	if len(inputs) > 0 {
		l.imports.Add("gitlab.com/mnm/duo/request")
	}
	return inputs
}

func (l *loader) loadActionInput(order int, param *parser.Param) *ActionInput {
	input := new(ActionInput)
	input.Name = l.loadActionInputName(order, param)
	input.Pascal = gotext.Pascal(input.Name)
	input.Snake = gotext.Snake(input.Name)
	input.Type = l.loadActionInputType(param)
	input.Variable = "in." + input.Pascal
	input.Snake = l.loadActionInputJSON(input.Snake)
	return input
}

func (l *loader) loadActionInputName(order int, param *parser.Param) string {
	name := param.Name()
	if name != "" {
		return name
	}
	// Handle inputs with no variable
	return "in" + strconv.Itoa(order)
}

func (l *loader) loadActionInputType(param *parser.Param) string {
	dt := param.Type()
	dtString := dt.String()
	// Do nothing with built-in types
	// TODO: It can't be any built-in type (e.g. chan)
	if is.Builtin(dtString) {
		return dtString
	}
	// Find the import path
	importPath, err := param.File().Import()
	if err != nil {
		l.Bail(err)
	}
	// Add the type's import
	name := l.imports.Add(importPath)
	dt = parser.Qualify(dt, name)
	return dt.String()
}

func (l *loader) loadActionInputJSON(snake string) string {
	out := "`json:\""
	if snake == "" {
		out += "-"
	} else {
		out += snake
		out += ",omitempty"
	}
	out += "\"`"
	return out
}

func (l *loader) loadActionOutputs(fn *parser.Function) (outputs []*ActionOutput) {
	for order, result := range fn.Results() {
		outputs = append(outputs, l.loadActionOutput(order, result))
	}
	return outputs
}

func (l *loader) loadActionOutput(order int, result *parser.Result) *ActionOutput {
	output := new(ActionOutput)
	output.Name = l.loadActionOutputName(order, result)
	output.Pascal = gotext.Pascal(output.Name)
	output.Named = result.Named()
	output.Snake = gotext.Snake(output.Name)
	output.Type = result.Type().String()
	output.Variable = gotext.Camel(output.Name)
	output.Methods = l.loadActionOutputMethods()
	output.Fields = l.loadActionOutputFields()
	// TODO: check for other types that implement error
	output.IsError = output.Type == "error"
	return output
}

func (l *loader) loadActionOutputName(order int, result *parser.Result) string {
	name := result.Name()
	if name != "" {
		return name
	}
	// Handle inputs with no variable
	return "in" + strconv.Itoa(order)
}

// TODO: Finish up
func (l *loader) loadActionOutputMethods() (methods []*ActionOutputMethod) {
	return methods
}

// TODO: Finish up
func (l *loader) loadActionOutputFields() (fields []*ActionOutputField) {
	return fields
}

func (l *loader) loadContext(fn *parser.Function) *Context {
	recv := fn.Receiver()
	if recv == nil {
		return nil
	}
	def, err := recv.Definition()
	if err != nil {
		l.Bail(err)
	}
	importPath, err := def.Package().Import()
	if err != nil {
		l.Bail(err)
	}
	provider, err := l.injector.Wire(&di.Function{
		Name:   "load" + def.Name(),
		Target: l.module.Import("bud", "action"),
		Hoist:  true,
		Results: []di.Dependency{
			&di.Type{
				Import: importPath,
				Type:   recv.Type().String(),
			},
		},
		Params: []di.Dependency{
			&di.Type{
				Import: "net/http",
				Type:   "ResponseWriter",
			},
			&di.Type{
				Import: "net/http",
				Type:   "*Request",
			},
		},
	})
	if err != nil {
		l.Bail(err)
	}
	// Add generated imports
	for _, imp := range provider.Imports {
		l.imports.AddNamed(imp.Name, imp.Path)
	}
	// Create the context
	fnName := "load" + def.Name()
	context := new(Context)
	context.Function = fnName
	context.Code = provider.Function()
	context.Inputs = l.loadContextInputs(provider)
	context.Outputs = l.loadContextOutputs(provider)
	// Add the context to the context set
	l.contexts.Add(context)
	return context
}

func (l *loader) loadContextInputs(provider *di.Provider) (inputs []*ContextInput) {
	for _, param := range provider.Externals {
		inputs = append(inputs, l.loadContextInput(param))
	}
	return inputs
}

func (l *loader) loadContextInput(param *di.External) *ContextInput {
	input := new(ContextInput)
	input.Name = param.Key
	input.Variable = param.Name
	input.Hoisted = param.Hoisted
	input.Type = param.Type
	return input
}

// func (l *loader) loadContextInputName(dataType string) (typeName string) {
// 	parts := strings.Split(dataType, ".")
// 	if len(parts) > 1 {
// 		typeName = parts[len(parts)-1]
// 	} else {
// 		typeName = parts[0]
// 	}
// 	return strings.TrimLeft(typeName, "[]*")
// }

func (l *loader) loadContextOutputs(provider *di.Provider) (outputs []*ContextOutput) {
	for _, result := range provider.Results {
		outputs = append(outputs, l.loadContextOutput(result))
	}
	return outputs
}

func (l *loader) loadContextOutput(result *di.Variable) *ContextOutput {
	output := new(ContextOutput)
	output.Variable = gotext.Camel(result.Name)
	return output
}

func newContextSet() *contextSet {
	return &contextSet{map[string]*Context{}}
}

type contextSet struct {
	contextMap map[string]*Context
}

func (c *contextSet) Add(context *Context) {
	c.contextMap[context.Function] = context
}

func (c *contextSet) List() (contexts []*Context) {
	for _, context := range c.contextMap {
		contexts = append(contexts, context)
	}
	sort.Slice(contexts, func(i, j int) bool {
		return contexts[i].Function < contexts[j].Function
	})
	return contexts
}
