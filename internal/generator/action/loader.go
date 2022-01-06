package action

import (
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"

	"gitlab.com/mnm/bud/internal/valid"
	"gitlab.com/mnm/bud/router"

	"github.com/matthewmueller/gotext"
	"github.com/matthewmueller/text"
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
}

// load fn
func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	state.Controller = l.loadController("action")
	state.Contexts = l.contexts.List()
	state.Imports = l.imports.List()
	return state, err
}

func (l *loader) loadController(actionPath string) *Controller {
	des, err := fs.ReadDir(l.module, actionPath)
	if err != nil {
		l.Bail(err)
	}
	controller := new(Controller)
	controllerPath := strings.TrimPrefix(actionPath, "action")
	controller.Name = l.loadControllerName(controllerPath)
	controller.Pascal = gotext.Pascal(controller.Name)
	// TODO: rename to route
	controller.Path = l.loadControllerPath(controllerPath)
	shouldParse := false
	for _, de := range des {
		if !de.IsDir() && valid.ActionFile(de.Name()) {
			shouldParse = true
			continue
		}
		if de.IsDir() && valid.Dir(de.Name()) {
			subController := l.loadController(path.Join(actionPath, de.Name()))
			if subController == nil {
				continue
			}
			controller.Controllers = append(controller.Controllers, subController)
			continue
		}
	}
	if !shouldParse {
		return controller
	}
	pkg, err := l.parser.Parse(actionPath)
	if err != nil {
		l.Bail(err)
	}
	stct := pkg.Struct("Controller")
	if stct == nil {
		return controller
	}
	controller.Actions = l.loadActions(controller, stct)
	return controller
}

func (l *loader) loadControllerName(controllerPath string) string {
	return text.Space(controllerPath)
}

func (l *loader) loadControllerPath(controllerPath string) string {
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
	if path.Len() == 0 {
		return ""
	}
	return path.String()
}

func (l *loader) loadActions(controller *Controller, stct *parser.Struct) (actions []*Action) {
	for _, method := range stct.PublicMethods() {
		actions = append(actions, l.loadAction(controller, method))
	}
	// Add the imports if we have more than one action
	if len(actions) > 0 {
		importPath, err := stct.File().Import()
		if err != nil {
			l.Bail(err)
		}
		l.imports.Add(importPath)
		l.imports.Add("gitlab.com/mnm/duo/response")
		l.imports.Add("net/http")
	}
	return actions
}

func (l *loader) loadAction(controller *Controller, method *parser.Function) *Action {
	action := new(Action)
	action.Name = method.Name()
	action.Pascal = gotext.Pascal(action.Name)
	action.Camel = gotext.Camel(action.Name)
	action.Short = text.Lower(gotext.Short(action.Name))
	// action.View = l.loadView(action.Name)
	action.Key = l.loadActionKey(action.Name)
	action.Path = l.loadActionPath(controller.Path, action.Name)
	action.Method = l.loadActionMethod(action.Name)
	action.Params = l.loadActionParams(method.Params())
	action.Input = l.loadActionInput(action.Params)
	action.Results = l.loadActionResults(method)
	action.ResponseJSON = len(action.Results) > 0
	action.Context = l.loadContext(controller, method)
	action.Redirect = l.loadActionRedirect(action)
	return action
}

func (l *loader) loadActionKey(actionName string) string {
	return "/" + text.Lower(text.Path(actionName))
}

// Path is the route to the action
func (l *loader) loadActionPath(controllerPath, actionName string) string {
	switch actionName {
	case "Show", "Update", "Delete":
		return "/" + path.Join(controllerPath, ":id")
	case "New":
		return "/" + path.Join(controllerPath, "new")
	case "Edit":
		return "/" + path.Join(controllerPath, ":id", "edit")
	case "Index", "Create":
		return "/" + controllerPath
	default:
		return "/" + path.Join(controllerPath, text.Path(actionName))
	}
}

// Method is the HTTP method for this controller
func (l *loader) loadActionMethod(actionName string) string {
	switch actionName {
	case "Create":
		return http.MethodPost
	case "Update":
		return http.MethodPatch
	case "Delete":
		return http.MethodDelete
	default:
		return http.MethodGet
	}
}

func (l *loader) loadActionParams(params []*parser.Param) (inputs []*ActionParam) {
	numParams := len(params)
	for nth, param := range params {
		inputs = append(inputs, l.loadActionParam(param, nth, numParams))
	}
	if len(inputs) > 0 {
		l.imports.Add("gitlab.com/mnm/duo/request")
	}
	return inputs
}

func (l *loader) loadActionParam(param *parser.Param, nth, numParams int) *ActionParam {
	dec, err := param.Definition()
	if err != nil {
		l.Bail(fmt.Errorf("action: unable to find param definition: %w", err))
	}
	ap := new(ActionParam)
	ap.Name = l.loadActionParamName(param, nth)
	ap.Pascal = gotext.Pascal(ap.Name)
	ap.Snake = gotext.Lower(gotext.Snake(ap.Name))
	ap.Type = l.loadType(param.Type(), dec)
	ap.Tag = fmt.Sprintf("`json:\"%[1]s\"`", tagValue(ap.Snake))
	ap.Kind = string(dec.Kind())
	// Single struct input
	if numParams == 1 && dec.Kind() == parser.KindStruct {
		ap.Variable = "in"
	} else {
		ap.Variable = "in." + ap.Pascal
	}
	return ap
}

func (l *loader) loadActionParamName(param *parser.Param, nth int) string {
	name := param.Name()
	if name != "" {
		return name
	}
	// Handle inputs with no variable
	return "in" + strconv.Itoa(nth)
}

func (l *loader) loadType(dt parser.Type, dec parser.Declaration) string {
	// TODO: Error out for certain built-ins (e.g. chan)
	if dec.Kind() == parser.KindBuiltin {
		return dt.String()
	}
	// Find the import path
	importPath, err := dec.Package().Import()
	if err != nil {
		l.Bail(err)
	}
	// Add the type's import
	name := l.imports.Add(importPath)
	dt = parser.Qualify(dt, name)
	return dt.String()
}

func (l *loader) loadActionInput(params []*ActionParam) string {
	if len(params) == 1 && params[0].Kind == parser.KindStruct {
		return params[0].Type
	}
	return l.loadActionInputStruct(params)
}

func (l *loader) loadActionInputStruct(params []*ActionParam) string {
	b := new(strings.Builder)
	b.WriteString("struct {")
	for _, param := range params {
		b.WriteString("\n")
		b.WriteString("\t\t" + param.Pascal)
		b.WriteString(" ")
		b.WriteString(param.Type)
		b.WriteString(" ")
		b.WriteString(param.Tag)
	}
	b.WriteString("\n\t}")
	return b.String()
}

func (l *loader) loadActionResults(method *parser.Function) (outputs []*ActionResult) {
	for order, result := range method.Results() {
		outputs = append(outputs, l.loadActionResult(order, result))
	}
	return outputs
}

func (l *loader) loadActionResult(order int, result *parser.Result) *ActionResult {
	output := new(ActionResult)
	output.Name = l.loadActionResultName(order, result)
	output.Pascal = gotext.Pascal(output.Name)
	output.Named = result.Named()
	output.Snake = gotext.Snake(output.Name)
	output.Type = result.Type().String()
	output.Variable = gotext.Camel(output.Name)
	output.Methods = l.loadActionResultMethods()
	output.Fields = l.loadActionResultFields()
	// TODO: check for other types that implement error
	output.IsError = output.Type == "error"
	return output
}

func (l *loader) loadActionResultName(order int, result *parser.Result) string {
	name := result.Name()
	if name != "" {
		return name
	}
	// Handle inputs with no variable
	return "in" + strconv.Itoa(order)
}

// TODO: Finish up
func (l *loader) loadActionResultFields() (fields []*ActionResultField) {
	return fields
}

// TODO: Finish up
func (l *loader) loadActionResultMethods() (methods []*ActionResultMethod) {
	return methods
}

// TODO: wrap this up
func (l *loader) loadActionRedirect(action *Action) string {
	switch action.Method {
	case http.MethodPatch, http.MethodDelete:
		return ""
		// return l.replacePath(action.Path, l.inputsToStrings(action.Inputs...)...)
	case http.MethodPost:
		return ""
		// return l.replacePath(action.Path, l.outputsToStrings(action.Outputs...)...)
	default: // Don't need to redirect on GET requests
		return ""
	}
}

func (l *loader) replacePath(route string, variables ...string) string {
	tokens := router.Parse(route)
	for _, token := range tokens {
		fmt.Println(token.String())
	}
	return route
}

func (l *loader) inputsToStrings(ins ...*ActionParam) (inputs []string) {
	for _, in := range ins {
		inputs = append(inputs, l.variableToString(in.Type, in.Variable))
	}
	return inputs
}

func (l *loader) outputsToStrings(results ...*ActionResult) (outputs []string) {
	for _, result := range results {
		outputs = append(outputs, l.variableToString(result.Type, result.Variable))
	}
	return outputs
}

func (l *loader) variableToString(dataType, variable string) string {
	switch dataType {
	case "int":
		l.imports.AddStd("strconv")
		return fmt.Sprintf(`strconv.Itoa(%s)`, variable)
	case "string":
		return variable
	}
	l.Bail(fmt.Errorf("action: unhandled type %q", dataType))
	return ""
}

func (l *loader) loadContext(controller *Controller, method *parser.Function) *Context {
	recv := method.Receiver()
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
	fnName := gotext.Camel("load " + controller.Name + " " + def.Name())
	provider, err := l.injector.Wire(&di.Function{
		Name:   fnName,
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
	context := new(Context)
	context.Function = fnName
	context.Code = provider.Function()
	context.Fields = l.loadContextInputs(provider)
	context.Results = l.loadContextResults(provider)
	// Add the context to the context set
	l.contexts.Add(context)
	return context
}

func (l *loader) loadContextInputs(provider *di.Provider) (fields []*ContextField) {
	for _, param := range provider.Externals {
		fields = append(fields, l.loadContextField(param))
	}
	return fields
}

func (l *loader) loadContextField(param *di.External) *ContextField {
	field := new(ContextField)
	field.Name = param.Key
	field.Variable = param.Name
	field.Hoisted = param.Hoisted
	field.Type = param.Type
	return field
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

func (l *loader) loadContextResults(provider *di.Provider) (outputs []*ContextResult) {
	for _, result := range provider.Results {
		outputs = append(outputs, l.loadContextResult(result))
	}
	return outputs
}

func (l *loader) loadContextResult(result *di.Variable) *ContextResult {
	output := new(ContextResult)
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

func tagValue(snake string) (out string) {
	if snake == "" {
		out += "-"
	} else {
		out += snake
		out += ",omitempty"
	}
	return out
}
