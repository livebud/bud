package controller

import (
	"strconv"
	"strings"

	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/parser"
	"github.com/matthewmueller/gotext"
	"github.com/matthewmueller/text"
)

type State struct {
	Imports    []*imports.Import
	Controller *Controller
	Providers  []*di.Provider
}

// Controller is the target controller state
type Controller struct {
	Name        string
	Pascal      string
	JSON        string
	Path        string // Path to controller without action dir
	Route       string
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
	Name        string
	Pascal      string
	Camel       string
	Short       string
	View        *View
	Key         string // Key is an extension-less path
	Route       string // Route to this action
	Redirect    string
	Method      string
	Provider    *di.Provider
	Params      []*ActionParam
	HandlerFunc bool
	Input       string
	Results     ActionResults
	RespondJSON bool
	RespondHTML bool
	PropsKey    string
}

// View struct
type View struct {
	Route string
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

func (ap *ActionParam) IsContext() bool {
	return ap.Type == "context.Context"
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

func (results ActionResults) propsKey() string {
	for _, result := range results {
		if result.IsError {
			continue
		} else if result.Named {
			return result.Name
		}
		dataType := strings.TrimPrefix(result.Type, "[]*")
		if isList(result.Type) {
			// e.g. []*UserStory => userStories
			return gotext.Camel(text.Plural(dataType))
		}
		// e.g. *UserStory => userStory
		return gotext.Camel(dataType)
	}
	return ""
}

func isList(dataType string) bool {
	return strings.HasPrefix(dataType, "[]") ||
		strings.HasPrefix(dataType, "map")
}

func (results ActionResults) ViewResult() string {
	propsKey := results.propsKey()
	out := new(strings.Builder)
	out.WriteString(`map[string]interface{}{`)
	if propsKey != "" {
		out.WriteString(strconv.Quote(propsKey))
		out.WriteString(":")
		out.WriteString(results.Result())
	}
	out.WriteString(`},`)
	return out.String()
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

// Error expression is only return
func (results ActionResults) IsOnlyError() bool {
	return len(results) == 1 && results[0].IsError
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
