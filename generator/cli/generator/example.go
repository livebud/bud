package generator

import "gitlab.com/mnm/bud/pkg/gen"

type GeneratorExample struct {
}

func (g *GeneratorExample) Generator(generator map[string]gen.Generator) {
	generator["main.go"] = gen.GenerateFile(g.mainGo)
	generator["main.go"] = gen.GenerateFile(g.mainGo)
}

func (g *GeneratorExample) mainGo(_ gen.F, file *gen.File) error {
	return nil
}
