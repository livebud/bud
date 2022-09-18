package controller

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/livebud/bud/internal/valid"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
	"github.com/matthewmueller/gotext"
	"github.com/matthewmueller/text"
)

func Load(fsys fs.FS, injector *di.Injector, module *gomod.Module, parser *parser.Parser) (*State, error) {
	if files, err := fs.Glob(fsys, "{controller/**.go,view/**}"); err != nil {
		return nil, err
	} else if len(files) == 0 {
		return nil, fs.ErrNotExist
	}
	loader := &loader{
		fsys:      fsys,
		providers: newProviderSet(),
		imports:   imports.New(),
		injector:  injector,
		module:    module,
		parser:    parser,
	}
	return loader.Load()
}

// loader struct
type loader struct {
	bail.Struct
	fsys      fs.FS
	injector  *di.Injector
	imports   *imports.Set
	providers *providerSet
	module    *gomod.Module
	parser    *parser.Parser
}

// load fn
func (l *loader) Load() (state *State, err error) {
	defer l.Recover2(&err, "controller: unable to load")
	state = new(State)
	state.Controller = l.loadController("controller")
	state.Providers = l.providers.List()
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadController(controllerPath string) *Controller {
	des, err := fs.ReadDir(l.fsys, controllerPath)
	if err != nil {
		l.Bail(err)
	} else if len(des) == 0 {
		l.Bail(fs.ErrNotExist)
	}
	controller := new(Controller)
	controller.Path = l.loadControllerPath(controllerPath)
	controller.Name = l.loadControllerName(controller.Path)
	controller.Pascal = gotext.Pascal(controller.Name)
	// TODO: rename to route
	controller.Route = l.loadControllerRoute(controller.Path)
	shouldParse := false
	for _, de := range des {
		if !de.IsDir() && valid.ControllerFile(de.Name()) {
			shouldParse = true
			continue
		}
		if de.IsDir() && valid.Dir(de.Name()) {
			subController := l.loadController(path.Join(controllerPath, de.Name()))
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
	pkg, err := l.parser.Parse(controllerPath)
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

func (l *loader) loadControllerPath(controllerPath string) string {
	parts := strings.SplitN(controllerPath, "/", 2)
	if len(parts) == 1 {
		return "/"
	}
	return "/" + parts[1]
}

func (l *loader) loadControllerName(controllerPath string) string {
	return text.Space(controllerPath)
}

func (l *loader) loadControllerRoute(controllerPath string) string {
	segments := strings.Split(strings.TrimPrefix(controllerPath, "/"), "/")
	path := new(strings.Builder)
	for i := 0; i < len(segments); i++ {
		if i%2 != 0 {
			path.WriteString("/")
			path.WriteString(":" + text.Snake(text.Singular(segments[i-1])) + "_id")
			path.WriteString("/")
		}
		path.WriteString(text.Snake(segments[i]))
	}
	if path.Len() == 0 {
		return "/"
	}
	return "/" + path.String()
}

func (l *loader) loadActions(controller *Controller, stct *parser.Struct) (actions []*Action) {
	var usesResponse bool
	for _, method := range stct.PublicMethods() {
		action := l.loadAction(controller, method)
		if !action.HandlerFunc {
			usesResponse = true
		}
		actions = append(actions, action)
	}
	// Add the imports if we have more than one action
	if len(actions) > 0 {
		importPath, err := stct.File().Import()
		if err != nil {
			l.Bail(err)
		}
		l.imports.Add(importPath)
		l.imports.Add("net/http")
		if usesResponse {
			l.imports.Add("github.com/livebud/bud/framework/controller/controllerrt/response")
		}
	}
	return actions
}

func (l *loader) loadAction(controller *Controller, method *parser.Function) *Action {
	action := new(Action)
	action.Name = method.Name()
	action.Pascal = gotext.Pascal(action.Name)
	action.Camel = gotext.Camel(action.Name)
	action.Short = text.Lower(gotext.Short(action.Name))
	action.Route = l.loadActionRoute(controller.Route, action.Name)
	action.Key = l.loadActionKey(controller.Path, action.Name)
	action.View = l.loadView(controller.Path, action.Key, action.Route)
	action.Method = l.loadActionMethod(action.Name)
	params := method.Params()
	results := method.Results()
	action.HandlerFunc = l.isHandlerFunc(params, results)
	if !action.HandlerFunc {
		action.Params = l.loadActionParams(params)
		action.Input = l.loadActionInput(action.Params)
		action.Results = l.loadActionResults(results)
	}
	action.RespondJSON = len(action.Results) > 0
	action.RespondHTML = l.loadRespondHTML(action.Results)
	action.Provider = l.loadProvider(controller, method)
	action.Redirect = l.loadActionRedirect(action)
	return action
}

func (l *loader) loadActionKey(controllerPath, actionName string) string {
	return path.Join(controllerPath, text.Lower(text.Snake(actionName)))
}

func (l *loader) isHandlerFunc(params []*parser.Param, results []*parser.Result) bool {
	if len(params) != 2 || len(results) != 0 {
		return false
	}

	maybeW := params[0].Type()
	isW, err := parser.IsImportType(maybeW, "net/http", "ResponseWriter")
	if err != nil {
		l.Bail(err)
	}

	if !isW {
		return false
	}

	maybeR := params[1].Type()
	isR, err := parser.IsImportType(maybeR, "net/http", "Request")
	if err != nil {
		l.Bail(err)
	}

	if !isR {
		return false
	}

	return true
}

// Route to the action
func (l *loader) loadActionRoute(controllerRoute, actionName string) string {
	switch actionName {
	case "Show", "Update", "Delete":
		return path.Join(controllerRoute, ":id")
	case "New":
		return path.Join(controllerRoute, "new")
	case "Edit":
		return path.Join(controllerRoute, ":id", "edit")
	case "Index", "Create":
		return controllerRoute
	default:
		return path.Join(controllerRoute, text.Lower(text.Snake(actionName)))
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

func (l *loader) loadView(controllerKey, actionKey, actionRoute string) *View {
	viewDir := path.Join("view", controllerKey)
	des, err := fs.ReadDir(l.fsys, viewDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		l.Bail(fmt.Errorf("controller: unable read view directory %q . %w", viewDir, err))
	}
	for _, de := range des {
		name := de.Name()
		ext := path.Ext(name)
		if ext != ".svelte" {
			continue
		}
		base := strings.TrimSuffix(path.Base(name), ext)
		key := path.Base(actionKey)
		if base != key {
			continue
		}
		l.imports.Add(l.module.Import("bud/internal/app/view"))
		return &View{
			Route: actionRoute,
		}
	}
	return nil
}

func (l *loader) loadActionParams(params []*parser.Param) (inputs []*ActionParam) {
	numParams := len(params)
	for nth, param := range params {
		inputs = append(inputs, l.loadActionParam(param, nth, numParams))
	}
	if len(inputs) > 0 {
		l.imports.Add("github.com/livebud/bud/framework/controller/controllerrt/request")
	}
	return inputs
}

func (l *loader) loadActionParam(param *parser.Param, nth, numParams int) *ActionParam {
	dec, err := param.Definition()
	if err != nil {
		l.Bail(fmt.Errorf("controller: unable to find param definition for %s. %w", param.Type(), err))
	}
	ap := new(ActionParam)
	ap.Name = l.loadActionParamName(param, nth)
	ap.Pascal = gotext.Pascal(ap.Name)
	ap.Snake = gotext.Lower(gotext.Snake(ap.Name))
	ap.Type = l.loadType(param.Type(), dec)
	ap.Tag = fmt.Sprintf("`json:\"%[1]s\"`", tagValue(ap.Snake))
	ap.Kind = string(dec.Kind())
	switch {
	// Single struct input
	case numParams == 1 && dec.Kind() == parser.KindStruct:
		ap.Variable = "in"
	// Handle context.Context
	case ap.IsContext():
		ap.Variable = `httpRequest.Context()`
	default:
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
	// Standard library
	if strings.HasPrefix(importPath, "std/") {
		dt := parser.Requalify(dt, imports.AssumedName(importPath))
		return dt.String()
	}
	// Add the type's import
	name := l.imports.Add(importPath)
	dt = parser.Qualify(dt, name)
	return dt.String()
}

func (l *loader) loadActionInput(params []*ActionParam) string {
	if len(params) == 1 && params[0].Kind == string(parser.KindStruct) {
		return params[0].Type
	}
	return l.loadActionInputStruct(params)
}

func (l *loader) loadActionInputStruct(params []*ActionParam) string {
	b := new(strings.Builder)
	b.WriteString("struct {")
	for _, param := range params {
		if param.IsContext() {
			continue
		}
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

func (l *loader) loadActionResults(results []*parser.Result) (outputs []*ActionResult) {
	for order, result := range results {
		outputs = append(outputs, l.loadActionResult(order, result))
	}
	return outputs
}

func (l *loader) loadActionResult(order int, result *parser.Result) *ActionResult {
	def, err := result.Definition()
	if err != nil {
		l.Bail(fmt.Errorf("controller: unable to load result definition for %s . %w", result.Type(), err))
	}
	output := new(ActionResult)
	output.Name = l.loadActionResultName(order, result)
	output.Pascal = gotext.Pascal(output.Name)
	output.Named = result.Named()
	output.Snake = gotext.Snake(output.Name)
	output.Type = parser.Unqualify(result.Type()).String()
	output.Kind = def.Kind()
	output.Variable = l.loadActionResultVariable(order, result)
	output.Fields = l.loadActionResultFields(result, def)
	// TODO: check for other types that implement error
	output.IsError = output.Type == "error"
	return output
}

func (l *loader) loadActionResultName(order int, result *parser.Result) string {
	name := result.Name()
	if name != "" {
		return name
	}
	// Handle outputs with no variable
	return "out" + strconv.Itoa(order)
}

// Load the variable name
func (l *loader) loadActionResultVariable(order int, result *parser.Result) string {
	name := result.Name()
	if name != "" {
		return name + "Out"
	}
	// Handle outputs with no variable
	return "out" + strconv.Itoa(order)
}

// Load the inner fields of the result type, if it's a struct
func (l *loader) loadActionResultFields(result *parser.Result, def parser.Declaration) (fields []*ActionResultField) {
	// Fields should be empty if the definition isn't a struct
	if def.Kind() != parser.KindStruct {
		return fields
	}
	// Find the struct in the package
	stct := def.Package().Struct(def.Name())
	if stct == nil {
		l.Bail(fmt.Errorf("controller: unable to find struct for %s", result.Type()))
	}
	for _, field := range stct.PublicFields() {
		fields = append(fields, &ActionResultField{
			Name: field.Name(),
			Type: field.Type().String(),
		})
	}
	return fields
}

// TODO: Clean this up, the logic is quite complicated and could be simplified
// with better methods
func (l *loader) loadActionRedirect(action *Action) string {
	// Redirect for non-create methods is an empty string
	if action.Method != http.MethodPost {
		return `""`
	}
	results := action.Results
	if isSingleStruct(results) {
		for _, field := range results[0].Fields {
			if field.Name != "ID" {
				continue
			}
			return l.variableToString(field.Type, results[0].Variable+"."+field.Name)
		}
	}
	for _, result := range results {
		if result.Name != "id" {
			continue
		}
		return l.variableToString(result.Type, result.Variable)
	}
	return `""`
}

func (l *loader) variableToString(dataType string, variable string) string {
	switch dataType {
	case "string":
		return variable
	case "int":
		l.imports.AddStd("strconv")
		return `strconv.Itoa(` + variable + `)`
	default:
		l.Bail(fmt.Errorf("controller: unable to generate string from %s", dataType))
		return ""
	}
}

func isSingleStruct(results ActionResults) bool {
	switch len(results) {
	case 0:
		return false
	case 1:
		result := results[0]
		if result.IsError {
			return false
		}
		return result.Kind == parser.KindStruct
	case 2:
		if !results[1].IsError {
			return false
		}
		return results[0].Kind == parser.KindStruct
	default:
		return false
	}
}

func (l *loader) loadRespondHTML(results ActionResults) bool {
	n := len(results)
	if n == 1 || n == 2 {
		if results[0].Named || results[0].Type != "string" {
			return false
		}
		if n == 2 && !results[1].IsError {
			return false
		}
		return true
	}
	return false
}

func (l *loader) loadProvider(controller *Controller, method *parser.Function) *di.Provider {
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
		Name:    fnName,
		Target:  l.module.Import("bud", "controller"),
		Imports: l.imports,
		Hoist:   true,
		Results: []di.Dependency{
			&di.Type{
				Import: importPath,
				Type:   recv.Type().String(),
			},
			&di.Error{},
		},
		Params: []*di.Param{
			{Import: "context", Type: "Context", Hoist: true},
			{Import: "net/http", Type: "*Request"},
			{Import: "net/http", Type: "ResponseWriter"},
		},
		Aliases: di.Aliases{},
	})
	if err != nil {
		l.Bail(err)
	}
	// Add generated imports
	for _, imp := range provider.Imports {
		l.imports.AddNamed(imp.Name, imp.Path)
	}
	// Add the context to the provider set
	l.providers.Add(provider)
	return provider
}

func newProviderSet() *providerSet {
	return &providerSet{map[string]*di.Provider{}}
}

type providerSet struct {
	providerMap map[string]*di.Provider
}

func (c *providerSet) Add(provider *di.Provider) {
	c.providerMap[provider.Name] = provider
}

func (c *providerSet) List() (providers []*di.Provider) {
	for _, provider := range c.providerMap {
		providers = append(providers, provider)
	}
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Name < providers[j].Name
	})
	return providers
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
