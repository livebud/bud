package config

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/budfs"
	"github.com/livebud/bud/internal/current"
	"github.com/livebud/bud/internal/prompter"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/shell"
	"github.com/livebud/bud/package/budhttp/budsvr"
	"github.com/livebud/bud/package/gomod"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/levelfilter"
	"github.com/livebud/bud/package/socket"
	"golang.org/x/mod/semver"
)

type Provide interface {
	Module() (*gomod.Module, error)
	Logger() (*log.Logger, error)
	Command() *shell.Command
	BudFileSystem() (budfs.FileSystem, error)
	BudServer() (*budsvr.Server, error)
	BudListener() (socket.Listener, error)
	WebListener() (socket.Listener, error)
	Prompter() (*prompter.Prompter, error)
	Bus() pubsub.Client
	V8() (*v8.VM, error)
}

func New() *Config {
	return &Config{
		Dir:    ".",
		Log:    "info",
		Listen: ":3000",
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    os.Environ(),
	}
}

type Config struct {
	Log    string
	Dir    string
	Listen string

	Flag framework.Flag

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    []string

	// Currently passed in only for testing
	BudLn  socket.Listener // Can be nil
	WebLn  socket.Listener // Can be nil
	PubSub pubsub.Client   // Can be nil

	budfs    budfs.FileSystem
	budsvr   *budsvr.Server
	logger   *log.Logger
	module   *gomod.Module
	prompter *prompter.Prompter
	v8       *v8.VM
}

var _ Provide = (*Config)(nil)

func (c *Config) Module() (*gomod.Module, error) {
	if c.module != nil {
		return c.module, nil
	}
	module, err := gomod.Find(c.Dir)
	if err != nil {
		return nil, err
	}
	c.module = module
	return c.module, nil
}

func (c *Config) Logger() (*log.Logger, error) {
	if c.logger != nil {
		return c.logger, nil
	}
	level, err := log.ParseLevel(c.Log)
	if err != nil {
		return nil, err
	}
	c.logger = log.New(levelfilter.New(console.New(c.Stderr), level))
	return c.logger, nil
}

func (c *Config) Command() *shell.Command {
	// Create a clean command everytime
	return &shell.Command{
		Dir:    c.Dir,
		Stdin:  c.Stdin,
		Stdout: c.Stdout,
		Stderr: c.Stderr,
		Env:    append([]string{}, c.Env...),
	}
}

func (c *Config) BudServer() (*budsvr.Server, error) {
	if c.budsvr != nil {
		return c.budsvr, nil
	}
	budln, err := c.BudListener()
	if err != nil {
		return nil, err
	}
	bus := c.Bus()
	budfs, err := c.BudFileSystem()
	if err != nil {
		return nil, err
	}
	log, err := c.Logger()
	if err != nil {
		return nil, err
	}
	vm, err := c.V8()
	if err != nil {
		return nil, err
	}
	c.budsvr = budsvr.New(budln, bus, &c.Flag, budfs, log, vm)
	c.budsvr.Start(context.Background())
	return c.budsvr, nil
}

func (c *Config) BudListener() (socket.Listener, error) {
	if c.BudLn != nil {
		return c.BudLn, nil
	}
	budln, err := socket.Listen(":35729")
	if err != nil {
		return nil, err
	}
	c.BudLn = budln
	return budln, nil
}

func (c *Config) BudFileSystem() (budfs.FileSystem, error) {
	if c.budfs != nil {
		return c.budfs, nil
	}
	budln, err := c.BudListener()
	if err != nil {
		return nil, err
	}
	cmd := c.Command()
	cmd.Env = append(cmd.Env, "BUD_LISTEN="+budln.Addr().String())
	module, err := c.Module()
	if err != nil {
		return nil, err
	}
	log, err := c.Logger()
	if err != nil {
		return nil, err
	}
	c.budfs, err = budfs.Load(budln, cmd, &c.Flag, module, log)
	if err != nil {
		return nil, err
	}
	return c.budfs, nil
}

func (c *Config) WebListener() (socket.Listener, error) {
	if c.WebLn != nil {
		return c.WebLn, nil
	}
	// Listen and increment if the port is already in use up to 10 times
	webln, err := socket.ListenUp(c.Listen, 10)
	if err != nil {
		return nil, err
	}
	c.WebLn = webln
	return c.WebLn, nil
}

func (c *Config) Bus() pubsub.Client {
	if c.PubSub != nil {
		return c.PubSub
	}
	c.PubSub = pubsub.New()
	return c.PubSub
}

func (c *Config) V8() (*v8.VM, error) {
	if c.v8 != nil {
		return c.v8, nil
	}
	v8, err := v8.Load()
	if err != nil {
		return nil, err
	}
	c.v8 = v8
	return c.v8, nil
}

func (c *Config) Prompter() (*prompter.Prompter, error) {
	if c.prompter != nil {
		return c.prompter, nil
	}
	webln, err := c.WebListener()
	if err != nil {
		return nil, err
	}
	var prompter prompter.Prompter
	c.Stdout = io.MultiWriter(c.Stdout, &prompter.StdOut)
	c.Stderr = io.MultiWriter(c.Stderr, &prompter.StdErr)
	c.prompter = &prompter
	c.prompter.Init(webln.Addr().String())
	return c.prompter, nil
}

// BudModule finds the module not of your app, but of bud itself
func BudModule() (*gomod.Module, error) {
	dirname, err := current.Directory()
	if err != nil {
		return nil, err
	}
	return gomod.Find(dirname)
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
