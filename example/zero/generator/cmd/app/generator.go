package app

import (
	"fmt"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/gotemplate"
	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/runtime/generator"
)

func New(module *gomod.Module) *Generator {
	return &Generator{module}
}

type Generator struct {
	module *gomod.Module
}

const appTemplate = `package main

{{- if $.Imports }}

import (
	{{- range $import := $.Imports }}
	{{$import.Name}} "{{$import.Path}}"
	{{- end }}
)
{{- end }}

func main() {
	app.Main()
}
`

var appGen = gotemplate.MustParse("app.gotext", appTemplate)

func (g *Generator) GenerateFile(fsys generator.FS, file *generator.File) error {
	type State struct {
		Imports []*imports.Import
	}
	imset := imports.New()
	imset.Add(g.module.Import("generator/cmd/app"))
	code, err := appGen.Generate(State{
		Imports: imset.List(),
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}

func Main() {
	fmt.Println("Hello, from app!!")
}
