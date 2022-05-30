package program

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/vfs"
	"github.com/livebud/bud/runtime/command"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/gomod"
)

//go:embed program.gotext
var template string

var generator = gotemplate.MustParse("program.gotext", template)

type Program struct {
	Flag     *command.Flag
	Module   *gomod.Module
	Injector *di.Injector
}

type State struct {
	Imports  []*imports.Import
	Provider *di.Provider
}

func (p *Program) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	if err := vfs.Exist(fsys, "bud/.app/command/command.go"); err != nil {
		return err
	}
	// Add the imports
	imports := imports.New()
	imports.AddStd("errors", "context")
	imports.AddNamed("console", "github.com/livebud/bud/package/log/console")
	imports.Add(p.Module.Import("bud/.app/command"))
	jsVM := di.ToType("github.com/livebud/bud/package/js", "VM")
	loadApp := &di.Function{
		Name:   "loadApp",
		Target: p.Module.Import("bud", "program"),
		Params: []di.Dependency{
			di.ToType("github.com/livebud/bud/package/gomod", "*Module"),
			di.ToType("context", "Context"),
		},
		Results: []di.Dependency{
			di.ToType(p.Module.Import("bud", ".app", "command"), "*CLI"),
			&di.Error{},
		},
		Aliases: di.Aliases{
			di.ToType("io/fs", "FS"): di.ToType("github.com/livebud/bud/package/overlay", "*FileSystem"),
			jsVM:                     di.ToType("github.com/livebud/bud/package/js/v8client", "*Client"),
			di.ToType("github.com/livebud/bud/runtime/view", "Renderer"): di.ToType("github.com/livebud/bud/runtime/view", "*Server"),
		},
	}
	if p.Flag.Embed {
		loadApp.Aliases[jsVM] = di.ToType("github.com/livebud/bud/package/js/v8", "*VM")
	}
	provider, err := p.Injector.Wire(loadApp)
	if err != nil {
		// Don't wrap on purpose, this error gets swallowed up easily
		return fmt.Errorf("program: unable to wire. %s", err)
	}
	for _, im := range provider.Imports {
		imports.AddNamed(im.Name, im.Path)
	}
	code, err := generator.Generate(State{
		Imports:  imports.List(),
		Provider: provider,
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
