package cli

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/internal/buildcache"
	"github.com/livebud/bud/internal/cli/create"
	tool_cache_clean "github.com/livebud/bud/internal/cli/tool/cache/clean"
	tool_di "github.com/livebud/bud/internal/cli/tool/di"
	tool_v8 "github.com/livebud/bud/internal/cli/tool/v8"
	tool_v8_client "github.com/livebud/bud/internal/cli/tool/v8/client"
	"github.com/livebud/bud/internal/cli/version"
	"github.com/livebud/bud/internal/envs"
	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/internal/generator/command"
	"github.com/livebud/bud/internal/generator/generator"
	"github.com/livebud/bud/internal/generator/importfile"
	"github.com/livebud/bud/internal/generator/mainfile"
	"github.com/livebud/bud/internal/generator/program"
	"github.com/livebud/bud/internal/generator/transform"
	"github.com/livebud/bud/package/commander"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/modcache"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/parser"
	"github.com/mattn/go-isatty"
)

func Run(ctx context.Context, args ...string) int {
	if err := New(".").Run(ctx, args...); err != nil {
		if !errors.Is(err, context.Canceled) && !isExitStatus(err) {
			console.Error(err.Error())
		}
		return 1
	}
	return 0
}

// New CLI
func New(dir string) *CLI {
	return &CLI{
		dir:    dir,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  stdin(),
		Env: envs.Map{
			"HOME":       os.Getenv("HOME"),
			"PATH":       os.Getenv("PATH"),
			"GOPATH":     os.Getenv("GOPATH"),
			"GOMODCACHE": modcache.Default().Directory(),
			"TMPDIR":     os.TempDir(),
		},
	}
}

type CLI struct {
	dir        string
	help       bool
	args       []string
	Stdout     io.Writer
	Stderr     io.Writer
	Stdin      io.Reader
	Env        envs.Map
	ExtraFiles []*os.File
}

// Dir returns the configured directory
func (c *CLI) Dir() string {
	return c.dir
}

// Inject extra files into the command
func (c *CLI) Inject(prefix string, files ...extrafile.File) error {
	extras, env, err := extrafile.Prepare(prefix, len(c.ExtraFiles), files...)
	if err != nil {
		return err
	}
	c.ExtraFiles = append(c.ExtraFiles, extras...)
	c.Env = c.Env.Append(env...)
	return nil
}

// Run the CLI and wait for the command to finish
func (c *CLI) Run(ctx context.Context, args ...string) error {
	return c.parse(ctx, args, func(ctx context.Context) error {
		cmd, err := c.compile(ctx)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				// When we're outside a go module, we're outside a bud application.
				// In this case, show the bud cli's usage, not your project cli's usage.
				return commander.Usage()
			}
			return err
		}
		return cmd.Run()
	})
}

// Start the CLI but don't wait for the command to finish. This is typically
// used for testing purposes.
func (c *CLI) Start(ctx context.Context, args ...string) (cmd *exe.Cmd, err error) {
	err = c.parse(ctx, args, func(ctx context.Context) error {
		cmd, err = c.compile(ctx)
		if err != nil {
			return err
		}
		return cmd.Start()
	})
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func (c *CLI) parse(ctx context.Context, args []string, fn func(ctx context.Context) error) error {
	cli := commander.New("bud").Writer(c.Stdout)
	cli.Flag("chdir", "change the working directory").Short('C').String(&c.dir).Default(c.dir)
	cli.Flag("help", "show the help message").Short('h').Bool(&c.help).Default(false)
	cli.Args("args").Strings(&c.args).Default(c.args...)
	cli.Run(fn)

	{ // $ bud create <dir>
		cmd := &create.Command{}
		cli := cli.Command("create", "create a new project")
		cli.Arg("dir").String(&cmd.Dir)
		cli.Run(c.create(cmd))
	}

	{ // $ bud tool
		cli := cli.Command("tool", "extra tools")

		{ // $ bud tool di
			// TODO: move into the project CLI since it depends on being within a Go
			// module.
			cmd := &tool_di.Command{}
			cli := cli.Command("di", "dependency injection generator")
			cli.Flag("dependency", "generate dependency provider").Short('d').Strings(&cmd.Dependencies)
			cli.Flag("external", "mark dependency as external").Short('e').Strings(&cmd.Externals).Optional()
			cli.Flag("map", "map interface types to concrete types").Short('m').StringMap(&cmd.Map).Optional()
			cli.Flag("target", "target import path").Short('t').String(&cmd.Target)
			cli.Flag("hoist", "hoist dependencies that depend on externals").Bool(&cmd.Hoist).Default(false)
			cli.Flag("verbose", "verbose logging").Short('v').Bool(&cmd.Verbose).Default(false)
			cli.Run(c.tool_di(cmd))
		}

		{ // $ bud tool v8
			cmd := &tool_v8.Command{Stdin: c.Stdin, Stdout: c.Stdout}
			cli := cli.Command("v8", "Execute Javascript with V8 from stdin")
			cli.Run(cmd.Run)

			{ // $ bud tool v8 client
				cmd := &tool_v8_client.Command{}
				cli := cli.Command("client", "V8 client used during development")
				cli.Run(cmd.Run)
			}
		}

		{ // $ bud tool cache
			cli := cli.Command("cache", "Manage the build cache")

			{ // $ bud tool cache clean
				// TODO: move into the project CLI since it depends on a project
				// existing anyways.
				cmd := &tool_cache_clean.Command{}
				cli := cli.Command("clean", "Clear the cache directory")
				cli.Run(c.tool_cache_clean(cmd))
			}
		}
	}

	{ // $ bud version
		cmd := version.Command{}
		cli := cli.Command("version", "Show package versions")
		cli.Arg("key").String(&cmd.Key).Default("")
		cli.Run(cmd.Run)
	}

	return cli.Parse(ctx, args)
}

func (c *CLI) compile(ctx context.Context) (*exe.Cmd, error) {
	// Find the go.mod
	module, err := gomod.Find(c.dir)
	if err != nil {
		return nil, err
	}

	// Initialize generator dependencies
	genfs, err := overlay.Load(module)
	if err != nil {
		return nil, err
	}
	parser := parser.New(genfs, module)
	injector := di.New(genfs, module, parser)

	// Setup the macros
	genfs.FileGenerator("bud/import.go", importfile.New(module))
	genfs.FileGenerator("bud/.cli/main.go", mainfile.New(module))
	genfs.FileGenerator("bud/.cli/program/program.go", program.New(injector, module))
	genfs.FileGenerator("bud/.cli/command/command.go", command.New(injector, module, parser))
	genfs.FileGenerator("bud/.cli/generator/generator.go", generator.New(genfs, module, parser))
	genfs.FileGenerator("bud/.cli/transform/transform.go", transform.New(module))

	// Synchronize the macros with the filesystem
	if err := genfs.Sync("bud/.cli"); err != nil {
		return nil, err
	}

	// Write the import generator
	importFile, err := fs.ReadFile(genfs, "bud/import.go")
	if err != nil {
		return nil, err
	}
	if err := module.DirFS().WriteFile("bud/import.go", importFile, 0644); err != nil {
		return nil, err
	}

	// Ensure that main.go exists
	if _, err := fs.Stat(module, "bud/.cli/main.go"); err != nil {
		return nil, err
	}

	// Run go build on bud/.cli/main.go outputing to bud/cli
	bcache := buildcache.Default(module)
	if err := bcache.Build(ctx, module, "bud/.cli/main.go", "bud/cli"); err != nil {
		return nil, err
	}

	// Run the project CLI `bud/cli [args...]`
	cmd := exe.Command(ctx, "bud/cli", c.args...)
	cmd.Dir = c.dir
	cmd.Env = c.Env.List()
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	cmd.ExtraFiles = c.ExtraFiles
	return cmd, nil
}

// Wrapper for `bud create`
func (c *CLI) create(cmd *create.Command) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// Incorporate --chdir in dir
		cmd.Dir = filepath.Join(c.dir, cmd.Dir)
		return cmd.Run(ctx)
	}
}

// Wrapper for `bud tool di`
func (c *CLI) tool_di(cmd *tool_di.Command) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// Take into account --chdir
		cmd.Dir = c.dir
		return cmd.Run(ctx)
	}
}

// Wrapper for `bud tool cache clean`
func (c *CLI) tool_cache_clean(cmd *tool_cache_clean.Command) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// Take into account --chdir
		cmd.Dir = c.dir
		return cmd.Run(ctx)
	}
}

func isExitStatus(err error) bool {
	return strings.Contains(err.Error(), "exit status ")
}

// Input from stdin or empty reader by default.
func stdin() io.Reader {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		return strings.NewReader("")
	}
	return os.Stdin
}
