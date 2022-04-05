package program

import (
	"context"
	_ "embed"
	"fmt"

	"gitlab.com/mnm/bud/package/di"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/package/vfs"
	"gitlab.com/mnm/bud/runtime/bud"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/gomod"
)

//go:embed program.gotext
var template string

var generator = gotemplate.MustParse("program.gotext", template)

type Program struct {
	Flag     *bud.Flag
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
	imports.AddNamed("console", "gitlab.com/mnm/bud/package/log/console")
	imports.Add(p.Module.Import("bud/.app/command"))
	jsVM := di.ToType("gitlab.com/mnm/bud/package/js", "VM")
	loadCLI := &di.Function{
		Name:   "loadCLI",
		Target: p.Module.Import("bud", "program"),
		Params: []di.Dependency{
			di.ToType("gitlab.com/mnm/bud/package/gomod", "*Module"),
			di.ToType("context", "Context"),
		},
		Results: []di.Dependency{
			di.ToType(p.Module.Import("bud", ".app", "command"), "*CLI"),
			&di.Error{},
		},
		Aliases: di.Aliases{
			di.ToType("io/fs", "FS"): di.ToType("gitlab.com/mnm/bud/package/overlay", "*FileSystem"),
			jsVM:                     di.ToType("gitlab.com/mnm/bud/package/js/v8client", "*Client"),
			di.ToType("gitlab.com/mnm/bud/runtime/view", "Renderer"): di.ToType("gitlab.com/mnm/bud/runtime/view", "*Server"),
		},
	}
	if p.Flag.Embed {
		loadCLI.Aliases[jsVM] = di.ToType("gitlab.com/mnm/bud/package/js/v8", "*VM")
	}
	provider, err := p.Injector.Wire(loadCLI)
	if err != nil {
		// Don't wrap on purpose, this error gets swallowed up easily
		return fmt.Errorf("program: unable to wire > %s", err)
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
