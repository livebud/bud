package vcache

import (
	"fmt"
	"io"
	"io/fs"
	"path"

	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/virtual"
)

func Wrap(cache Cache, fsys fs.FS, log log.Interface) fs.FS {
	return &cachedfs{cache, fsys, log}
}

type cachedfs struct {
	cache Cache
	fsys  fs.FS
	log   log.Interface
}

func (f *cachedfs) Open(name string) (fs.File, error) {
	f.log.Debug("vcache: open", "name", name)
	entry, ok := f.cache.Get(name)
	if ok {
		f.log.Debug("vcache: cache hit", "name", name)
		return entry.Open(), nil
	}
	f.log.Debug("vcache: cache miss", "name", name)
	file, err := f.fsys.Open(name)
	if err != nil {
		return nil, err
	}
	entry, err = f.toEntry(name, file)
	if err != nil {
		return nil, err
	}
	f.cache.Set(name, entry)
	return entry.Open(), nil
}

// toEntry converts a fs.File into a reusable virtual entry
func (f *cachedfs) toEntry(fullpath string, file fs.File) (virtual.Entry, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	// Handle files
	if !stat.IsDir() {
		// Read the data fully
		data, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		return &virtual.File{
			Path:    fullpath,
			Data:    data,
			ModTime: stat.ModTime(),
			Mode:    stat.Mode(),
			Sys:     stat.Sys(),
		}, nil
	}
	vdir := &virtual.Dir{
		Path:    fullpath,
		ModTime: stat.ModTime(),
		Mode:    stat.Mode(),
		Sys:     stat.Sys(),
	}
	dir, ok := file.(fs.ReadDirFile)
	if !ok {
		return nil, fmt.Errorf("vcache: dir does not have ReadDir: %s", fullpath)
	}
	des, err := dir.ReadDir(-1)
	if err != nil {
		return nil, err
	}
	for _, de := range des {
		stat, err := f.toEntryInfo(fullpath, de)
		if err != nil {
			return nil, err
		}
		vdir.Entries = append(vdir.Entries, &virtual.DirEntry{
			Path:    de.Name(),
			ModTime: stat.ModTime(),
			Mode:    stat.Mode(),
			Sys:     stat.Sys(),
			Size:    stat.Size(),
		})
	}
	return vdir, nil
}

func (f *cachedfs) toEntryInfo(fullpath string, de fs.DirEntry) (fs.FileInfo, error) {
	entryPath := path.Join(fullpath, de.Name())
	f.log.Debug("vcache: entry info", "name", entryPath)
	entry, ok := f.cache.Get(entryPath)
	if ok {
		f.log.Debug("vcache: cache hit", "name", entryPath)
		return entry.Open().Stat()
	}
	f.log.Debug("vcache: cache miss", "name", entryPath)
	return de.Info()
}
