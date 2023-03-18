package web

import (
	"net/http"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/gotemplate"
	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/package/middleware"
	"github.com/livebud/bud/runtime/generator"
)

func New(module *gomod.Module) *Generator {
	return &Generator{module}
}

type Generator struct {
	module *gomod.Module
}

func (g *Generator) Extend(gen generator.FileSystem) {
	gen.GenerateFile("bud/pkg/web/web.go", g.generateFile)
}

const template = `package web

{{- if $.Imports }}

import (
	{{- range $import := $.Imports }}
	{{$import.Name}} "{{$import.Path}}"
	{{- end }}
)
{{- end }}

func New(web *web.Web) *Server {
	middleware := web.Middleware(middleware.Stack{})
	router := web.Router(router.New())
	handler := web.Handler(middleware.Middleware(router))
	server := web.Server(&http.Server{
		Handler: handler,
	})
	return &Server{server}
}

type Server struct {
	server *http.Server
}

func (s *Server) Listen(ctx context.Context, addr string) error {
	ln, err := socket.Listen(addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	return s.Serve(ctx, ln)
}

func (s *Server) Serve(ctx context.Context, ln net.Listener) error {
	return s.server.Serve(ln)
}
`

var gen = gotemplate.MustParse("web.gotext", template)

type State struct {
	Imports []*imports.Import
}

func (g *Generator) generateFile(fsys generator.FS, file *generator.File) error {
	imset := imports.New()
	imset.AddStd("context", "net/http", "net")
	// imset.AddNamed("runweb", g.module.Import("generator/web"))
	imset.AddNamed("web", g.module.Import("web"))
	imset.AddNamed("router", "github.com/livebud/bud/package/router")
	imset.AddNamed("socket", "github.com/livebud/bud/package/socket")
	imset.AddNamed("middleware", "github.com/livebud/bud/package/middleware")

	code, err := gen.Generate(&State{
		Imports: imset.List(),
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}

/////////////////////////////
// Runtime
/////////////////////////////

// type Server interface {
// 	Listen(ctx context.Context, addr string) error
// 	Serve(ctx context.Context, ln net.Listener) error
// }

// TODO: move this into *router.Router
type Router interface {
	http.Handler
	middleware.Middleware
	Add(method string, route string, handler http.Handler) error
	Get(route string, handler http.Handler) error
	Post(route string, handler http.Handler) error
	Put(route string, handler http.Handler) error
	Patch(route string, handler http.Handler) error
	Delete(route string, handler http.Handler) error
}

type Middlewares = middleware.Stack
