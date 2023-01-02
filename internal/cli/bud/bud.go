package bud

import (
	"context"
	"io"

	"github.com/livebud/bud/internal/config"
	"github.com/livebud/bud/internal/pubsub"

	"github.com/livebud/bud/package/commander"
	"github.com/livebud/bud/package/socket"
)

// Input contains the configuration that gets passed into the commands
type Input struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    []string

	// Currently passed in only for testing
	Dir   string          // Can be empty
	BudLn socket.Listener // Can be nil
	WebLn socket.Listener // Can be nil
	Bus   pubsub.Client   // Can be nil
}

func New(provide config.Provide) *Command {
	return &Command{provide: provide}
}

type Command struct {
	provide config.Provide
	Args    []string
	Help    bool
}

// Run a custom command
// TODO: finish supporting custom commands
//  1. Compile
//     a. Generate generator (later!)
//     i. Generate bud/internal/generator
//     ii. Build bud/generator
//     iii. Run bud/generator
//     b. Generate custom command
//     i. Generate bud/internal/command/${name}/
//     ii. Build bud/command/${name}
//  2. Run bud/command/${name}
func (c *Command) Run(ctx context.Context) error {
	return commander.Usage()
}

// // Module finds the go.mod file for the application
// func Module(dir string) (*gomod.Module, error) {
// 	return gomod.Find(dir)
// }

// // BudModule finds the module not of your app, but of bud itself
// func BudModule() (*gomod.Module, error) {
// 	dirname, err := current.Directory()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return gomod.Find(dirname)
// }

// func Log(stderr io.Writer, logFilter string) (log.Log, error) {
// 	level, err := log.ParseLevel(logFilter)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return log.New(levelfilter.New(console.New(stderr), level)), nil
// }

// // EnsureVersionAlignment ensures that the CLI and runtime versions are aligned.
// // If they're not aligned, the CLI will correct the go.mod file to align them.
// func EnsureVersionAlignment(ctx context.Context, module *gomod.Module, budVersion string) error {
// 	modfile := module.File()
// 	// Do nothing for the latest version
// 	if budVersion == "latest" {
// 		// If the module file already replaces bud, don't do anything.
// 		if modfile.Replace(`github.com/livebud/bud`) != nil {
// 			return nil
// 		}
// 		// Best effort attempt to replace bud with the latest version.
// 		budModule, err := BudModule()
// 		if err != nil {
// 			return nil
// 		}
// 		// Replace bud with the local version if we found it.
// 		if err := modfile.AddReplace("github.com/livebud/bud", "", budModule.Directory(), ""); err != nil {
// 			return err
// 		}
// 		// Write the go.mod file back to disk.
// 		if err := os.WriteFile(module.Directory("go.mod"), modfile.Format(), 0644); err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	target := "v" + budVersion
// 	require := modfile.Require("github.com/livebud/bud")
// 	// We're good, the CLI matches the runtime version
// 	if require != nil && require.Version == target {
// 		return nil
// 	}
// 	// Otherwise, update go.mod to match the CLI's version
// 	if err := modfile.AddRequire("github.com/livebud/bud", target); err != nil {
// 		return err
// 	}
// 	if err := os.WriteFile(module.Directory("go.mod"), modfile.Format(), 0644); err != nil {
// 		return err
// 	}
// 	// Run `go mod download`
// 	cmd := exec.CommandContext(ctx, "go", "mod", "download")
// 	cmd.Dir = module.Directory()
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	cmd.Env = os.Environ()
// 	if err := cmd.Run(); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func BudListener(in *Input) (net.Listener, error) {
// 	if in.BudLn != nil {
// 		return in.BudLn, nil
// 	}
// 	return socket.Listen(":35729")
// }

// func StartBudServer(ctx context.Context, budln net.Listener, fsys fs.FS, log log.Log) (*budsvr.Server, error) {
// 	bus := pubsub.New()
// 	vm, err := v8.Load()
// 	if err != nil {
// 		return nil, err
// 	}
// 	budServer := budsvr.New(budln, bus, fsys, log, vm)
// 	budServer.Start(ctx)
// 	return budServer, nil
// }

// func BudServer() *budServer {
// 	return &budServer{}
// }

// type budServer struct {
// }

// func (b *budServer) Start(ctx context.Context) error {
// 	return nil
// }

// func (b *budServer) Close() error {
// 	return nil
// }

// func Shell(in *Input, module *gomod.Module) *shell.Command {
// 	cmd := &shell.Command{
// 		Dir:    module.Directory(),
// 		Stdin:  in.Stdin,
// 		Stdout: in.Stdout,
// 		Stderr: in.Stderr,
// 		Env:    []string{},
// 	}
// 	cmd.Env = append(cmd.Env, "GOMODCACHE="+module.ModCache())
// 	cmd.Env = append(cmd.Env, in.Env...)
// 	return cmd
// }
