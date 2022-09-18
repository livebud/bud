package transform

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/livebud/bud/package/di"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/internal/scan"
	"github.com/livebud/bud/internal/valid"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"
	"github.com/matthewmueller/gotext"
	"github.com/matthewmueller/text"
)

func Load(fsys fs.FS, injector *di.Injector, log log.Interface, module *gomod.Module, parser *parser.Parser) (*State, error) {
	loader := &loader{
		imports:  imports.New(),
		injector: injector,
		log:      log,
		module:   module,
		parser:   parser,
	}
	return loader.Load(fsys)
}

type loader struct {
	bail.Struct
	imports  *imports.Set
	injector *di.Injector
	log      log.Interface
	module   *gomod.Module
	parser   *parser.Parser
}

// Load the command state
func (l *loader) Load(fsys fs.FS) (state *State, err error) {
	// TODO: for cases like this, we just want to watch, we don't need to
	// return the files.
	if files, err := fs.Glob(fsys, "transform/**.go"); err != nil {
		return nil, err
	} else if len(files) == 0 {
		return nil, fmt.Errorf("framework/transform: no transforms found. %w", fs.ErrNotExist)
	}
	state = new(State)
	state.Transformers = l.loadTransformers(fsys)
	state.Provider = l.loadProvider(state.Transformers)
	l.imports.AddStd("context", "os", "errors")
	l.imports.AddNamed("transformrt", "github.com/livebud/bud/framework/transform2/transformrt")
	l.imports.AddNamed("gomod", "github.com/livebud/bud/package/gomod")
	l.imports.AddNamed("log", "github.com/livebud/bud/package/log")
	l.imports.AddNamed("filter", "github.com/livebud/bud/package/log/filter")
	l.imports.AddNamed("console", "github.com/livebud/bud/package/log/console")
	l.imports.AddNamed("commander", "github.com/livebud/bud/package/commander")
	l.imports.AddNamed("remotefs", "github.com/livebud/bud/package/remotefs")
	l.imports.AddNamed("budfs", "github.com/livebud/bud/package/budfs")
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadTransformers(fsys fs.FS) (transformers []*Transformer) {
	transformDirs, err := scan.List(fsys, "transform", func(de fs.DirEntry) bool {
		if de.IsDir() {
			return valid.Dir(de.Name())
		} else {
			return valid.GoFile(de.Name())
		}
	})
	if err != nil {
		l.Bail(err)
	} else if len(transformDirs) == 0 {
		l.Bail(fmt.Errorf("framework/transform: no transforms found. %w", fs.ErrNotExist))
	}
	for _, transformDir := range transformDirs {
		importPath := l.module.Import(transformDir)
		pkg, err := l.parser.Parse(transformDir)
		if err != nil {
			l.Bail(err)
		}
		// Ensure the package has a transform
		stct := pkg.Struct("Transform")
		if stct == nil {
			l.log.Warn("No Transform struct in %q. Skipping.", importPath)
			continue
		}
		rootlessGenerator := strings.TrimPrefix(transformDir, "transform/")
		transformers = append(transformers, &Transformer{
			Import: &imports.Import{
				Name: l.imports.Add(importPath),
				Path: importPath,
			},
			Path:       rootlessGenerator,
			Pascal:     gotext.Pascal(rootlessGenerator),
			Transforms: l.loadTransforms(stct),
		})
	}
	return transformers
}

func (l *loader) loadTransforms(stct *parser.Struct) (transforms []*Transform) {
	for _, method := range stct.Methods() {
		if method.Private() {
			continue
		}
		parts := strings.SplitN(text.Lower(text.Space(method.Name())), " to ", 2)
		if len(parts) != 2 {
			continue
		}
		transforms = append(transforms, &Transform{
			Name: method.Name(),
			From: "." + text.Dot(parts[0]),
			To:   "." + text.Dot(parts[1]),
		})
	}
	return transforms
}

func (l *loader) loadProvider(transformers []*Transformer) *di.Provider {
	structFields := make([]*di.StructField, len(transformers))
	for i, transformer := range transformers {
		structFields[i] = &di.StructField{
			Name:   transformer.Pascal,
			Import: transformer.Import.Path,
			Type:   "*Transform",
		}
	}
	provider, err := l.injector.Wire(&di.Function{
		Name:    "loadTransformers",
		Target:  l.module.Import("bud/transform"),
		Imports: l.imports,
		Params: []*di.Param{
			{Import: "github.com/livebud/bud/package/log", Type: "Interface"},
			{Import: "github.com/livebud/bud/package/gomod", Type: "*Module"},
			{Import: "context", Type: "Context"},
		},
		Results: []di.Dependency{
			&di.Struct{
				Import: l.module.Import("bud/transform"),
				Type:   "*Transformers",
				Fields: structFields,
			},
			&di.Error{},
		},
	})
	if err != nil {
		l.Bail(err)
	}
	// Add generated imports
	for _, imp := range provider.Imports {
		l.imports.AddNamed(imp.Name, imp.Path)
	}
	return provider
}
