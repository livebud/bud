package bud

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/livebud/bud/internal/gobuild"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/bfs"
	"github.com/livebud/bud/internal/current"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/versions"
	"golang.org/x/mod/semver"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/filter"
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
func (c *Command) Run(ctx context.Context) error {
	// Find go.mod
	module, err := Module(c.Dir)
	if err != nil {
		return err
	}
	// Ensure we have version alignment between the CLI and the runtime
	if err := EnsureVersionAlignment(ctx, module, versions.Bud); err != nil {
		return err
	}
	// Setup the logger
	log, err := Log(c.in.Stderr, c.Log)
	if err != nil {
		return err
	}
	// TODO: this should probably be configurable
	flag := &framework.Flag{
		Embed:  false,
		Minify: false,
		Hot:    false,
		Stdin:  c.in.Stdin,
		Stdout: c.in.Stdout,
		Stderr: c.in.Stderr,
		Env:    c.in.Env,
	}
	// Load the filesystem
	bfs, err := bfs.Load(flag, log, module)
	if err != nil {
		return err
	}
	defer bfs.Close()
	// Generate the application
	if err := bfs.Sync(); err != nil {
		return err
	}
	// Build the app
	builder := gobuild.New(module)
	if err := builder.Build(ctx, "bud/internal/app/main.go", "bud/app"); err != nil {
		return err
	}
	// Run the command
	cmd := exec.Command("bud/app", c.Args...)
	cmd.Dir = c.Dir
	cmd.Stdin = c.in.Stdin
	cmd.Stdout = c.in.Stdout
	cmd.Stderr = c.in.Stderr
	cmd.Env = c.in.Env
	return cmd.Run()
}

const minGoVersion = "v1.17"

// ErrMinGoVersion error is returned when Bud needs a newer version of Go
var ErrMinGoVersion = fmt.Errorf("bud requires Go %s or later", minGoVersion)

// CheckGoVersion checks if the current version of Go is greater than the
// minimum required Go version.
func CheckGoVersion(currentVersion string) error {
	currentVersion = "v" + strings.TrimPrefix(currentVersion, "go")
	// If we encounter an invalid version, it's probably a development version of
	// Go. We'll let those pass through. Reference:
	// https://github.com/golang/go/blob/3cf79d96105d890d7097d274804644b2a2093df1/src/runtime/extern.go#L273-L275
	if !semver.IsValid(currentVersion) {
		return nil
	}
	if semver.Compare(currentVersion, minGoVersion) < 0 {
		return ErrMinGoVersion
	}
	return nil
}

// Module finds the go.mod file for the application
func Module(dir string) (*gomod.Module, error) {
	return gomod.Find(dir)
}

// BudModule finds the module not of your app, but of bud itself
func BudModule() (*gomod.Module, error) {
	dirname, err := current.Directory()
	if err != nil {
		return nil, err
	}
	return gomod.Find(dirname)
}

func Log(stderr io.Writer, logFilter string) (log.Interface, error) {
	handler, err := filter.Load(console.New(stderr), logFilter)
	if err != nil {
		return nil, err
	}
	return log.New(handler), nil
}

// EnsureVersionAlignment ensures that the CLI and runtime versions are aligned.
// If they're not aligned, the CLI will correct the go.mod file to align them.
func EnsureVersionAlignment(ctx context.Context, module *gomod.Module, budVersion string) error {
	modfile := module.File()
	// Do nothing for the latest version
	if budVersion == "latest" {
		// If the module file already replaces bud, don't do anything.
		if modfile.Replace(`github.com/livebud/bud`) != nil {
			return nil
		}
		// Best effort attempt to replace bud with the latest version.
		budModule, err := BudModule()
		if err != nil {
			return nil
		}
		// Replace bud with the local version if we found it.
		if err := modfile.AddReplace("github.com/livebud/bud", "", budModule.Directory(), ""); err != nil {
			return err
		}
		// Write the go.mod file back to disk.
		if err := os.WriteFile(module.Directory("go.mod"), modfile.Format(), 0644); err != nil {
			return err
		}
		return nil
	}
	target := "v" + budVersion
	require := modfile.Require("github.com/livebud/bud")
	// We're good, the CLI matches the runtime version
	if require != nil && require.Version == target {
		return nil
	}
	// Otherwise, update go.mod to match the CLI's version
	if err := modfile.AddRequire("github.com/livebud/bud", target); err != nil {
		return err
	}
	if err := os.WriteFile(module.Directory("go.mod"), modfile.Format(), 0644); err != nil {
		return err
	}
	// Run `go mod download`
	cmd := exec.CommandContext(ctx, "go", "mod", "download")
	cmd.Dir = module.Directory()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
