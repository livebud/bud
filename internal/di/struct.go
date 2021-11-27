package di

import (
	"fmt"
	"strings"

	"gitlab.com/mnm/bud/go/is"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/parser"
)

// Check to see if the struct initializes the dependency.
//
// Given the following dependency: Web, tryStruct will match on the
// following functions:
//
//   type Web struct { ... }
//
func tryStruct(stct *parser.Struct, dep *Dependency) (*Struct, error) {
	if stct.Private() {
		return nil, ErrNoMatch
	}
	if strings.TrimPrefix(dep.Type, "*") != stct.Name() {
		return nil, ErrNoMatch
	}
	importPath, err := stct.File().Import()
	if err != nil {
		return nil, err
	}
	decl := &Struct{
		Import:   importPath,
		Name:     stct.Name(),
		NeedsRef: strings.HasPrefix(dep.Type, "*"),
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
		modFile, err := def.Package().Modfile()
		if err != nil {
			return nil, err
		}
		decl.Fields = append(decl.Fields, &StructField{
			ModFile: modFile,
			Import:  importPath,
			Name:    field.Name(),
			Type:    t.String(),
		})
	}
	return decl, nil
}

type Struct struct {
	Import   string
	Name     string
	NeedsRef bool
	Fields   []*StructField
}

var _ Declaration = (*Struct)(nil)

type StructField struct {
	ModFile *mod.File // Module file that contains this struct
	Import  string    // Import path
	Type    string    // Field type
	Name    string    // Field or parameter name
}

func (s *Struct) ID() string {
	return `"` + s.Import + `".` + s.Name
}

func (s *Struct) Dependencies() (deps []*Dependency) {
	for _, field := range s.Fields {
		deps = append(deps, &Dependency{
			ModFile: field.ModFile,
			Import:  field.Import,
			Type:    field.Type,
		})
	}
	return deps
}

func (s *Struct) Generate(gen *Generator, inputs []*Variable) (outputs []*Variable) {
	var params []string
	for i, input := range inputs {
		field := s.Fields[i]
		value := maybePrefixField(field, input)
		params = append(params, field.Name+": "+value)
	}
	identifier := gen.Identifier(s.Import, s.Name)
	result := gen.Variable(s.Import, s.Name)
	output := &Variable{
		Import: s.Import,
		Name:   result,
		Type:   s.Name,
	}
	if s.NeedsRef {
		identifier = "&" + identifier
		output.Type = "*" + output.Type
	}
	fmt.Fprintf(gen.Code, "%s := %s{%s}\n", result, identifier, strings.Join(params, ", "))
	return append(outputs, output)
}

// maybePrefix allows us to reference and derefence values during generate so
// the result type doesn't need to be exact.
func maybePrefixField(field *StructField, input *Variable) string {
	if field.Type == input.Type {
		return input.Name
	}
	// Want *T, got T. Need to reference.
	if strings.HasPrefix(field.Type, "*") && !strings.HasPrefix(input.Type, "*") {
		return "&" + input.Name
	}
	// Want T, got*T. Need to dereference.
	if !strings.HasPrefix(field.Type, "*") && strings.HasPrefix(input.Type, "*") {
		return "*" + input.Name
	}
	// We really shouldn't reach here.
	return input.Name
}
