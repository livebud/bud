package app

import (
	"fmt"
	"io/fs"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/vfs"
)

func Load(fsys fs.FS, injector *di.Injector, module *gomod.Module, flag *framework.Flag) (*State, error) {
	if err := vfs.Exist(fsys, "bud/internal/app/web/web.go"); err != nil {
		return nil, err
	}
	return (&loader{
		fsys:     fsys,
		injector: injector,
		module:   module,
		flag:     flag,
		imports:  imports.New(),
	}).Load()
}

type loader struct {
	fsys     fs.FS
	injector *di.Injector
	module   *gomod.Module
	flag     *framework.Flag

	imports *imports.Set
	bail.Struct
}

func (l *loader) Load() (state *State, err error) {
	defer l.Recover2(&err, "app: unable to load state")
	state = new(State)
	l.imports.AddStd("os", "context", "errors")
	l.imports.AddNamed("commander", "github.com/livebud/bud/package/commander")
	l.imports.AddNamed("budhttp", "github.com/livebud/bud/package/budhttp")
	l.imports.AddNamed("console", "github.com/livebud/bud/package/log/console")
	l.imports.AddNamed("log", "github.com/livebud/bud/package/log")
	l.imports.AddNamed("filter", "github.com/livebud/bud/package/log/filter")
	l.imports.Add(l.module.Import("bud/internal/app/web"))
	state.Provider = l.loadProvider()
	state.Flag = l.flag
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadProvider() *di.Provider {
	jsVM := di.ToType("github.com/livebud/bud/package/js", "VM")
	// TODO: the public generator should be able to configure this
	publicServer := di.ToType("github.com/livebud/bud/framework/public/publicrt", "Server")
	fn := &di.Function{
		Name:    "loadWeb",
		Imports: l.imports,
		Target:  l.module.Import("bud", "program"),
		Params: []*di.Param{
			{Import: "github.com/livebud/bud/package/log", Type: "Interface"},
			{Import: "github.com/livebud/bud/package/gomod", Type: "*Module"},
			{Import: "github.com/livebud/bud/package/budhttp", Type: "Client"},
			{Import: "context", Type: "Context"},
		},
		Results: []di.Dependency{
			di.ToType(l.module.Import("bud/internal/app/web"), "*Server"),
			&di.Error{},
		},
		Aliases: di.Aliases{
			publicServer: di.ToType("github.com/livebud/bud/framework/public/publicrt", "*LiveServer"),
			jsVM:         di.ToType("github.com/livebud/bud/package/budhttp", "Client"),
		},
	}
	if l.flag.Embed {
		fn.Aliases[jsVM] = di.ToType("github.com/livebud/bud/package/js/v8", "*VM")
		fn.Aliases[publicServer] = di.ToType("github.com/livebud/bud/framework/public/publicrt", "*StaticServer")
	}
	provider, err := l.injector.Wire(fn)
	if err != nil {
		// Intentionally don't wrap the error. The error gets swallowed up too
		// easily
		l.Bail(fmt.Errorf("app: unable to wire. %s", err))
	}
	return provider
}
