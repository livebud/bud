package generator

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/matthewmueller/gotext"

	"github.com/livebud/bud/package/finder"

	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/package/valid"
)

type loader struct {
	*Generator
	imports *imports.Set
}

func (l *loader) Load(fsys fs.FS) (state *State, err error) {
	state, err = l.load(fsys)
	if err != nil {
		return nil, fmt.Errorf("generator: %w", err)
	}
	return state, nil
}

func (l *loader) load(fsys fs.FS) (state *State, err error) {
	state = new(State)
	generatorDirs, err := finder.Find(fsys, l.Glob, func(path string, isDir bool) (entries []string) {
		if !isDir && valid.GoFile(path) && !isUserDefinedInternalGenerator(path) {
			entries = append(entries, filepath.Dir(path))
		}
		return entries
	})
	if err != nil {
		return nil, err
	}

	exists := make(map[string]bool)

	// Load passed in file generators
	for path, generator := range l.FileGenerators {
		if !exists[path] {
			exists[path] = true
			name := l.imports.Add(generator.Import)
			state.FileGenerators = append(state.FileGenerators, &CodeGenerator{
				Import: &imports.Import{
					Name: name,
					Path: generator.Import,
				},
				Path:  path,
				Camel: gotext.Camel(name),
			})
		}
	}

	// Load passed in file servers
	for path, generator := range l.FileServers {
		if !exists[path] {
			exists[path] = true
			name := l.imports.Add(generator.Import)
			state.FileServers = append(state.FileServers, &CodeGenerator{
				Import: &imports.Import{
					Name: name,
					Path: generator.Import,
				},
				Path:  path,
				Camel: gotext.Camel(name),
			})
		}
	}

	// Load the custom generators
	for _, generatorDir := range generatorDirs {
		importPath := l.module.Import(generatorDir)

		// Skip over conflicting generators
		if exists[importPath] {
			continue
		}

		// Parse the generator package
		pkg, err := l.parser.Parse(generatorDir)
		if err != nil {
			return nil, err
		}
		// Skip packages that don't have the expected signature
		generator := pkg.Struct("Generator")
		if generator == nil {
			l.log.Debug("framework/generator: skipping package because there's no Generator struct")
			continue
		}
		key := strings.TrimPrefix(generatorDir, "bud/internal/")
		key = strings.TrimPrefix(key, "generator/")

		// Support generating directories into the bud/internal directory
		if generator.Method("Generate") != nil {
			generatorPath := path.Join("bud", "internal", key)
			if !exists[generatorPath] {
				exists[generatorPath] = true
				name := l.imports.Add(importPath)
				state.GenerateDirs = append(state.GenerateDirs, &CodeGenerator{
					Import: &imports.Import{
						Name: name,
						Path: importPath,
					},
					Path:  generatorPath,
					Camel: gotext.Camel(name),
				})
			}
		}

		// Support serving files from the bud/internal directory
		if generator.Method("Serve") != nil {
			generatorPath := path.Join("bud", "internal", key)
			if !exists[generatorPath] {
				exists[generatorPath] = true
				name := l.imports.Add(importPath)
				state.ServeFiles = append(state.ServeFiles, &CodeGenerator{
					Import: &imports.Import{
						Name: name,
						Path: importPath,
					},
					Path:  generatorPath,
					Camel: gotext.Camel(name),
				})
			}
		}

		// Support generating directories into the bud/cmd directory
		if generator.Method("GenerateCmd") != nil {
			generatorPath := path.Join("bud", "cmd", key)
			if !exists[generatorPath] {
				exists[importPath] = true
				name := l.imports.Add(importPath)
				state.GenerateDirs = append(state.GenerateDirs, &CodeGenerator{
					Import: &imports.Import{
						Name: name,
						Path: importPath,
					},
					Path:  generatorPath,
					Camel: gotext.Camel(name),
				})
			}
		}

		// Support generating directories into the bud/pkg directory
		if generator.Method("GeneratePkg") != nil {
			generatorPath := path.Join("bud", "pkg", key)
			if !exists[generatorPath] {
				exists[importPath] = true
				name := l.imports.Add(importPath)
				state.GenerateDirs = append(state.GenerateDirs, &CodeGenerator{
					Import: &imports.Import{
						Name: name,
						Path: importPath,
					},
					Path:  generatorPath,
					Camel: gotext.Camel(name),
				})
			}
		}
	}
	l.imports.AddStd("io/fs")
	l.imports.AddNamed("genfs", "github.com/livebud/bud/package/genfs")
	state.Imports = l.imports.List()
	return state, nil
}

// Ignore packages in $APP/generator/generator/... because they're meant for
// internal generators, not user-defined generators.
func isUserDefinedInternalGenerator(path string) bool {
	return strings.HasPrefix(path, "generator/generator/")
}
