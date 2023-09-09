package generator

import (
	app "github.com/livebud/bud/example/zero/generator/app"
	command "github.com/livebud/bud/example/zero/generator/command"
	controller "github.com/livebud/bud/example/zero/generator/controller"
	env "github.com/livebud/bud/example/zero/generator/env"
	middleware "github.com/livebud/bud/example/zero/generator/middleware"
	session "github.com/livebud/bud/example/zero/generator/session"
	view "github.com/livebud/bud/example/zero/generator/view"
	web "github.com/livebud/bud/example/zero/generator/web"
	log "github.com/livebud/bud/package/log"
	generator "github.com/livebud/bud/runtime/generator"
)

func NewGenerator(
	genfs generator.FileSystem,
	log log.Log,
	app *app.Generator,
	command *command.Generator,
	web *web.Generator,
	controller *controller.Generator,
	view *view.Generator,
	env *env.Generator,
	middleware *middleware.Generator,
	session *session.Generator,
) *Generator {
	return generator.NewGenerator(
		genfs,
		log,
		app,
		command,
		web,
		controller,
		view,
		env,
		middleware,
		session,
	)
}

type Generator = generator.Generator
