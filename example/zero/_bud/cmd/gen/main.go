package main

import (
	generator "github.com/livebud/bud/example/zero/bud/internal/generator"
	transpiler "github.com/livebud/bud/example/zero/bud/pkg/transpiler"
	viewer "github.com/livebud/bud/example/zero/bud/pkg/viewer"
	app "github.com/livebud/bud/example/zero/generator/app"
	command "github.com/livebud/bud/example/zero/generator/command"
	controller "github.com/livebud/bud/example/zero/generator/controller"
	env "github.com/livebud/bud/example/zero/generator/env"
	middleware "github.com/livebud/bud/example/zero/generator/middleware"
	session "github.com/livebud/bud/example/zero/generator/session"
	view "github.com/livebud/bud/example/zero/generator/view"
	web "github.com/livebud/bud/example/zero/generator/web"
	goldmark "github.com/livebud/bud/example/zero/transpiler/goldmark"
	tailwind "github.com/livebud/bud/example/zero/transpiler/tailwind"
	gohtml "github.com/livebud/bud/example/zero/viewer/gohtml"
	framework "github.com/livebud/bud/framework"
	genfs "github.com/livebud/bud/package/genfs"
	gomod "github.com/livebud/bud/package/gomod"
	log "github.com/livebud/bud/package/log"
	gen "github.com/livebud/bud/runtime/gen"
)

func main() {
	gen.Main(loadGenerator)
}

func loadGenerator(frameworkFlag *framework.Flag, genfsFileSystem genfs.FileSystem, gomodModule *gomod.Module, logLog log.Log) (*generator.Generator, error) {
	genParser := gen.ProvideParser(genfsFileSystem, gomodModule)
	genInjector := gen.ProvideInjector(genfsFileSystem, logLog, gomodModule, genParser)
	appGenerator := app.New(frameworkFlag, genInjector, gomodModule)
	commandGenerator := command.NewGenerator(gomodModule)
	webGenerator := web.New(gomodModule)
	controllerGenerator := controller.New(genInjector, gomodModule)
	tailwindTranspiler := &tailwind.Transpiler{Log: logLog}
	goldmarkTranspiler := &goldmark.Transpiler{Log: logLog}
	transpilerTranspiler := transpiler.Load(tailwindTranspiler, goldmarkTranspiler)
	gohtmlViewer := gohtml.New(logLog, transpilerTranspiler)
	viewerViewer := viewer.New(transpilerTranspiler, gohtmlViewer)
	viewGenerator := view.New(frameworkFlag, gomodModule, viewerViewer)
	envGenerator := env.New(gomodModule)
	middlewareGenerator := middleware.New(gomodModule)
	sessionGenerator := session.New(gomodModule)
	generatorGenerator := generator.NewGenerator(genfsFileSystem, logLog, appGenerator, commandGenerator, webGenerator, controllerGenerator, viewGenerator, envGenerator, middlewareGenerator, sessionGenerator)
	return generatorGenerator, nil
}
