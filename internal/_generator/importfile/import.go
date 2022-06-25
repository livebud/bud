package importfile

import (
	"context"
	_ "embed"
	"path"
	"strings"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/overlay"
)

//go:embed import.gotext
var template string

var generator = gotemplate.MustParse("import.gotext", template)

// State for the generator
type State struct {
	Programs []*imports.Import
	Plugins  []*imports.Import // TODO
}

func New(module *gomod.Module) *Generator {
	return &Generator{module}
}

type Generator struct {
	module *gomod.Module
}

func (g *Generator) Parse(ctx context.Context) (*State, error) {
	state := new(State)
	state.Programs = []*imports.Import{
		{
			Name: "_",
			Path: g.module.Import("bud/.cli/program"),
		},
		{
			Name: "_",
			Path: g.module.Import("bud/.app/program"),
		},
	}
	// Range over requires
	// TODO: consolidate with pluginfs
	for _, req := range g.module.File().Requires() {
		// The last path in the module path needs to start with "bud-"
		if strings.HasPrefix(path.Base(req.Mod.Path), "bud-") {
			state.Plugins = append(state.Plugins, &imports.Import{
				Name: "_",
				Path: req.Mod.Path,
			})
		}
	}
	return state, nil
}

// Generate a main file
func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func (g *Generator) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	state, err := g.Parse(ctx)
	if err != nil {
		return err
	}
	// Generate code from the state
	code, err := Generate(state)
	if err != nil {
		return err
	}
	// Write to the file
	file.Data = code
	return nil
}
