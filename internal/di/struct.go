package di

import (
	"fmt"
	"strings"

	"gitlab.com/mnm/bud/go/is"
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
		if field.Private() {
			continue
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
		decl.Fields = append(decl.Fields, &Field{
			Import: importPath,
			Name:   field.Name(),
			Type:   t.String(),
		})
	}
	return decl, nil
}

type Struct struct {
	Import   string
	Name     string
	NeedsRef bool
	Fields   []*Field
}

var _ Declaration = (*Struct)(nil)

func (s *Struct) ID() string {
	return `"` + s.Import + `".` + s.Name
}

func (s *Struct) Dependencies() (deps []*Dependency) {
	for _, field := range s.Fields {
		deps = append(deps, &Dependency{
			Import: field.Import,
			Type:   field.Type,
		})
	}
	return deps
}

func (s *Struct) Generate(gen *Generator, inputs []*Variable) (outputs []*Variable) {
	var params []string
	for i, input := range inputs {
		key := s.Fields[i].Name
		params = append(params, key+": "+input.Name)
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
