package command

import (
	"context"
	"io"
	"path/filepath"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/internal/gobuild"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/exe"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/package/svelte"
	"github.com/livebud/bud/package/watcher"
	"github.com/livebud/bud/runtime/transform"

	"github.com/livebud/bud/framework/app"
	"github.com/livebud/bud/framework/controller"
	"github.com/livebud/bud/framework/public"
	"github.com/livebud/bud/framework/view"
	"github.com/livebud/bud/framework/web"
	"github.com/livebud/bud/package/overlay"

	"github.com/livebud/bud/package/commander"

	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/filter"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
)

// Bud command
type Bud struct {
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

func (c *Bud) Module() (*gomod.Module, error) {
	return gomod.Find(c.Dir)
}

func (c *Bud) Logger() (log.Interface, error) {
	handler, err := filter.Load(console.New(c.Stderr), c.Log)
	if err != nil {
		return nil, err
	}
	return log.New(handler), nil
}

// Generate the app
func (c *Bud) Generate(module *gomod.Module, flag *framework.Flag, outDir string) error {
	genfs, err := overlay.Load(module)
	if err != nil {
		return err
	}
	parser := parser.New(genfs, module)
	injector := di.New(genfs, module, parser)
	vm, err := v8.Load()
	if err != nil {
		return err
	}
	svelteCompiler, err := svelte.Load(vm)
	if err != nil {
		return err
	}
	transformMap, err := transform.Load(svelte.NewTransformable(svelteCompiler))
	if err != nil {
		return err
	}
	genfs.FileGenerator("bud/internal/app/main.go", app.New(injector, module, flag))
	genfs.FileGenerator("bud/internal/app/web/web.go", web.New(module, parser))
	genfs.FileGenerator("bud/internal/app/controller/controller.go", controller.New(injector, module, parser))
	genfs.FileGenerator("bud/internal/app/view/view.go", view.New(module, transformMap, flag))
	genfs.FileGenerator("bud/internal/app/public/public.go", public.New(flag))
	return genfs.Sync(outDir)
}

func (c *Bud) Build(ctx context.Context, module *gomod.Module, mainPath, outPath string) error {
	builder := gobuild.New(module)
	return builder.Build(ctx, mainPath, outPath)
}

func (c *Bud) Start(module *gomod.Module, webListener socket.Listener) (*exe.Cmd, error) {
	// Start the web server
	cmd := exe.Command(context.Background(), filepath.Join("bud", "app"))
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	cmd.Env = c.Env
	cmd.Dir = module.Directory()
	// Inject the web listener into the app
	webFile, err := webListener.File()
	if err != nil {
		return nil, err
	}
	// TODO: rename to WEB
	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "APP", webFile)
	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (c *Bud) Watch(ctx context.Context, module *gomod.Module, log log.Interface, fn func(isBoot, canHotReload bool) error) error {
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
func (c *Bud) Run(ctx context.Context) error {
	return commander.Usage()
}
