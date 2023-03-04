package transpiler

import (
	_ "embed"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/package/finder"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/gotemplate"
	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/valid"
	"github.com/matthewmueller/gotext"
	"github.com/matthewmueller/text"
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
	Camel   string
	Methods []*Method
}

type Method struct {
	Pascal string // Method name in pascal
	From   string // From extension
	To     string // To extension
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

func (g *Generator) Load(fsys fs.FS) (*State, error) {
	state := new(State)
	imset := imports.New()

	// Load the transpilers
	transpilerDirs, err := finder.Find(fsys, "{transpiler/**.go}", func(path string, isDir bool) (entries []string) {
		if !isDir && valid.GoFile(path) {
			entries = append(entries, filepath.Dir(path))
		}
		return entries
	})
	if err != nil {
		return nil, err
	}

	// Load the custom transpilers
	for _, transpilerDir := range transpilerDirs {
		// Parse the transpiler package
		pkg, err := g.parser.Parse(transpilerDir)
		if err != nil {
			return nil, err
		}
		stct := pkg.Struct("Transpiler")
		importPath := g.module.Import(transpilerDir)
		if stct == nil {
			g.log.Warn("No Transpiler struct in %q. Skipping.", importPath)
			continue
		}
		importName := imset.Add(importPath)
		state.Transpilers = append(state.Transpilers, &Transpiler{
			Import: &imports.Import{
				Name: importName,
				Path: importPath,
			},
			Camel:   gotext.Camel(strings.TrimPrefix(transpilerDir, "transpiler/")),
			Methods: loadMethods(stct),
		})
	}
	if len(state.Transpilers) == 0 {
		return nil, fs.ErrNotExist
	}

	// Setup the imports
	imset.AddNamed("genfs", "github.com/livebud/bud/package/genfs")
	imset.AddNamed("transpiler", "github.com/livebud/bud/runtime/transpiler")
	state.Imports = imset.List()

	return state, nil
}

func loadMethods(stct *parser.Struct) (methods []*Method) {
	for _, method := range stct.Methods() {
		if method.Private() {
			continue
		}
		parts := strings.SplitN(text.Lower(text.Space(method.Name())), " to ", 2)
		if len(parts) != 2 {
			continue
		}
		methods = append(methods, &Method{
			Pascal: gotext.Pascal(method.Name()),
			From:   "." + text.Dot(parts[0]),
			To:     "." + text.Dot(parts[1]),
		})
	}
	return methods
}
