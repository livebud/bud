package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/livebud/bud/internal/buildcache"
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
)

func Parse(ctx context.Context, args ...string) int {
	if err := New(".").Parse(ctx, args...); err != nil {
		if !errors.Is(err, context.Canceled) && !isExitStatus(err) {
			console.Error(err.Error())
		}
		return 1
	}
	return 0
}

// New `bud` command
func New(dir string) *Command {
	return &Command{
		dir:    dir,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
		Env: Env{
			"HOME":       os.Getenv("HOME"),
			"PATH":       os.Getenv("PATH"),
			"GOPATH":     os.Getenv("GOPATH"),
			"GOMODCACHE": modcache.Default().Directory(),
			"TMPDIR":     os.TempDir(),
		},
	}
}

type Command struct {
	dir        string
	args       []string
	Stdout     io.Writer
	Stderr     io.Writer
	Stdin      io.Reader
	Env        Env
	ExtraFiles []*os.File
}

func (c *Command) Parse(ctx context.Context, args ...string) error {
	cli := commander.New("bud")
	cli.Flag("chdir", "Change directory").Short('C').String(&c.dir).Default(c.dir)
	cli.Args("args").Strings(&c.args).Default(c.args...)
	if err := cli.Parse(ctx, args); err != nil {
		fmt.Println(c.dir)
		fmt.Println(err)
		return err
	}
	cli.Run(c.run)
	return nil
}

func (c *Command) run(ctx context.Context) error {

	// Find the go.mod
	module, err := gomod.Find(c.dir)
	if err != nil {
		return err
	}

	// Initialize generator dependencies
	genfs, err := overlay.Load(module)
	if err != nil {
		return err
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
		return err
	}

	// Write the import generator
	// TODO: add import writer back in
	// if err := c.writeImporter(ctx, overlay); err != nil {
	// 	return err
	// }

	// Ensure that main.go exists
	if _, err := fs.Stat(module, "bud/.cli/main.go"); err != nil {
		return err
	}

	// Run go build on bud/.cli/main.go outputing to bud/cli
	bcache := buildcache.Default(module)
	if err := bcache.Build(ctx, module, "bud/.cli/main.go", "bud/cli"); err != nil {
		return err
	}

	// Run the project CLI `bud/cli [args...]`
	cmd := exe.Command(ctx, "bud/cli", c.args...)
	cmd.Dir = c.dir
	cmd.Env = c.Env.List()
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	cmd.ExtraFiles = c.ExtraFiles
	return cmd.Run()
}

// func parse(args []string) error {
// 	// $ bud
// 	cmd := &custom.Command{
// 		Stderr: os.Stderr,
// 		Stdout: os.Stdout,
// 	}
// 	cli := commander.New("bud")
// 	cli.Flag("chdir", "Change the working directory").Short('C').String(&cmd.Dir).Default(".")
// 	cli.Args("args").Strings(&cmd.Args)
// 	cli.Run(cmd.Run)

// 	// { // $ bud create <dir>
// 	// 	cmd := createCommand(bud)
// 	// 	cli := cli.Command("create", "create a new project")
// 	// 	cli.Arg("dir").String(&cmd.Dir)
// 	// 	cli.Run(cmd.Run)
// 	// }

// 	// { // $ bud run
// 	// 	cmd := runCommand(bud)
// 	// 	cli := cli.Command("run", "run the development server")
// 	// 	cli.Flag("embed", "embed the assets").Bool(&cmd.Embed).Default(false)
// 	// 	cli.Flag("hot", "hot reload the frontend").String(&cmd.Hot).Default("35729")
// 	// 	cli.Flag("minify", "minify the assets").Bool(&cmd.Minify).Default(false)
// 	// 	cli.Flag("listen", "listen on an address").String(&cmd.Listen).Default("3000")
// 	// 	cli.Run(cmd.Run)
// 	// }

// 	// { // $ bud build
// 	// 	cmd := buildCommand(bud)
// 	// 	cli := cli.Command("build", "build the production server")
// 	// 	cli.Flag("embed", "embed the assets").Bool(&cmd.Embed).Default(true)
// 	// 	cli.Flag("hot", "hot reload the frontend").String(&cmd.Hot).Default("")
// 	// 	cli.Flag("minify", "minify the assets").Bool(&cmd.Minify).Default(true)
// 	// 	cli.Run(cmd.Run)
// 	// }

// 	// { // $ bud tool
// 	// 	cli := cli.Command("tool", "extra tools")

// 	// 	{ // $ bud tool di
// 	// 		cmd := diCommand(bud)
// 	// 		cli := cli.Command("di", "dependency injection generator")
// 	// 		cli.Flag("dependency", "generate dependency provider").Short('d').Strings(&cmd.Dependencies)
// 	// 		cli.Flag("external", "mark dependency as external").Short('e').Strings(&cmd.Externals).Optional()
// 	// 		cli.Flag("map", "map interface types to concrete types").Short('m').StringMap(&cmd.Map).Optional()
// 	// 		cli.Flag("target", "target import path").Short('t').String(&cmd.Target)
// 	// 		cli.Flag("hoist", "hoist dependencies that depend on externals").Bool(&cmd.Hoist).Default(false)
// 	// 		cli.Flag("verbose", "verbose logging").Short('v').Bool(&cmd.Verbose).Default(false)
// 	// 		cli.Run(cmd.Run)
// 	// 	}

// 	// 	{ // $ bud tool v8
// 	// 		cmd := v8.New()
// 	// 		cli := cli.Command("v8", "Execute Javascript with V8 from stdin")
// 	// 		cli.Run(cmd.Run)

// 	// 		{ // $ bud tool v8 client
// 	// 			cmd := v8client.New()
// 	// 			cli := cli.Command("client", "V8 client used during development")
// 	// 			cli.Run(cmd.Run)
// 	// 		}
// 	// 	}

// 	// 	{ // $ bud tool cache
// 	// 		cmd := cache.New(bud)
// 	// 		cli := cli.Command("cache", "Manage the build cache")

// 	// 		{ // $ bud tool cache clean
// 	// 			cli := cli.Command("clean", "Clear the cache directory")
// 	// 			cli.Run(cmd.Clean)
// 	// 		}
// 	// 	}
// 	// }

// 	// { // $ bud version
// 	// 	cmd := version.New()
// 	// 	cli := cli.Command("version", "Show package versions")
// 	// 	cli.Arg("key").String(&cmd.Key).Default("")
// 	// 	cli.Run(cmd.Run)
// 	// }

// 	ctx := context.Background()
// 	return cli.Parse(ctx, args)
// }

func isExitStatus(err error) bool {
	return strings.Contains(err.Error(), "exit status ")
}
