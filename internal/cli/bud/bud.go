package bud

import (
	"context"
	"io"
	"os/exec"
	"path/filepath"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/framework/app"
	"github.com/livebud/bud/framework/controller"
	"github.com/livebud/bud/framework/public"
	"github.com/livebud/bud/framework/view"
	"github.com/livebud/bud/framework/web"
	"github.com/livebud/bud/internal/current"
	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/internal/gobuild"
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
	"github.com/livebud/bud/package/watcher"
	"github.com/livebud/bud/runtime/transform"
	"github.com/livebud/bud/runtime/view/dom"
	"github.com/livebud/bud/runtime/view/ssr"
)

// Command is the root commandfor `bud`
type Command struct {
	// Flags
	Dir  string
	Log  string
	Args []string
	Help bool

	// Passed through the subprocesses
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Module finds the go.mod file for the application
func (c *Command) Module() (*gomod.Module, error) {
	return gomod.Find(c.Dir)
}

// BudModule finds the go.mod file for bud itself
func (c *Command) BudModule() (*gomod.Module, error) {
	currentDir, err := current.Directory()
	if err != nil {
		return nil, err
	}
	dir, err := gomod.Absolute(currentDir)
	if err != nil {
		return nil, err
	}
	return gomod.Find(dir)
}

func (c *Command) Logger() (log.Interface, error) {
	handler, err := filter.Load(console.New(c.Stderr), c.Log)
	if err != nil {
		return nil, err
	}
	return log.New(handler), nil
}

func (c *Command) FileSystem(module *gomod.Module, flag *framework.Flag) (*overlay.FileSystem, error) {
	genfs, err := overlay.Load(module)
	if err != nil {
		return nil, err
	}
	parser := parser.New(genfs, module)
	injector := di.New(genfs, module, parser)
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

func (c *Command) FileServer(module *gomod.Module, flag *framework.Flag) (*overlay.Server, error) {
	servefs, err := overlay.Serve(module)
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

// Generate the app
func (c *Command) Generate(genfs *overlay.FileSystem, outDir string) error {
	return genfs.Sync(outDir)
}

func (c *Command) Build(ctx context.Context, module *gomod.Module, mainPath, outPath string) error {
	builder := gobuild.New(module)
	return builder.Build(ctx, mainPath, outPath)
}

func (c *Command) Start(module *gomod.Module, webListener socket.Listener, budListener socket.Listener, flag *framework.Flag) (*exec.Cmd, error) {
	// Start the web server
	cmd := exec.Command(filepath.Join("bud", "app"))
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	cmd.Env = c.Env
	// Add the bud address as an env var if we're not embedding assets
	if !flag.Embed {
		cmd.Env = append(cmd.Env, "BUD_LISTEN="+budListener.Addr().String())
	}
	cmd.Dir = module.Directory()
	// Inject the web listener into the app
	webFile, err := webListener.File()
	if err != nil {
		return nil, err
	}
	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "WEB", webFile)
	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (c *Command) Watch(ctx context.Context, module *gomod.Module, log log.Interface, fn func(isBoot, canHotReload bool) error) error {
	// Wrap the function
	watchFn := func(paths []string) error {
		if err := fn(false, canHotReload(paths)); err != nil {
			log.Error(err.Error())
		}
		return nil
	}
	// Call the function once manually to boot
	if err := fn(true, false); err != nil {
		log.Error(err.Error())
	}
	// Regardless of success, watch for changes
	return watcher.Watch(ctx, module.Directory(), watchFn)
}

// canHotReload returns true if we can incrementally reload a page
func canHotReload(paths []string) bool {
	for _, path := range paths {
		if filepath.Ext(path) == ".go" {
			return false
		}
	}
	return true
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
