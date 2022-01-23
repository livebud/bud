package budfs

import (
	"io/fs"

	"gitlab.com/mnm/bud/2/fscache"
	"gitlab.com/mnm/bud/2/genfs"
	"gitlab.com/mnm/bud/2/mod"
	"gitlab.com/mnm/bud/2/pluginfs"
	"gitlab.com/mnm/bud/internal/modcache"
)

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

func Load(fmap *fscache.Cache, module *mod.Module, options ...Option) (*FS, error) {
	plugin, err := pluginfs.Load(module, pluginfs.WithFSCache(fmap))
	if err != nil {
		return nil, err
	}
	// cache2 := cachefs.New(plugin, singleflight.New(), store)
	genfs := genfs.New(plugin, genfs.WithFSCache(fmap))
	// cache3 := cachefs.New(genfs, singleflight.New(), store)
	// Cache is what we should read from, but we also need access to the generator
	// filesystem to be able to add generators.
	return &FS{
		gen: genfs,
	}, nil
}

type FS struct {
	cache    *modcache.Cache
	replaces []*Replace
	gen      *genfs.FileSystem
}

func (f *FS) Open(name string) (fs.File, error) {
	return f.gen.Open(name)
}

func (f *FS) Entry(name string, generator genfs.Generator) {
	f.gen.Add(map[string]genfs.Generator{
		name: generator,
	})
}

// // Merge the filesystems into one
// func merge(filesystems ...fs.FS) fs.FS {
// 	switch len(filesystems) {
// 	case 0:
// 		return emptyFS{}
// 	case 1:
// 		return filesystems[0]
// 	default:
// 		var next fs.FS = filesystems[0]
// 		for _, plugin := range filesystems[1:] {
// 			next = mergefs.NewMergedFS(next, plugin)
// 		}
// 		return next
// 	}
// }

// type emptyFS struct {
// }

// func (emptyFS) Open(name string) (fs.File, error) {
// 	return nil, fs.ErrNotExist
// }
