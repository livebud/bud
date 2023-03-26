package middleware

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
	gen.GenerateFile("bud/pkg/web/middleware/middleware.go", g.generateFile)
}

const template = `package middleware

// Code generated by bud; DO NOT EDIT.

{{- if $.Imports }}

import (
	{{- range $import := $.Imports }}
	{{$import.Name}} "{{$import.Path}}"
	{{- end }}
)
{{- end }}

type Interface = middleware.Middleware
type Stack = middleware.Stack

func New(
	csrf *csrf.Middleware,
	session *session.Middleware,
) *Middleware {
	return &Middleware{
		csrf,
		session,
		Stack{
			csrf,
			session,
		},
	}
}

type Middleware struct {
	CSRF    *csrf.Middleware
	Session *session.Middleware
	stack   Stack
}

var _ Interface = (*Middleware)(nil)

func (m *Middleware) Middleware(next http.Handler) http.Handler {
	return m.stack.Middleware(next)
}
`

var gen = gotemplate.MustParse("middleware.gotext", template)

type State struct {
	Imports     []*imports.Import
	Middlewares []*Middleware
}

type Middleware struct {
}

func (g *Generator) generateFile(fsys generator.FS, file *generator.File) error {
	imset := imports.New()
	imset.AddStd("net/http")
	imset.AddNamed("middleware", "github.com/livebud/bud/package/middleware")
	// TODO: generate these
	imset.AddNamed("csrf", g.module.Import("middleware/csrf"))
	imset.AddNamed("session", g.module.Import("middleware/session"))
	code, err := gen.Generate(&State{
		Imports: imset.List(),
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}

//////////////////////////////////
// RUNTIME
//////////////////////////////////
