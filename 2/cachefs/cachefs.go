package cachefs

import (
	"io"
	"io/fs"
	"testing/fstest"

	"gitlab.com/mnm/bud/2/singleflight"
)

func New(fsys fs.FS, loader *singleflight.Loader, cache *Store) *FS {
	return &FS{fsys: fsys, loader: loader, cache: cache}
}

type FS struct {
	fsys   fs.FS
	loader *singleflight.Loader
	cache  *Store
}

func (f *FS) Open(name string) (fs.File, error) {
	// Try reading from the cache
	if f.cache.Has(name) {
		return f.cache.Open(name)
	}
	// Load the resource
	file, err := f.loader.Load(f.fsys, name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// Get the stats
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		f.cache.Set(name, &fstest.MapFile{
			ModTime: stat.ModTime(),
			Mode:    stat.Mode(),
			Sys:     stat.Sys(),
		})
		return file, nil
	}
	// Read the data fully
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	f.cache.Set(name, &fstest.MapFile{
		Data:    data,
		ModTime: stat.ModTime(),
		Mode:    stat.Mode(),
		Sys:     stat.Sys(),
	})
	// Seek back to the beginning
	if s, ok := file.(io.Seeker); ok {
		if _, err := s.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
	}
	// Open from the cache
	cached, err := f.cache.Open(name)
	if err != nil {
		return nil, err
	}
	return cached, nil
}

// Events
