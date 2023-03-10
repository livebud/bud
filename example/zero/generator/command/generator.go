package command

import (
	"fmt"

	"github.com/livebud/bud/package/gotemplate"
	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/runtime/generator"
)

type Generator struct {
}

const template = `package command

{{- if $.Imports }}

import (
	{{- range $import := $.Imports }}
	{{$import.Name}} "{{$import.Path}}"
	{{- end }}
)
{{- end }}

func New() *Command {
	fmt.Println("new command!")
	return &Command{}
}

type Command struct {
}
`

var gen = gotemplate.MustParse("command.gotext", template)

func (g *Generator) GenerateDir(fsys generator.FS, dir *generator.Dir) error {
	dir.GenerateFile("internal/command/command.go", g.generateFile)
	return nil
}

func (g *Generator) generateFile(fsys generator.FS, file *generator.File) error {
	fmt.Println("generating ./generator/command/!")
	type State struct {
		Imports []*imports.Import
	}
	imset := imports.New()
	imset.AddStd("fmt")
	code, err := gen.Generate(State{
		Imports: imset.List(),
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
