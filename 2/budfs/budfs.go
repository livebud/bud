package budfs

import (
	"io/fs"

	"gitlab.com/mnm/bud/2/cachefs"
	"gitlab.com/mnm/bud/2/genfs"
	"gitlab.com/mnm/bud/2/mod"
	"gitlab.com/mnm/bud/2/pluginfs"
	"gitlab.com/mnm/bud/2/singleflight"
	"gitlab.com/mnm/bud/internal/modcache"
)

type FS struct {
	cache    *modcache.Cache
	replaces []*Replace
	gen      *genfs.FileSystem
	fsys     fs.FS
}

type Replace struct {
	Old string
	New string
}

type Option func(*FS)

func WithReplace(from, to string) func(*FS) {
	return func(f *FS) {
		f.replaces = append(f.replaces, &Replace{Old: from, New: to})
	}
}

func WithCache(mc *modcache.Cache) func(*FS) {
	return func(f *FS) {
		f.cache = mc
	}
}

func Load(module *mod.Module, options ...Option) (*FS, error) {
	store := cachefs.Cache()
	// loader := singleflight.New()
	cache1 := cachefs.New(module, singleflight.New(), store)
	plugin, err := pluginfs.Load(cache1, module)
	if err != nil {
		return nil, err
	}
	cache2 := cachefs.New(plugin, singleflight.New(), store)
	genfs := genfs.New(cache2)
	cache3 := cachefs.New(genfs, singleflight.New(), store)
	// Cache is what we should read from, but we also need access to the generator
	// filesystem to be able to add generators.
	return &FS{
		fsys: cache3,
		gen:  genfs,
	}, nil
}

func (f *FS) Open(name string) (fs.File, error) {
	return f.fsys.Open(name)
}

func (f *FS) Entry(name string, generator genfs.Generator) {
	f.gen.Add(map[string]genfs.Generator{
		name: generator,
	})
}
