package buddy

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/livebud/bud/framework/appfs"
	generator "github.com/livebud/bud/framework/generator2"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/log/levelfilter"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/remotefs"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"

	"github.com/livebud/bud/internal/errs"
	"github.com/livebud/bud/internal/exe"
	"github.com/livebud/bud/internal/gobuild"
)

type Input struct {
	Embed  bool
	Minify bool
	Hot    bool
	Log    string
	Dir    string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    []string
}

func Load(ctx context.Context, in *Input) (*Driver, error) {
	// Load the logger
	lvl, err := log.ParseLevel(in.Log)
	if err != nil {
		return nil, err
	}
	log := log.New(levelfilter.New(console.New(in.Stderr), lvl))
	// Load the module
	module, err := gomod.Find(in.Dir)
	if err != nil {
		return nil, err
	}
	// Create a command template from the input for future commands
	exec := &exe.Template{
		Stdin:  in.Stdin,
		Stdout: in.Stdout,
		Stderr: in.Stderr,
		Env:    in.Env,
	}
	// Ensure the bud CLI is aligned with the runtime version
	if err := ensureVersionAlignment(ctx, module, exec); err != nil {
		return nil, err
	}
	// Setup go builder
	gobuild := gobuild.New(module)
	gobuild.Env = in.Env
	gobuild.Stdin = in.Stdin
	gobuild.Stderr = in.Stderr
	gobuild.Stdout = in.Stdout
	// Load the flags
	flag := &framework.Flag{
		Embed:  in.Embed,
		Minify: in.Minify,
		Hot:    in.Hot,
	}
	// Load the driver
	driver := &Driver{
		// bfs:      bfs,
		flag:    flag,
		gobuild: gobuild,
		// injector: injector,
		log:    log,
		module: module,
		// parser:   parser,
		exec: exec,
	}
	return driver, nil
}

type Driver struct {
	// bfs      *budfs.FileSystem
	flag    *framework.Flag
	gobuild *gobuild.Builder
	// injector *di.Injector
	log    log.Log
	module *gomod.Module
	// parser   *parser.Parser
	exec *exe.Template
}

// EnsureVersionAlignment ensures that the CLI and runtime versions are aligned.
// If they're not aligned, the CLI will correct the go.mod file to align them.
func ensureVersionAlignment(ctx context.Context, module *gomod.Module, exec *exe.Template) error {
	// modfile := module.File()
	// // Do nothing for the latest version
	// if budVersion == "latest" {
	// 	// If the module file already replaces bud, don't do anything.
	// 	if modfile.Replace(`github.com/livebud/bud`) != nil {
	// 		return nil
	// 	}
	// 	// Best effort attempt to replace bud with the latest version.
	// 	budModule, err := BudModule()
	// 	if err != nil {
	// 		return nil
	// 	}
	// 	// Replace bud with the local version if we found it.
	// 	if err := modfile.AddReplace("github.com/livebud/bud", "", budModule.Directory(), ""); err != nil {
	// 		return err
	// 	}
	// 	// Write the go.mod file back to disk.
	// 	if err := os.WriteFile(module.Directory("go.mod"), modfile.Format(), 0644); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }
	// target := "v" + budVersion
	// require := modfile.Require("github.com/livebud/bud")
	// // We're good, the CLI matches the runtime version
	// if require != nil && require.Version == target {
	// 	return nil
	// }
	// // Otherwise, update go.mod to match the CLI's version
	// if err := modfile.AddRequire("github.com/livebud/bud", target); err != nil {
	// 	return err
	// }
	// if err := os.WriteFile(module.Directory("go.mod"), modfile.Format(), 0644); err != nil {
	// 	return err
	// }
	// // Run `go mod download`
	// cmd := exec.CommandContext(ctx, "go", "mod", "download")
	// cmd.Dir = module.Directory()
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// cmd.Env = os.Environ()
	// if err := cmd.Run(); err != nil {
	// 	return err
	// }
	// return nil
	return nil
}

// Directories needed to be expanded
var expandDirs = []string{
	"bud/cmd/appfs",
	"bud/internal/generator",
}

// Generate the code. Process can be nil if it wasn't needed.
func (d *Driver) Generate(ctx context.Context, dirs ...string) (*remotefs.Process, error) {
	// Load the filesystem
	bfs := budfs.New(d.module, d.log)
	// Load the parser
	parser := parser.New(bfs, d.module)
	// Load the injector
	injector := di.New(bfs, d.log, d.module, parser)
	// Setup the initial file generators
	bfs.FileGenerator("bud/cmd/appfs/main.go", appfs.New(d.flag, injector, d.log, d.module))
	bfs.FileGenerator("bud/internal/generator/generator.go", generator.New(bfs, d.flag, injector, d.log, d.module, parser))
	// If all dirs are part of the appfs phase, so we can just sync them and
	// return early.
	budDirs := findBudDirs(dirs)
	if len(budDirs) > 0 && withinAppFS(budDirs) {
		if err := bfs.SyncMany(d.module, budDirs); err != nil {
			return nil, err
		}
		// Exit early
		return nil, nil
	}
	// Sync bud directories needed for appfs
	if err := bfs.SyncMany(d.module, expandDirs); err != nil {
		return nil, err
	}
	// Build the application file server
	if err := d.gobuild.Build(ctx, "bud/cmd/appfs/main.go", "bud/appfs"); err != nil {
		return nil, fmt.Errorf("buddy: error building appfs. %w", err)
	}
	// Start the application file server
	// TODO: use exe.Template within remotefs
	cmd := &remotefs.Command{
		Dir:    d.module.Directory(),
		Env:    append([]string{}, d.exec.Env...),
		Stderr: d.exec.Stderr,
		Stdout: d.exec.Stdout,
	}
	// Start the process
	appFS, err := cmd.Start(ctx, "bud/appfs")
	if err != nil {
		return nil, err
	}
	// Mount the remote application file server into the existing file system.
	bfs.Mount(appFS)
	// Sync only the requested files
	if len(budDirs) > 0 {
		if err := bfs.SyncMany(d.module, budDirs); err != nil {
			err = errs.Join(err, appFS.Close())
			return nil, err
		}
		return appFS, nil
	}
	// Otherwise sync all of bud
	if err := bfs.Sync(d.module, "bud"); err != nil {
		err = errs.Join(err, appFS.Close())
		return nil, err
	}
	return appFS, nil
}

func (d *Driver) Build(ctx context.Context) error {
	return nil
}

func (d *Driver) Run(ctx context.Context) error {
	if err := d.Build(ctx); err != nil {
		return err
	}
	// TODO: add run
	return nil
}

func (d *Driver) Watch(ctx context.Context) error {
	return nil
}

func findBudDirs(patterns []string) (paths []string) {
	for _, pattern := range patterns {
		// Only sync from within the bud directory
		if !strings.HasPrefix(pattern, "bud/") {
			continue
		}
		// Trim the wildcard suffix since SyncDirs is recursive already
		// TODO: support non-recursive syncs
		pattern = strings.TrimSuffix(pattern, "/...")
		// Add the file or directory
		paths = append(paths, pattern)
	}
	return paths
}

// Returns false if we need appfs in order to run the requested generators.
// This assumes that findBudDir was already run
func withinAppFS(dirs []string) bool {
	for _, dir := range dirs {
		if !inAppFS(dir) {
			return false
		}
	}
	return true
}

// inAppFS returns true if the given paths are used within appfs
func inAppFS(dir string) bool {
	return strings.HasPrefix(dir, "bud/cmd/appfs") ||
		strings.HasPrefix(dir, "bud/internal/generator")
}
