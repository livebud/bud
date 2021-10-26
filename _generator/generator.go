package generator

import (
	"os"

	"gitlab.com/mnm/bud/bfs"
	"gitlab.com/mnm/bud/go/mod"
)

func New(modfile mod.File) *Generator {
	return &Generator{modfile}
}

type Generator struct {
	modfile mod.File
}

type GenerateInput struct {
	Embed bool
}

func (g *Generator) Generate(in *GenerateInput) error {
	dirfs := os.DirFS(g.modfile.Directory())
	bf := bfs.New(dirfs)
	_ = bf
	return nil
}
