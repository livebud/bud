package budfs

import (
	"io/fs"

	"gitlab.com/mnm/bud/internal/fscache"
	"gitlab.com/mnm/bud/internal/pluginfs"
	"gitlab.com/mnm/bud/pkg/gen"
	"gitlab.com/mnm/bud/pkg/gomod"
)

type option struct {
	FSCache *fscache.Cache
}

type Option = func(*option)

func WithFSCache(fc *fscache.Cache) func(*option) {
	return func(option *option) {
		option.FSCache = fc
	}
}

func Load(module *gomod.Module, options ...Option) (*FileSystem, error) {
	opt := &option{
		FSCache: nil,
	}
	for _, option := range options {
		option(opt)
	}
	plugin, err := pluginfs.Load(module, pluginfs.WithFSCache(opt.FSCache))
	if err != nil {
		return nil, err
	}
	genfs := gen.New(plugin, gen.WithFSCache(opt.FSCache))
	return &FileSystem{
		gen: genfs,
	}, nil
}

type FS interface {
	fs.FS
}

type FileSystem struct {
	gen *gen.FileSystem
}

func (f *FileSystem) Open(name string) (fs.File, error) {
	return f.gen.Open(name)
}

func (f *FileSystem) Entry(name string, generator gen.Generator) {
	f.gen.Add(map[string]gen.Generator{
		name: generator,
	})
}
