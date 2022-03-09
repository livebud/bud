package maing

import (
	"context"
	_ "embed"

	"gitlab.com/mnm/bud/internal/imports"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/package/overlay"
)

//go:embed main.gotext
var template string

var generator = gotemplate.MustParse("main.gotext", template)

type state struct {
	Imports []*imports.Import
}

// func New

func New(programPath string) overlay.GenerateFile {
	return func(ctx context.Context, fsys overlay.F, file *overlay.File) error {
		// Add default imports
		imports := imports.New()
		imports.AddStd("os", "context")
		imports.AddNamed("program", programPath)
		// Generate code
		code, err := generator.Generate(state{imports.List()})
		if err != nil {
			return err
		}
		// Write code to file
		file.Data = code
		return nil
	}
}
