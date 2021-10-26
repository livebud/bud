package maingo

import (
	_ "embed"

	"gitlab.com/mnm/bud/bfs"
	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed maingo.gotext
var maingoTemplate string

// maingoGenerator
var maingoGenerator = gotemplate.MustParse("maingo.gotext", maingoTemplate)

type State struct {
}

func Generator() bfs.Generator {
	return bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
		code, err := maingoGenerator.Generate(State{})
		if err != nil {
			return err
		}
		file.Write(code)
		return nil
	})
}
