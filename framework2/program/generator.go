package program

import (
	_ "embed"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/pkg/di"
)

//go:embed program.gotext
var template string

var generator = gotemplate.MustParse("program.gotext", template)

// State of the program code
type State struct {
	Imports  []*imports.Import
	Provider *di.Provider
}

// Generate the program
func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

// // New program generator
// func New(injector *di.Injector, module *gomod.Module, fn *di.Function) overlay.GenerateFile {
// 	return func(ctx context.Context, fsys overlay.F, file *overlay.File) error {
// 		// Default  imports
// 		imports := imports.New()
// 		imports.AddStd("os", "errors", "context", "path/filepath", "runtime")
// 		imports.AddNamed("console", "gitlab.com/mnm/bud/pkg/log/console")
// 		imports.AddNamed("trace", "gitlab.com/mnm/bud/package/trace")
// 		// Write up the dependencies
// 		provider, err := injector.Wire(fn)
// 		if err != nil {
// 			// Don't wrap on purpose, this error gets swallowed up too easily
// 			return fmt.Errorf("programg: unable to wire > %s", err)
// 		}
// 		// Add additional imports that we brought in
// 		for _, im := range provider.Imports {
// 			imports.AddNamed(im.Name, im.Path)
// 		}
// 		// Generate code from the state
// 		code, err := generator.Generate(state{imports.List(), provider})
// 		if err != nil {
// 			return err
// 		}
// 		// Write to the file
// 		file.Data = code
// 		return nil
// 	}
// }
