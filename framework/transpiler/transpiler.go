package transpiler

import (
	_ "embed"
	"io/fs"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"
)

//go:embed transpiler.gotext
var template string

var generator = gotemplate.MustParse("framework/transpiler/transpiler.gotext", template)

// New transpiler generator
func New(flag *framework.Flag, log log.Log, module *gomod.Module, parser *parser.Parser) *Generator {
	return &Generator{flag, log, module, parser}
}

type Generator struct {
	flag   *framework.Flag
	log    log.Log
	module *gomod.Module
	parser *parser.Parser
}

type State struct {
	Imports     []*imports.Import
	Transpilers []*Transpiler
}

type Transpiler struct {
	Import  *imports.Import
	FromExt string
	ToExt   string
	Method  string
	Camel   string
}

func (g *Generator) GenerateFile(fsys genfs.FS, file *genfs.File) error {
	state, err := g.Load(fsys)
	if err != nil {
		return err
	}
	code, err := generator.Generate(state)
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}

// file.Data = []byte(`
// package transpiler
// import (
// 	"github.com/livebud/bud/package/genfs"
// 	"github.com/livebud/bud/runtime/transpiler"
// 	"app.com/transpiler/doubler"
// 	"io/fs"
// )
// func Load(doubler *doubler.Transpiler) *Generator {
// 	tr := transpiler.New()
// 	tr.Add(".svelte", ".svelte", doubler.SvelteToSvelte)
// 	return &Generator{tr}
// }
// type Generator struct {
// 	tr transpiler.Interface
// }
// func (g *Generator) Serve(fsys genfs.FS, file *genfs.File) error {
// 	toExt, inputPath := transpiler.SplitRoot(file.Relative())
// 	input, err := fs.ReadFile(fsys, inputPath)
// 	if err != nil {
// 		return err
// 	}
// 	output, err := g.tr.Transpile(file.Ext(), toExt, input)
// 	if err != nil {
// 		return err
// 	}
// 	file.Data = output
// 	return nil
// }
// `)

func (g *Generator) Load(fsys fs.FS) (*State, error) {
	state := new(State)
	imports := imports.New()
	imports.AddNamed("genfs", "github.com/livebud/bud/package/genfs")
	imports.AddNamed("transpiler", "github.com/livebud/bud/runtime/transpiler")
	state.Imports = imports.List()
	return state, nil
}
