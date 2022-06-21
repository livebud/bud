package bud

import (
	"context"
	"io"

	"github.com/livebud/bud/internal/pubsub"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/framework/app"
	"github.com/livebud/bud/framework/controller"
	"github.com/livebud/bud/framework/public"
	"github.com/livebud/bud/framework/view"
	"github.com/livebud/bud/framework/web"
	"github.com/livebud/bud/package/commander"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/filter"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/package/svelte"
	"github.com/livebud/bud/runtime/transform"
	"github.com/livebud/bud/runtime/view/dom"
	"github.com/livebud/bud/runtime/view/ssr"
)

// Input contains the configuration that gets passed into the commands
type Input struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    []string
	BudLn  socket.Listener // Can be nil
	WebLn  socket.Listener // Can be nil
	Bus    pubsub.Client   // Can be nil
}

func New(in *Input) *Command {
	return &Command{in: in}
}

type Command struct {
	in   *Input
	Dir  string
	Log  string
	Args []string
	Help bool
}

// Run a custom command
// TODO: finish supporting custom commands
// 1. Compile
//   a. Generate generator (later!)
//   	 i. Generate bud/internal/generator
//     ii. Build bud/generator
//     iii. Run bud/generator
//   b. Generate custom command
//     i. Generate bud/internal/command/${name}/
//     ii. Build bud/command/${name}
// 2. Run bud/command/${name}
func (c *Command) Run(ctx context.Context) error {
	return commander.Usage()
}

// Module finds the go.mod file for the application
func Module(dir string) (*gomod.Module, error) {
	return gomod.Find(dir)
}

func Log(stderr io.Writer, logFilter string) (log.Interface, error) {
	handler, err := filter.Load(console.New(stderr), logFilter)
	if err != nil {
		return nil, err
	}
	return log.New(handler), nil
}

func FileSystem(log log.Interface, module *gomod.Module, flag *framework.Flag) (*overlay.FileSystem, error) {
	genfs, err := overlay.Load(log, module)
	if err != nil {
		return nil, err
	}
	parser := parser.New(genfs, module)
	injector := di.New(genfs, log, module, parser)
	vm, err := v8.Load()
	if err != nil {
		return nil, err
	}
	svelteCompiler, err := svelte.Load(vm)
	if err != nil {
		return nil, err
	}
	transforms, err := transform.Load(svelte.NewTransformable(svelteCompiler))
	if err != nil {
		return nil, err
	}
	genfs.FileGenerator("bud/internal/app/main.go", app.New(injector, module, flag))
	genfs.FileGenerator("bud/internal/app/web/web.go", web.New(module, parser))
	genfs.FileGenerator("bud/internal/app/controller/controller.go", controller.New(injector, module, parser))
	genfs.FileGenerator("bud/internal/app/view/view.go", view.New(module, transforms, flag))
	genfs.FileGenerator("bud/internal/app/public/public.go", public.New(flag))
	return genfs, nil
}

func FileServer(log log.Interface, module *gomod.Module, flag *framework.Flag) (*overlay.Server, error) {
	servefs, err := overlay.Serve(log, module)
	if err != nil {
		return nil, err
	}
	vm, err := v8.Load()
	if err != nil {
		return nil, err
	}
	svelteCompiler, err := svelte.Load(vm)
	if err != nil {
		return nil, err
	}
	transforms, err := transform.Load(svelte.NewTransformable(svelteCompiler))
	if err != nil {
		return nil, err
	}
	servefs.FileGenerator("bud/view/_ssr.js", ssr.New(module, transforms.SSR))
	servefs.FileServer("bud/view", dom.New(module, transforms.DOM))
	servefs.FileServer("bud/node_modules", dom.NodeModules(module))
	return servefs, nil
}

// // Generate the app
// func (c *Command) Generate(genfs *overlay.FileSystem, outDir string) error {
// 	return genfs.Sync(outDir)
// }

// func (c *Command) Build(ctx context.Context, module *gomod.Module, mainPath, outPath string) error {
// 	builder := gobuild.New(module)
// 	return builder.Build(ctx, mainPath, outPath)
// }

// func (c *Command) Watch(ctx context.Context, module *gomod.Module, log log.Interface, fn func(isBoot, canHotReload bool) error) error {
// 	// Wrap the function
// 	watchFn := func(paths []string) error {
// 		if err := fn(false, canHotReload(paths)); err != nil {
// 			log.Error(err.Error())
// 		}
// 		return nil
// 	}
// 	// Call the function once manually to boot
// 	if err := fn(true, false); err != nil {
// 		log.Error(err.Error())
// 	}
// 	// Regardless of success, watch for changes
// 	return watcher.Watch(ctx, module.Directory(), watchFn)
// }

// // canHotReload returns true if we can incrementally reload a page
// func canHotReload(paths []string) bool {
// 	for _, path := range paths {
// 		if filepath.Ext(path) == ".go" {
// 			return false
// 		}
// 	}
// 	return true
// }
