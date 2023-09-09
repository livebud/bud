package main

import (
	command "github.com/livebud/bud/example/zero/bud/internal/command"
	env "github.com/livebud/bud/example/zero/bud/internal/env"
	transpiler "github.com/livebud/bud/example/zero/bud/pkg/transpiler"
	viewer "github.com/livebud/bud/example/zero/bud/pkg/viewer"
	web1 "github.com/livebud/bud/example/zero/bud/pkg/web"
	controller "github.com/livebud/bud/example/zero/bud/pkg/web/controller"
	middleware "github.com/livebud/bud/example/zero/bud/pkg/web/middleware"
	view "github.com/livebud/bud/example/zero/bud/pkg/web/view"
	web2 "github.com/livebud/bud/example/zero/command/web"
	posts "github.com/livebud/bud/example/zero/controller/posts"
	sessions "github.com/livebud/bud/example/zero/controller/sessions"
	users "github.com/livebud/bud/example/zero/controller/users"
	app "github.com/livebud/bud/example/zero/generator/app"
	csrf "github.com/livebud/bud/example/zero/middleware/csrf"
	wraprw "github.com/livebud/bud/example/zero/middleware/wraprw"
	mw "github.com/livebud/bud/example/zero/mw"
	goldmark "github.com/livebud/bud/example/zero/transpiler/goldmark"
	tailwind "github.com/livebud/bud/example/zero/transpiler/tailwind"
	gohtml "github.com/livebud/bud/example/zero/viewer/gohtml"
	web "github.com/livebud/bud/example/zero/web"
	gomod "github.com/livebud/bud/package/gomod"
	log "github.com/livebud/bud/package/log"
)

func main() {
	app.Main(loadCLI)
}

func loadCLI(gomodModule *gomod.Module, logLog log.Log) (*command.CLI, error) {
	envEnv, err := env.Load()
	if err != nil {
		return nil, err
	}
	tailwindTranspiler := &tailwind.Transpiler{Log: logLog}
	goldmarkTranspiler := &goldmark.Transpiler{Log: logLog}
	transpilerTranspiler := transpiler.Load(tailwindTranspiler, goldmarkTranspiler)
	gohtmlViewer := gohtml.New(logLog, transpilerTranspiler)
	viewerViewer := viewer.New(transpilerTranspiler, gohtmlViewer)
	viewView := view.New(gomodModule, viewerViewer)
	postsController := &posts.Controller{}
	sessionsController := &sessions.Controller{}
	usersController := &users.Controller{Env: envEnv}
	controllerController := controller.New(viewView, postsController, sessionsController, usersController)
	csrfMiddleware := &csrf.Middleware{Env: envEnv}
	wraprwMiddleware := &wraprw.Middleware{}
	middlewareMiddleware := middleware.New(csrfMiddleware, wraprwMiddleware)
	mwMiddleware := &mw.Middleware{Env: envEnv}
	webWeb := &web.Web{Env: envEnv, Controller: controllerController, Middleware: middlewareMiddleware, MW: mwMiddleware, View: viewView}
	web1Server := web1.New(webWeb)
	web2Command := web2.New(envEnv, logLog, web1Server)
	commandCLI := command.New(logLog, web2Command)
	return commandCLI, err
}
