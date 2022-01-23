package budfs

import (
	"io/fs"

	"gitlab.com/mnm/bud/2/fscache"
	"gitlab.com/mnm/bud/2/genfs"
	"gitlab.com/mnm/bud/2/mod"
	"gitlab.com/mnm/bud/2/pluginfs"
)

func Load(fsCache *fscache.Cache, module *mod.Module) (*FS, error) {
	plugin, err := pluginfs.Load(module, pluginfs.WithFSCache(fsCache))
	if err != nil {
		return nil, err
	}
	genfs := genfs.New(plugin, genfs.WithFSCache(fsCache))
	return &FS{
		gen: genfs,
	}, nil
}

type FS struct {
	gen *genfs.FileSystem
}

func (f *FS) Open(name string) (fs.File, error) {
	return f.gen.Open(name)
}

func (f *FS) Entry(name string, generator genfs.Generator) {
	f.gen.Add(map[string]genfs.Generator{
		name: generator,
	})
}
