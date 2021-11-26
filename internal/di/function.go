package di

import (
	"fmt"
	"strings"

	"gitlab.com/mnm/bud/go/is"
	"gitlab.com/mnm/bud/internal/parser"
)

// Check to see if the function initializes the dependency.
//
// Given the following dependency: *Web, tryFunction will match on the
// following functions:
//
//   func ...(...) *Web
//   func ...(...) Web
//   func ...(...) (*Web, error)
//   func ...(...) (Web, error)
//
func tryFunction(fn *parser.Function, dep *Dependency) (*Function, error) {
	if fn.Private() {
		return nil, ErrNoMatch
	}
	results := fn.Results()
	if len(results) < 1 || len(results) > 2 {
		return nil, ErrNoMatch
	}
	resultType := results[0].Type()
	innerType := parser.Unqualify(resultType).String()
	innerName := strings.TrimPrefix(innerType, "*")
	depName := strings.TrimPrefix(dep.Type, "*")
	if innerName != depName {
		return nil, ErrNoMatch
	}
	typeImport, err := parser.ImportPath(resultType)
	if err != nil {
		return nil, err
	}
	if typeImport != dep.Import {
		return nil, ErrNoMatch
	}
	importPath, err := fn.File().Import()
	if err != nil {
		return nil, err
	}
	function := &Function{
		Import: importPath,
		Name:   fn.Name(),
	}
	for _, param := range fn.Params() {
		pt := param.Type()
		// Ensure there are no builtin types (e.g. string) as parameters
		if is.Builtin(pt.String()) {
			return nil, ErrNoMatch
		}
		importPath, err := parser.ImportPath(pt)
		if err != nil {
			return nil, err
		}
		t := parser.Unqualify(pt)
		def, err := param.Definition()
		if err != nil {
			return nil, err
		}
		modFile, err := def.Package().Modfile()
		if err != nil {
			return nil, err
		}
		function.Params = append(function.Params, &Dependency{
			Import:  importPath,
			Type:    t.String(),
			ModFile: modFile,
		})
	}
	for _, result := range results {
		resultType := result.Type()
		name := result.Name()
		if name == "" {
			name = parser.TypeName(resultType)
		}
		importPath, err := parser.ImportPath(resultType)
		if err != nil {
			return nil, err
		}
		unqualified := parser.Unqualify(resultType)
		function.Results = append(function.Results, &ResultField{
			Import: importPath,
			Name:   name,
			Type:   unqualified.String(),
		})
		continue
	}
	return function, nil
}

type Function struct {
	Import  string
	Name    string
	Params  []*Dependency
	Results []*ResultField
}

type ResultField struct {
	Import string // Import path
	Type   string // Result type
	Name   string // Result name
}

var _ Declaration = (*Function)(nil)

func (fn *Function) ID() string {
	return `"` + fn.Import + `".` + fn.Name
}

// Dependencies are the values that the function depends on to run
func (fn *Function) Dependencies() []*Dependency {
	return fn.Params
}

// Returns true if the 2nd result is an error
func (fn *Function) hasError() bool {
	l := len(fn.Results)
	if l == 0 {
		return false
	}
	last := fn.Results[l-1]
	return last.Type == "error"
}

// maybePrefix allows us to reference and derefence values during generate so
// the result type doesn't need to be exact.
func maybePrefix(param *Dependency, input *Variable) string {
	if param.Type == input.Type {
		return input.Name
	}
	// Want *T, got T. Need to reference.
	if strings.HasPrefix(param.Type, "*") && !strings.HasPrefix(input.Type, "*") {
		return "&" + input.Name
	}
	// Want T, got*T. Need to dereference.
	if !strings.HasPrefix(param.Type, "*") && strings.HasPrefix(input.Type, "*") {
		return "*" + input.Name
	}
	// We really shouldn't reach here.
	return input.Name
}

// Generate a statement that calls this function and returns the results.
func (fn *Function) Generate(gen *Generator, inputs []*Variable) (outputs []*Variable) {
	var params []string
	for i, input := range inputs {
		params = append(params, maybePrefix(fn.Params[i], input))
	}
	identifier := gen.Identifier(fn.Import, fn.Name)
	var results []string
	for _, result := range fn.Results {
		name := gen.Variable(result.Import, result.Type)
		results = append(results, name)
		outputs = append(outputs, &Variable{
			Import: result.Import,
			Type:   result.Type,
			Name:   name,
		})
	}
	fmt.Fprintf(gen.Code, "%s := %s(%s)\n", strings.Join(results, ", "), identifier, strings.Join(params, ", "))
	if fn.hasError() {
		gen.HasError = true
		errvar := outputs[len(outputs)-1]
		fmt.Fprintf(gen.Code, "if %[1]s != nil {\n\treturn nil, %[1]s\n}\n", errvar.Name)
	}
	return outputs
}
