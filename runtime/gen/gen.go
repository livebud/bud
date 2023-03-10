package gen

import (
	"os"

	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/runtime/generator"
)

type load = func() (*generator.Generator, error)

func FindModule() (*Module, error) {
	return gomod.Find(".")
}

type Module = gomod.Module

func NewLog() Log {
	return log.New(console.New(os.Stderr))
}

type Log = log.Log

func NewParser(module *Module) *Parser {
	// TODO: do we need to pass genfs in here?
	return parser.New(module, module)
}

type Parser = parser.Parser

func NewInjector(log Log, module *Module, parser *Parser) *Injector {
	// TODO: do we need to pass genfs in here?
	return di.New(module, log, module, parser)
}

type Injector = di.Injector

func Main(load load) {
	log := log.New(console.New(os.Stderr))
	if err := run(load); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func run(load load) error {
	generator, err := load()
	if err != nil {
		return nil
	}
	return generator.Generate()
}
