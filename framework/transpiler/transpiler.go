package transpiler

import (
	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"
)

// New transform generator
func New(flag *framework.Flag, log log.Log, module *gomod.Module, parser *parser.Parser) *Generator {
	return &Generator{flag, log, module, parser}
}

type Generator struct {
	flag   *framework.Flag
	log    log.Log
	module *gomod.Module
	parser *parser.Parser
}

func (g *Generator) GenerateFile(fsys genfs.FS, file *genfs.File) error {
	file.Data = []byte(`
		package transpiler
		import (
			"github.com/livebud/bud/package/genfs"
			"github.com/livebud/bud/runtime/transpiler"
			"app.com/transpiler/doubler"
			"io/fs"
		)
		func Load(doubler *doubler.Transpiler) *Generator {
			tr := transpiler.New()
			tr.Add(".svelte", ".svelte", doubler.SvelteToSvelte)
			return &Generator{tr}
		}
		type Generator struct {
			tr transpiler.Interface
		}
		func (g *Generator) Serve(fsys genfs.FS, file *genfs.File) error {
			toExt, inputPath := transpiler.SplitRoot(file.Relative())
			input, err := fs.ReadFile(fsys, inputPath)
			if err != nil {
				return err
			}
			output, err := g.tr.Transpile(file.Ext(), toExt, input)
			if err != nil {
				return err
			}
			file.Data = output
			return nil
		}
	`)
	return nil
}
