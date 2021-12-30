package action

import (
	"strconv"
	"strings"

	"github.com/matthewmueller/gotext"
	"github.com/matthewmueller/text"
	"gitlab.com/mnm/bud/internal/imports"
)

// State is the target controller state
type State struct {
	Imports    []*imports.Import
	Controller *Controller
}

type Controller struct {
	Name        string
	Route       string
	Import      *imports.Import
	Actions     []*Action
	Context     *Context
	Controllers []*Controller
}

func (c *Controller) Pascal() string {
	return gotext.Pascal(c.Name)
}

// Action is the target action state
type Action struct {
	Name         string
	Key          string
	Route        string
	View         bool
	Method       string
	Context      *Context
	Params       ActionParams
	Results      ActionResults
	ResponseJSON bool
}

func (a *Action) Pascal() string {
	return gotext.Pascal(a.Name)
}

func (a *Action) Short() string {
	return gotext.Short(a.Name)
}

// type Function struct {
// 	Params  []*ActionParam
// 	Results ActionResults
// }

// // View struct
// type View struct {
// 	Path string
// }

type ActionParams []*ActionParam

func (params ActionParams) Publics() (publics ActionParams) {
	for _, param := range params {
		if param.Private() {
			continue
		}
		publics = append(publics, param)
	}
	return publics
}

// ActionParam struct
type ActionParam struct {
	Name     string
	Type     string
	Import   *imports.Import
	Kind     string
	Optional bool
}

func (a *ActionParam) Pascal() string {
	return gotext.Pascal(a.Name)
}

func (a *ActionParam) Variable() string {
	return gotext.Camel(a.Type)
}

func (a *ActionParam) Tag() string {
	w := new(strings.Builder)

	// json tag
	w.WriteString(`json:"`)
	if a.Name == "" {
		w.WriteString(`-`)
	} else {
		w.WriteString(text.Snake(a.Name))
	}
	w.WriteString(`,omitempty"`)

	// validate tag
	w.WriteString(` validate:"`)
	if !a.Optional {
		w.WriteString(`required`)
	}
	w.WriteString(`"`)

	return w.String()
}

// func (a *ActionParam) isContext() bool {
// 	return a.Import != nil &&
// 		a.Import.Path == "context" &&
// 		a.Type == "Context"
// }

func (a *ActionParam) isContext() bool {
	return a.Type == "context.Context"
}

func (a *ActionParam) Private() bool {
	return a.isContext()
}

func (a *ActionParam) Public() bool {
	return !a.Private()
}

// ActionResults fn
type ActionResults []*ActionResult

// Set helper
func (outputs ActionResults) List() string {
	if len(outputs) == 0 {
		return ""
	}
	variables := make([]string, len(outputs))
	for i, output := range outputs {
		variables[i] = output.Variable()
	}
	return strings.Join(variables, ", ")
}

// Result expression if there is one
func (outputs ActionResults) Result() string {
	var results ActionResults
	for _, output := range outputs {
		if output.IsError() {
			continue
		}
		results = append(results, output)
	}
	// Return nothing
	if len(results) == 0 {
		return ""
	}
	if len(results) == 1 {
		return results[0].Variable()
	}
	if results.isArray() {
		out := "[]interface{}{"
		for _, result := range results {
			out += result.Variable()
			out += ","
		}
		out += "}"
		return out
	}
	if results.isObject() {
		out := "map[string]interface{}{"
		for _, result := range results {
			out += strconv.Quote(result.Snake())
			out += ": "
			out += result.Variable()
			out += ","
		}
		out += "}"
		return out
	}
	return ""
}

func (outputs ActionResults) isArray() bool {
	for _, output := range outputs {
		if !output.Named() {
			return true
		}
	}
	return false
}

func (outputs ActionResults) isObject() bool {
	for _, output := range outputs {
		if output.Pascal() == "" {
			return false
		}
	}
	return true
}

// Error expression if there is one
func (outputs ActionResults) Error() string {
	for _, output := range outputs {
		if output.IsError() {
			return output.Variable()
		}
	}
	return ""
}

// ActionResult struct
type ActionResult struct {
	Name    string
	Type    string
	Kind    string
	Fields  []*ActionResultField
	Methods []*ActionResultMethod
}

func (a *ActionResult) Pascal() string {
	return gotext.Pascal(a.Name)
}

func (a *ActionResult) Snake() string {
	return text.Snake(a.Name)
}

func (a *ActionResult) Variable() string {
	return gotext.Camel(a.Type)
}

func (a *ActionResult) Named() bool {
	return a.Name != ""
}

func (a *ActionResult) IsError() bool {
	return a.Type == "error"
}

// ActionResultField struct
type ActionResultField struct {
}

// ActionResultMethod struct
type ActionResultMethod struct {
}

// Context is the target context state
type Context struct {
	Function string // Name of the function
	Code     string // Function code
	Inputs   []*ContextInput
	Results  ContextResults
}

// ContextInput struct
type ContextInput struct {
	Name     string
	Variable string
	Hoisted  bool
	Type     string
}

// ContextResult struct
type ContextResult struct {
	Variable string
}

// ContextResults is a list of context results
type ContextResults []*ContextResult

// List joins the outputs with a comma
func (outputs ContextResults) List() string {
	var outs []string
	for _, output := range outputs {
		outs = append(outs, output.Variable)
	}
	return strings.Join(outs, ", ")
}

// Result returns the result variable if there is one
func (outputs ContextResults) Result() string {
	if len(outputs) > 0 {
		return gotext.Camel(outputs[0].Variable)
	}
	return ""
}

// Error returns the error variable if there is one
func (outputs ContextResults) Error() string {
	if len(outputs) > 1 {
		return gotext.Camel(outputs[1].Variable)
	}
	return ""
}
