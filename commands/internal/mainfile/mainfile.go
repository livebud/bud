package mainfile

import (
	_ "embed"

	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/gotemplate"
	"github.com/livebud/bud/package/imports"
)

//go:embed mainfile.gotext
var mainCode string

var mainFile = gotemplate.MustParse("mainfile.gotext", mainCode)

// MainFile generates a main file. Typically used for testing purposes.
func New(injector *di.Injector, module *gomod.Module) genfs.GenerateFile {
	return func(fsys genfs.FS, file *genfs.File) error {
		imset := imports.New()
		imset.AddStd("os", "context", "fmt")
		provider, err := injector.Wire(&di.Function{
			Name:    "loadCommands",
			Imports: imset,
			Results: []di.Dependency{
				di.ToType(module.Import("bud", "command"), "*CLI"),
				&di.Error{},
			},
		})
		if err != nil {
			return err
		}
		code, err := mainFile.Generate(struct {
			Imports  *imports.Set
			Provider *di.Provider
		}{
			Imports:  imset,
			Provider: provider,
		})
		if err != nil {
			return err
		}
		file.Data = code
		return nil
	}
}
