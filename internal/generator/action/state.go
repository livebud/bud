package action

import (
	"strconv"
	"strings"

	"github.com/matthewmueller/gotext"
	"gitlab.com/mnm/bud/internal/imports"
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

// Action is the target action state
type Action struct {
	Name         string
	Pascal       string
	Camel        string
	Short        string
	View         *View
	Key          string
	Path         string
	Method       string
	Context      *Context
	Inputs       []*ActionInput
	Outputs      ActionOutputs
	ResponseJSON bool
}

// View struct
type View struct {
	Path string
}

// ActionInput struct
type ActionInput struct {
	Name     string
	Pascal   string
	Snake    string
	Type     string
	Variable string
	JSON     string
}

// ActionOutputs fn
type ActionOutputs []*ActionOutput

// Set helper
func (outputs ActionOutputs) Set() string {
	if len(outputs) == 0 {
		return ""
	}
	variables := make([]string, len(outputs))
	for i, output := range outputs {
		variables[i] = output.Variable
	}
	results := strings.Join(variables, ", ")
	// Tack on the operator
	return results + " := "
}

// Result expression if there is one
func (outputs ActionOutputs) Result() string {
	var results ActionOutputs
	for _, output := range outputs {
		if output.IsError {
			continue
		}
		results = append(results, output)
	}
	// Return nothing
	if len(results) == 0 {
		return ""
	}
	if len(results) == 1 {
		return results[0].Variable
	}
	if results.isArray() {
		out := "[]interface{}{"
		for _, result := range results {
			out += result.Variable
			out += ","
		}
		out += "}"
		return out
	}
	if results.isObject() {
		out := "map[string]interface{}{"
		for _, result := range results {
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

func (outputs ActionOutputs) isArray() bool {
	for _, output := range outputs {
		if !output.Named {
			return true
		}
	}
	return false
}

func (outputs ActionOutputs) isObject() bool {
	for _, output := range outputs {
		if output.Pascal == "" {
			return false
		}
	}
	return true
}

// Error expression if there is one
func (outputs ActionOutputs) Error() string {
	for _, output := range outputs {
		if output.IsError {
			return output.Variable
		}
	}
	return ""
}

// ActionOutput struct
type ActionOutput struct {
	Name     string
	Pascal   string
	Named    bool
	Snake    string
	Type     string
	Variable string
	IsError  bool
	Fields   []*ActionOutputField
	Methods  []*ActionOutputMethod
}

// ActionOutputField struct
type ActionOutputField struct {
}

// ActionOutputMethod struct
type ActionOutputMethod struct {
}

// Context is the target context state
type Context struct {
	Function string // Name of the function
	Code     string // Function code
	Inputs   []*ContextInput
	Outputs  ContextOutputs
}

// ContextInput struct
type ContextInput struct {
	Name     string
	Variable string
	Hoisted  bool
	Type     string
}

// ContextOutput struct
type ContextOutput struct {
	Variable string
}

// ContextOutputs is a list of context results
type ContextOutputs []*ContextOutput

// List joins the outputs with a comma
func (outputs ContextOutputs) List() string {
	var outs []string
	for _, output := range outputs {
		outs = append(outs, output.Variable)
	}
	return strings.Join(outs, ", ")
}

// Result returns the result variable if there is one
func (outputs ContextOutputs) Result() string {
	if len(outputs) > 0 {
		return gotext.Camel(outputs[0].Variable)
	}
	return ""
}

// Error returns the error variable if there is one
func (outputs ContextOutputs) Error() string {
	if len(outputs) > 1 {
		return gotext.Camel(outputs[1].Variable)
	}
	return ""
}
