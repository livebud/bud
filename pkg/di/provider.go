package di

import (
	"fmt"
	"sort"
	"strings"

	"gitlab.com/mnm/bud/internal/imports"
)

// Provider is the result of generating. Provider can generate functions or
// files or be used for it's template variables.
type Provider struct {
	Name      string            // Name of the function
	Target    string            // Target import path
	Imports   []*imports.Import // Imports needed
	Externals []*External       // External variables
	Code      string            // Body of the generated code
	Results   []*Variable       // Return variables
}

// Function wraps the body code in a function
func (p *Provider) Function() string {
	c := new(strings.Builder)
	params := p.Params()
	var resultTypes []string
	var resultNames []string
	for _, field := range p.Results {
		resultTypes = append(resultTypes, field.Type)
		resultNames = append(resultNames, field.Name)
	}
	resultList := strings.Join(resultTypes, ", ")
	if len(resultTypes) > 1 {
		resultList = "(" + resultList + ")"
	}
	fmt.Fprintf(c, "func %s(%s) %s {\n", p.Name, params, resultList)
	fmt.Fprintf(c, "\t%s", strings.Join(strings.Split(p.Code, "\n"), "\n\t"))
	fmt.Fprintf(c, "return %s\n", strings.Join(resultNames, ", "))
	fmt.Fprintf(c, "}\n")
	return c.String()
}

// Sort the variables by name so the order is always consistent.
func sortByName(externals []*External) []*External {
	sort.Slice(externals, func(i, j int) bool {
		return externals[i].Name < externals[j].Name
	})
	return externals
}

func (p *Provider) Params() (params Params) {
	for _, external := range sortByName(p.Externals) {
		params = append(params, external.Name+" "+external.Type)
	}
	return params
}

type Params []string

func (params Params) String() string {
	return strings.Join(params, ", ")
}

// File wraps the body code in a file
func (p *Provider) File() string {
	c := new(strings.Builder)
	body := p.Function()
	c.WriteString(`package ` + imports.AssumedName(p.Target) + "\n\n")
	c.WriteString("// GENERATED. DO NOT EDIT.\n\n")
	c.WriteString("import (\n")
	for _, im := range p.Imports {
		c.WriteString("\t" + im.Name + ` "` + im.Path + `"` + "\n")
	}
	c.WriteString(")\n\n")
	c.WriteString(body)
	c.WriteString("\n")
	return c.String()
}
