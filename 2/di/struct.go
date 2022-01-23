package di

import (
	"fmt"
	"strings"

	"gitlab.com/mnm/bud/2/mod"
	"gitlab.com/mnm/bud/2/parser"
	"gitlab.com/mnm/bud/go/is"
)

// Struct is a dependency that can be defined in memory. Struct is also a
// declaration that can be referenced and be used to generate initializers.
type Struct struct {
	Import string
	Type   string
	Fields []*StructField
}

var _ Dependency = (*Struct)(nil)
var _ Declaration = (*Struct)(nil)

func (s *Struct) ID() string {
	return `"` + s.Import + `".` + s.Type
}

func (s *Struct) ImportPath() string {
	return s.Import
}

func (s *Struct) TypeName() string {
	return s.Type
}

// Find a declaration that provides this type
func (s *Struct) Find(finder Finder) (Declaration, error) {
	return s, nil
}

func (s *Struct) Dependencies() (deps []Dependency) {
	for _, field := range s.Fields {
		deps = append(deps, field)
	}
	return deps
}

func (s *Struct) Generate(gen Generator, inputs []*Variable) (outputs []*Variable) {
	var params []string
	for i, input := range inputs {
		field := s.Fields[i]
		value := maybePrefixField(field, input)
		params = append(params, field.Name+": "+value)
	}
	identifier := gen.Identifier(s.Import, s.Type)
	result := gen.Variable(s.Import, s.Type)
	output := &Variable{
		Import: s.Import,
		Name:   result,
		Type:   s.Type,
		Kind:   parser.KindStruct,
	}
	if strings.HasPrefix(s.Type, "*") {
		identifier = "&" + identifier
	}
	gen.WriteString(fmt.Sprintf("%s := %s{%s}\n", result, identifier, strings.Join(params, ", ")))
	return append(outputs, output)
}

type StructField struct {
	Name   string
	Import string
	Type   string

	module *mod.Module // Module containing this type
	kind   parser.Kind // Kind of type
}

var _ Dependency = (*StructField)(nil)

func (s *StructField) ID() string {
	return getID(s.Import, s.Type)
}

func (s *StructField) ImportPath() string {
	return s.Import
}

func (s *StructField) TypeName() string {
	return s.Type
}

func (s *StructField) Find(finder Finder) (Declaration, error) {
	return finder.Find(s.module, s)
}

// Check to see if the struct initializes the dependency.
//
// Given the following dependency: Web, tryStruct will match on the
// following functions:
//
//   type Web struct { ... }
//
func tryStruct(stct *parser.Struct, dataType string) (*Struct, error) {
	if stct.Private() {
		return nil, ErrNoMatch
	}
	if strings.TrimPrefix(dataType, "*") != stct.Name() {
		return nil, ErrNoMatch
	}
	importPath, err := stct.File().Import()
	if err != nil {
		return nil, err
	}
	decl := &Struct{
		Import: importPath,
		Type:   dataType,
		// needsRef: strings.HasPrefix(dataType, "*"),
	}
	for _, field := range stct.Fields() {
		// Disallow any private fields. This is restrictive but it makes sure
		// that the struct is usable if we initialize it automatically. If you need
		// to use private fields, use a function.
		if field.Private() {
			return nil, ErrNoMatch
		}
		ft := field.Type()
		// Ensure there are no builtin types (e.g. string) as field types
		if is.Builtin(ft.String()) {
			return nil, ErrNoMatch
		}
		importPath, err := parser.ImportPath(ft)
		if err != nil {
			return nil, err
		}
		t := parser.Unqualify(ft)
		def, err := field.Definition()
		if err != nil {
			return nil, err
		}
		pkg := def.Package()
		module := pkg.Module()
		decl.Fields = append(decl.Fields, &StructField{
			Name:   field.Name(),
			Import: importPath,
			Type:   t.String(),
			kind:   def.Kind(),
			module: module,
		})
	}
	return decl, nil
}

// maybePrefix allows us to reference and derefence values during generate so
// the result type doesn't need to be exact.
func maybePrefixField(field *StructField, input *Variable) string {
	if field.Type == input.Type {
		if isInterface(field.kind) && !isInterface(input.Kind) {
			// Create a pointer to the input when field is an interface type, but the
			// input is not an interface.
			return "&" + input.Name
		}
		return input.Name
	}
	// Want *T, got T. Need to reference.
	if strings.HasPrefix(field.Type, "*") && !strings.HasPrefix(input.Type, "*") {
		return "&" + input.Name
	}
	// Want T, got *T. Need to dereference.
	if !strings.HasPrefix(field.Type, "*") && strings.HasPrefix(input.Type, "*") {
		if isInterface(field.kind) {
			// Don't dereference the type when the field is an interface type
			return input.Name
		}
		return "*" + input.Name
	}
	// We really shouldn't reach here.
	return input.Name
}
