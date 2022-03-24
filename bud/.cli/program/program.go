package program

import (
	context "context"
	errors "errors"
	command "gitlab.com/mnm/bud/bud/.cli/command"
	generator "gitlab.com/mnm/bud/bud/.cli/generator"
	di "gitlab.com/mnm/bud/package/di"
	gomod "gitlab.com/mnm/bud/package/gomod"
	console "gitlab.com/mnm/bud/package/log/console"
	overlay "gitlab.com/mnm/bud/package/overlay"
	parser "gitlab.com/mnm/bud/package/parser"
	action "gitlab.com/mnm/bud/runtime/generator/action"
	command1 "gitlab.com/mnm/bud/runtime/generator/command"
	mainfile "gitlab.com/mnm/bud/runtime/generator/mainfile"
	program "gitlab.com/mnm/bud/runtime/generator/program"
	public "gitlab.com/mnm/bud/runtime/generator/public"
	transform "gitlab.com/mnm/bud/runtime/generator/transform"
	view "gitlab.com/mnm/bud/runtime/generator/view"
	web "gitlab.com/mnm/bud/runtime/generator/web"
)

func Run(ctx context.Context, args ...string) int {
	if err := run(ctx, args...); err != nil {
		if errors.Is(err, context.Canceled) {
			// Unfortunately interrupts like SIGINT trigger a non-zero exit code,
			// regardless of whether you do os.Exit(0) or not. We're going to use exit
			// code 3 to distinguish between non-zero exit codes so "bud run" can know
			// that we exited cleanly on an interrupt.
			return 3
		}
		console.Error(err.Error())
		return 1
	}
	return 0
}

func run(ctx context.Context, args ...string) error {
	program, err := Load()
	if err != nil {
		return err
	}
	return program.Run(ctx, args...)
}

func Load() (*Program, error) {
	gomodModule, err := gomod.Find(".")
	if err != nil {
		return nil, err
	}
	cli, err := loadCLI(gomodModule)
	if err != nil {
		return nil, err
	}
	return &Program{cli}, nil
}

type Program struct {
	cli *command.CLI
}

func (p *Program) Run(ctx context.Context, args ...string) error {
	return p.cli.Parse(ctx, args...)
}

func loadCLI(gomodModule *gomod.Module) (*command.CLI, error) {
	overlayFileSystem, err := overlay.Load(gomodModule)
	if err != nil {
		return nil, err
	}
	mainfileMain := mainfile.New(overlayFileSystem, gomodModule)
	parserParser := parser.New(overlayFileSystem, gomodModule)
	diInjector := di.New(overlayFileSystem, gomodModule, parserParser)
	programProgram := &program.Program{FS: overlayFileSystem, Module: gomodModule, Injector: diInjector}
	command1Generator := &command1.Generator{FS: overlayFileSystem, Module: gomodModule, Parser: parserParser}
	webGenerator := &web.Generator{FS: overlayFileSystem, Module: gomodModule, Parser: parserParser}
	publicGenerator := public.New(overlayFileSystem, gomodModule)
	actionGenerator := &action.Generator{FS: overlayFileSystem, Injector: diInjector, Module: gomodModule, Parser: parserParser}
	viewGenerator := &view.Generator{FS: overlayFileSystem, Module: gomodModule}
	transformGenerator := &transform.Generator{FS: overlayFileSystem, Module: gomodModule}
	generatorFS := generator.New(overlayFileSystem, mainfileMain, programProgram, command1Generator, webGenerator, publicGenerator, actionGenerator, viewGenerator, transformGenerator)
	commandCLI := command.New(generatorFS, gomodModule)
	return commandCLI, err
}
