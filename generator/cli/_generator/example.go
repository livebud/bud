package generator

import "gitlab.com/mnm/bud/pkg/gen"

type GeneratorExample struct {
}

// Move to runtime
type Loader struct {
}

func (l *Loader) File(path string, fn func(f gen.F, file *gen.File) error) {}

func (l *Loader) Dir(path string, fn func(f gen.F, dir *gen.Dir) error) {}

// Different generators must not have files within the same directory
// Each generator must occupy a unique directory or unique tree of directories
func (g *GeneratorExample) Load(generate *Loader) {
	generate.File("main.go", g.mainGo)
	generate.Dir("bud/model", g.model)
}

func (g *GeneratorExample) mainGo(_ gen.F, file *gen.File) error {
	return nil
}

func (g *GeneratorExample) model(_ gen.F, dir *gen.Dir) error {
	return nil
}
