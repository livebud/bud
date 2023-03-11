package gen

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/livebud/bud/package/virtual"

	"github.com/livebud/bud/internal/dag"

	"github.com/livebud/bud/package/genfs"

	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/runtime/generator"
)

type load = func(genfs.FileSystem, *gomod.Module, log.Log) (*generator.Generator, error)

func ProvideParser(genfs genfs.FileSystem, module *gomod.Module) *Parser {
	return parser.New(genfs, module)
}

type Parser = parser.Parser

func ProvideInjector(genfs genfs.FileSystem, log log.Log, module *gomod.Module, parser *Parser) *Injector {
	return di.New(genfs, log, module, parser)
}

type Injector = di.Injector

func Main(load load) {
	log := log.New(console.New(os.Stderr))
	if err := run(log, load); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func run(log log.Log, load load) error {
	module, err := gomod.Find(".")
	if err != nil {
		return err
	}
	fsys := virtual.Exclude(module, exclude)
	gen := genfs.New(dag.Discard, fsys, log)
	generator, err := load(gen, module, log)
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
		strings.HasPrefix(p, "bud/internal/gen") ||
		strings.HasPrefix(p, "bud/pkg/gen")
}
