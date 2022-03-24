package generator

import (
	overlay "gitlab.com/mnm/bud/package/overlay"
	action "gitlab.com/mnm/bud/runtime/generator/action"
	command "gitlab.com/mnm/bud/runtime/generator/command"
	mainfile "gitlab.com/mnm/bud/runtime/generator/mainfile"
	program "gitlab.com/mnm/bud/runtime/generator/program"
	public "gitlab.com/mnm/bud/runtime/generator/public"
	transform "gitlab.com/mnm/bud/runtime/generator/transform"
	view "gitlab.com/mnm/bud/runtime/generator/view"
	web "gitlab.com/mnm/bud/runtime/generator/web"
	fs "io/fs"
)

func New(
	overlay *overlay.FileSystem,
	main *mainfile.Main,
	program *program.Program,
	command *command.Generator,
	web *web.Generator,
	public *public.Generator,
	action *action.Generator,
	view *view.Generator,
	transform *transform.Generator,
) FS {
	overlay.FileGenerator("bud/.app/main.go", main)
	overlay.FileGenerator("bud/.app/program/program.go", program)
	overlay.FileGenerator("bud/.app/command/command.go", command)
	overlay.FileGenerator("bud/.app/web/web.go", web)
	overlay.FileGenerator("bud/.app/public/public.go", public)
	overlay.FileGenerator("bud/.app/action/action.go", action)
	overlay.FileGenerator("bud/.app/view/view.go", view)
	overlay.FileGenerator("bud/.app/transform/transform.go", transform)
	return overlay
}

type FS = fs.FS
