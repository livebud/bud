package gen

import (
	"context"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/levelfilter"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/package/commander"
	"github.com/livebud/bud/package/virtual"

	"github.com/livebud/bud/internal/dag"

	"github.com/livebud/bud/package/genfs"

	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/runtime/generator"
)

type loadFn = func(*framework.Flag, genfs.FileSystem, *gomod.Module, log.Log) (*generator.Generator, error)

func ProvideParser(genfs genfs.FileSystem, module *gomod.Module) *Parser {
	return parser.New(genfs, module)
}

type Parser = parser.Parser

func ProvideInjector(genfs genfs.FileSystem, log log.Log, module *gomod.Module, parser *Parser) *Injector {
	return di.New(genfs, log, module, parser)
}

type Injector = di.Injector

func Main(load loadFn) {
	ctx := context.Background()
	if err := run(ctx, load, os.Args[1:]); err != nil {
		console.Error(err)
		os.Exit(1)
	}
}

func run(ctx context.Context, load loadFn, args []string) error {
	cmd := &Generate{new(framework.Flag), "info", load}
	cli := commander.New("gen", "generator")
	cli.Flag("embed", "embed assets").Bool(&cmd.flag.Embed).Default(false)
	cli.Flag("hot", "hot reloading").Bool(&cmd.flag.Hot).Default(true)
	cli.Flag("minify", "minify assets").Bool(&cmd.flag.Minify).Default(false)
	cli.Flag("log", "filter logs with this pattern").Short('L').String(&cmd.lvl).Default("info")
	cli.Run(cmd.Run)
	return cli.Parse(ctx, args...)
}

// Generate command
type Generate struct {
	flag *framework.Flag
	lvl  string
	load loadFn
}

// Run the generator
func (g *Generate) Run(ctx context.Context) error {
	lvl, err := log.ParseLevel(g.lvl)
	if err != nil {
		return err
	}
	log := log.New(levelfilter.New(console.New(os.Stderr), lvl))
	module, err := gomod.Find(".")
	if err != nil {
		return err
	}
	fsys := virtual.Exclude(module, exclude)
	gen := genfs.New(dag.Discard, fsys, log)
	generator, err := g.load(g.flag, gen, module, log)
	if err != nil {
		return err
	}
	// Generate the application packages like bud/cmd/app/main.go
	if err := generator.Generate(module, "bud"); err != nil {
		return err
	}
	// Build bud/cmd/app
	cmd := exec.Command("go", "build", "-mod=mod", "-o=bud/app", "./bud/cmd/app")
	cmd.Dir = module.Directory()
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// Avoid deleting files that were synced or built earlier
func exclude(path string) bool {
	if isGenPath(path) || isBudChild(path) {
		return false
	} else if isBudPath(path) {
		return true
	}
	return false
}

// Exclude everything in bud/
func isBudPath(p string) bool {
	return strings.HasPrefix(p, "bud/")
}

func isBudChild(p string) bool {
	return path.Dir(p) == "bud"
}

func isGenPath(p string) bool {
	return strings.HasPrefix(p, "bud/cmd/gen") ||
		strings.HasPrefix(p, "bud/internal/generator") ||
		strings.HasPrefix(p, "bud/pkg/transpiler") ||
		strings.HasPrefix(p, "bud/pkg/viewer")
}
