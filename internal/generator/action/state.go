package action

import (
	"strconv"
	"strings"

	"github.com/matthewmueller/gotext"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/internal/parser"
)

type State struct {
	Imports    []*imports.Import
	Controller *Controller
	Contexts   []*Context
}

// Controller is the target controller state
type Controller struct {
	Name        string
	Pascal      string
	Path        string
	Actions     []*Action
	Controllers []*Controller
}

func (c *Controller) Last() Name {
	names := strings.Split(c.Name, " ")
	return Name(names[len(names)-1])
}

type Name string

func (n Name) Pascal() string {
	return gotext.Pascal(string(n))
}

// Action is the target action state
type Action struct {
	Name         string
	Pascal       string
	Camel        string
	Short        string
	View         *View
	Key          string
	Path         string
	Redirect     string
	Method       string
	Context      *Context
	Params       []*ActionParam
	Input        string
	Results      ActionResults
	ResponseJSON bool
}

// View struct
type View struct {
	Path string
}

// ActionParam struct
type ActionParam struct {
	Name     string
	Pascal   string
	Snake    string
	Type     string
	Kind     string
	Variable string
	Tag      string
}

// ActionResults fn
type ActionResults []*ActionResult

// Set helper
func (results ActionResults) Set() string {
	if len(results) == 0 {
		return ""
	}
	variables := make([]string, len(results))
	for i, output := range results {
		variables[i] = output.Variable
	}
	resultString := strings.Join(variables, ", ")
	// Tack on the operator
	return resultString + " := "
}

// Result expression if there is one
func (results ActionResults) Result() string {
	var list ActionResults
	for _, result := range results {
		if result.IsError {
			continue
		}
		list = append(list, result)
	}
	// Return nothing
	if len(list) == 0 {
		return ""
	}
	if len(list) == 1 {
		return list[0].Variable
	}
	if list.isArray() {
		out := "[]interface{}{"
		for _, result := range list {
			out += result.Variable
			out += ","
		}
		out += "}"
		return out
	}
	if list.isObject() {
		out := "map[string]interface{}{"
		for _, result := range list {
			out += strconv.Quote(result.Snake)
			out += ": "
			out += result.Variable
			out += ","
		}
		out += "}"
		return out
	}
	return ""
}

func (results ActionResults) isArray() bool {
	for _, result := range results {
		if !result.Named {
			return true
		}
	}
	return false
}

func (results ActionResults) isObject() bool {
	for _, result := range results {
		if result.Pascal == "" {
			return false
		}
	}
	return true
}

// Error expression if there is one
func (results ActionResults) Error() string {
	for _, result := range results {
		if result.IsError {
			return result.Variable
		}
	}
	return ""
}

// ActionResult struct
type ActionResult struct {
	Name     string
	Pascal   string
	Named    bool
	Snake    string
	Type     string
	Kind     parser.Kind
	Variable string
	IsError  bool
	Fields   []*ActionResultField
	Methods  []*ActionResultMethod
}

// ActionResultField struct
type ActionResultField struct {
	Name string
	Type string
	Tag  string
}

// ActionResultMethod struct
type ActionResultMethod struct {
}

// Context is the target context state
type Context struct {
	Function string // Name of the function
	Code     string // Function code
	Fields   []*ContextField
	Results  ContextResults
}

// ContextField struct
type ContextField struct {
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
