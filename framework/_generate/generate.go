package generate

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/overlay"
)

//go:embed main.gotext
var template string

var generator = gotemplate.MustParse("framework/generate/main.gotext", template)

func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func New(module *gomod.Module) *Generator {
	return &Generator{module}
}

type Generator struct {
	module *gomod.Module
}

func (g *Generator) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	state, err := Load(fsys, g.module)
	if err != nil {
		return err
	}
	code, err := Generate(state)
	if err != nil {
		return err
	}
	fmt.Println(string(code))
	file.Data = code
	return nil
}

// func (g *Generator) GenerateDir(ctx context.Context, fsys overlay.F, dir *overlay.Dir) error {
// 	generators, err := g.findGenerators(fsys)
// 	if err != nil {
// 		return err
// 	}
// 	dir.FileGenerator("main.go", &generate{generators})
// 	for _, generator := range generators {
// 		dir.GenerateDir(generator, func(ctx context.Context, fsys overlay.F, dir *overlay.Dir) error {
// 			fmt.Println("generating...", dir.Path())
// 			return nil
// 		})
// 	}
// 	// if dir.Path() != "bud/internal/generator" {
// 	// 	return nil
// 	// }
// 	// fmt.Println(dir.Path())
// 	// state, err := Load(g.module)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// code, err := Generate(state)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// file.Data = code
// 	return nil
// }

// func (g *Generator) findGenerators(fsys fs.FS) (generators []string, err error) {
// 	err = fs.WalkDir(fsys, "generator", fs.WalkDirFunc(func(path string, de fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		} else if !de.IsDir() || !valid.Dir(path) {
// 			return nil
// 		} else if path == "generator" {
// 			// TODO: support the root directory
// 			return nil
// 		}
// 		generators = append(generators, path)
// 		return nil
// 	}))
// 	return generators, nil
// }

// // main.go

// type generate struct {
// 	generators []string
// }
