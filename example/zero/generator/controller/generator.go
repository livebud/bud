package controller

import (
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

func (g *Generator) Extend(gen generator.FileSystem) {
	gen.GenerateFile("bud/pkg/web/controller/controller.go", g.generateFile)
}

const template = `package controller

{{- if $.Imports }}

import (
	{{- range $import := $.Imports }}
	{{ $import.Name }} "{{ $import.Path }}"
	{{- end }}
)
{{- end }}

func New(
	controller *controller.Controller,
) *Controller {
	return &Controller{
		&IndexAction{controller},
	}
}

type Controller struct {
	Index *IndexAction
}

// TODO: use a router.Router interface
func (c *Controller) Mount(r *router.Router) error {
	r.Get("/", c.Index)
	return nil
}

type IndexAction struct {
	controller *controller.Controller
}

func (a *IndexAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := a.controller.Index()
	w.Write([]byte(res))
}
`

var gen = gotemplate.MustParse("controller.gotext", template)

type State struct {
	Imports []*imports.Import
}

func (g *Generator) generateFile(fsys generator.FS, file *generator.File) error {
	imset := imports.New()
	imset.AddStd("net/http")
	imset.AddNamed("router", "github.com/livebud/bud/package/router")
	imset.AddNamed("controller", g.module.Import("controller"))
	code, err := gen.Generate(&State{
		Imports: imset.List(),
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
