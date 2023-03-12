package web

import (
	"io/fs"
	"path"
	"path/filepath"

	"github.com/livebud/bud/package/valid"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/package/finder"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/package/parser"
	"github.com/matthewmueller/gotext"
)

func Load(fsys fs.FS, module *gomod.Module, parser *parser.Parser) (*State, error) {
	loader := &loader{
		imports: imports.New(),
		fsys:    fsys,
		module:  module,
		parser:  parser,
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	imports *imports.Set
	fsys    fs.FS
	module  *gomod.Module
	parser  *parser.Parser
}

// Load the command state
func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	// Load all the web handlers
	webDirs, err := finder.Find(l.fsys, "bud/internal/web/*/**.go", func(path string, isDir bool) (entries []string) {
		if !isDir && valid.GoFile(path) {
			entries = append(entries, filepath.Dir(path))
		}
		return entries
	})
	if err != nil {
		return nil, err
	}
	// Add initial imports
	l.imports.AddStd("net/http", "context")
	l.imports.AddNamed("middleware", "github.com/livebud/bud/package/middleware")
	l.imports.AddNamed("webrt", "github.com/livebud/bud/framework/web/webrt")
	l.imports.AddNamed("router", "github.com/livebud/bud/package/router")
	// Show the welcome page if we don't have any web resources
	showWelcome, err := shouldShowWelcome(l.fsys, webDirs)
	if err != nil {
		return nil, err
	}
	if showWelcome {
		const importPath = "github.com/livebud/bud/framework/web/welcome"
		state.Resources = append(state.Resources, &Resource{
			Camel: "welcome",
			Import: &imports.Import{
				Name: l.imports.Add(importPath),
				Path: importPath,
			},
		})
	}
	// Load the resources
	for _, webDir := range webDirs {
		state.Resources = append(state.Resources, l.loadResource(webDir))
	}
	// Load the imports
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadResource(webDir string) (resource *Resource) {
	resource = new(Resource)
	importPath := l.module.Import(webDir)
	resource.Import = &imports.Import{
		Name: l.imports.Add(importPath),
		Path: importPath,
	}
	packageName := path.Base(webDir)
	resource.Camel = gotext.Camel(packageName)
	return resource
}

func shouldShowWelcome(fsys fs.FS, webDirs []string) (bool, error) {
	if len(webDirs) == 0 {
		return true, nil
	} else if len(webDirs) > 1 || webDirs[0] != "bud/internal/web/public" {
		return false, nil
	}
	paths, err := fs.Glob(fsys, "public/**")
	if err != nil {
		return false, err
	}
	if len(paths) == 1 && paths[0] == "public/favicon.ico" {
		return true, nil
	}
	return false, nil
}
