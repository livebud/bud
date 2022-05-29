package di

import (
	"fmt"
	"strings"

	"github.com/livebud/bud/internal/gois"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/parser"
)

// Function is the top-level load function that we generate to provide all the
// dependencies
type Function struct {
	// Name of the function to generate
	Name string
	// Imports to pass through
	Imports *imports.Set
	// Params are the external parameters that are passed in
	Params []Dependency
	// Results are the dependencies that need to be loaded
	Results []Dependency
	// Hoist dependencies that don't depend on externals, turning them into
	// externals. This is to avoid initializing these inner deps every time.
	// Useful for per-request dependency injection.
	Hoist bool
	// Aliases allow you to map one dependency to another. Useful to supporting
	// interfaces as inputs that are mapped to a concrete value.
	Aliases Aliases
	// Target import path where this function will be generated to
	Target string
}

var _ Declaration = (*Function)(nil)

func (fn *Function) ID() string {
	return getID(fn.Target, fn.Name)
}

func (fn *Function) Validate() error {
	if fn.Name == "" {
		return fmt.Errorf("di: function must have a name")
	}
	return nil
}

func (fn *Function) Dependencies() []Dependency {
	return fn.Results
}

// Generate the function declaration that can be called to initialize the
// required dependencies
func (fn *Function) Generate(g Generator, ins []*Variable) (outs []*Variable) {
	return ins
}

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
func tryFunction(fn *parser.Function, importPath, dataType string) (*function, error) {
	if fn.Private() || fn.Receiver() != nil {
		return nil, ErrNoMatch
	}
	results := fn.Results()
	if len(results) < 1 || len(results) > 2 {
		return nil, ErrNoMatch
	}
	resultType := results[0].Type()
	innerType := parser.Unqualify(resultType).String()
	innerName := strings.TrimPrefix(innerType, "*")
	depName := strings.TrimPrefix(dataType, "*")
	if innerName != depName {
		return nil, ErrNoMatch
	}
	typeImport, err := parser.ImportPath(resultType)
	if err != nil {
		return nil, err
	}
	if typeImport != importPath {
		return nil, ErrNoMatch
	}
	fileImportPath, err := fn.File().Import()
	if err != nil {
		return nil, err
	}
	function := &function{
		Import: fileImportPath,
		Name:   fn.Name(),
	}
	for _, param := range fn.Params() {
		pt := param.Type()
		// Ensure there are no builtin types (e.g. string) as parameters
		if gois.Builtin(pt.String()) {
			return nil, ErrNoMatch
		}
		imPath, err := parser.ImportPath(pt)
		if err != nil {
			return nil, err
		}
		def, err := param.Definition()
		if err != nil {
			return nil, fmt.Errorf("di: unable to find definition for param %q.%s in %q.%s . %w", imPath, parser.Unqualify(pt).String(), importPath, dataType, err)
		}
		module := def.Package().Module()
		function.Params = append(function.Params, &Type{
			Import: imPath,
			Type:   parser.Unqualify(pt).String(),
			kind:   def.Kind(),
			module: module,
		})
	}
	for _, result := range results {
		rt := result.Type()
		name := result.Name()
		if name == "" {
			name = parser.TypeName(rt)
		}
		imPath, err := parser.ImportPath(rt)
		if err != nil {
			return nil, err
		}
		def, err := result.Definition()
		if err != nil {
			return nil, fmt.Errorf("di: unable to find definition for result %q.%s in %q.%s . %w", imPath, parser.Unqualify(rt).String(), importPath, dataType, err)
		}
		unqualified := parser.Unqualify(rt)
		function.Results = append(function.Results, &Type{
			Import: importPath,
			Type:   unqualified.String(),
			kind:   def.Kind(),
			name:   name,
		})
		continue
	}
	return function, nil
}

// Function is a declaration that can provide a dependency
type function struct {
	Import  string
	Name    string
	Params  []*Type
	Results []*Type
}

var _ Declaration = (*function)(nil)

func (fn *function) ID() string {
	return `"` + fn.Import + `".` + fn.Name
}

// Dependencies are the values that the funcDecl depends on to run
func (fn *function) Dependencies() (deps []Dependency) {
	for _, param := range fn.Params {
		deps = append(deps, param)
	}
	return deps
}

// Returns true if the 2nd result is an error
func (fn *function) hasError() bool {
	l := len(fn.Results)
	if l == 0 {
		return false
	}
	last := fn.Results[l-1]
	return last.Type == "error"
}

// Generate a caller that can initialize a dependency
func (fn *function) Generate(gen Generator, inputs []*Variable) (outputs []*Variable) {
	var params []string
	for i, input := range inputs {
		params = append(params, maybePrefixParam(fn.Params[i], input))
	}
	identifier := gen.Identifier(fn.Import, fn.Name)
	var results []string
	for _, result := range fn.Results {
		name := gen.Variable(result.Import, result.Type)
		results = append(results, name)
		outputs = append(outputs, &Variable{
			Import: result.Import,
			Type:   result.Type,
			Kind:   result.kind,
			Name:   name,
		})
	}
	gen.WriteString(fmt.Sprintf("%s := %s(%s)\n", strings.Join(results, ", "), identifier, strings.Join(params, ", ")))
	if fn.hasError() {
		// Mark the code as having an error
		gen.MarkError(true)
		errvar := outputs[len(outputs)-1]
		gen.WriteString(fmt.Sprintf("if %[1]s != nil {\n\treturn nil, %[1]s\n}\n", errvar.Name))
	}
	return outputs
}

// maybePrefix allows us to reference and derefence values during generate so
// the result type doesn't need to be exact.
func maybePrefixParam(param *Type, input *Variable) string {
	sameImport := param.Import == input.Import
	sameType := param.Type == input.Type
	// Nothing to change
	if sameImport && sameType {
		return input.Name
	}
	// Want *T, got T. Need to reference.
	if sameImport && strings.HasPrefix(param.Type, "*") && !strings.HasPrefix(input.Type, "*") {
		return "&" + input.Name
	}
	// Want T, got *T. Need to dereference.
	if sameImport && !strings.HasPrefix(param.Type, "*") && strings.HasPrefix(input.Type, "*") {
		if isInterface(param.kind) {
			// Don't dereference the type when the param is an interface type
			return input.Name
		}
		return "*" + input.Name
	}
	// Create a pointer to the input when param is an interface type, but the
	// input is not an interface.
	// TODO: This is hacky because input.Kind can't tell if it's a pointer to a
	// struct vs. a struct, so we need extra logic to tell.
	if !sameImport && isInterface(param.kind) && !isInterface(input.Kind) && !strings.HasPrefix(input.Type, "*") {
		return "&" + input.Name
	}
	// Passing through, not sure what to do.
	return input.Name
}
