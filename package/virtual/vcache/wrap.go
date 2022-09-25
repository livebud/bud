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
		return virtual.New(entry), nil
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
	return virtual.New(entry), nil
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
		}, nil
	}
	vdir := &virtual.Dir{
		Path:    fullpath,
		ModTime: stat.ModTime(),
		Mode:    stat.Mode(),
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
		entryPath := path.Join(fullpath, de.Name())
		vdir.Entries = append(vdir.Entries, &dirEntry{f, entryPath, de})
	}
	return vdir, nil
}

// cached fs.DirEntry
type dirEntry struct {
	f         *cachedfs
	entryPath string
	de        fs.DirEntry
}

var _ fs.DirEntry = (*dirEntry)(nil)

func (e *dirEntry) Name() string {
	return e.de.Name()
}

func (e *dirEntry) IsDir() bool {
	return e.de.IsDir()
}

func (e *dirEntry) Type() fs.FileMode {
	return e.de.Type()
}

func (e *dirEntry) Info() (fs.FileInfo, error) {
	e.f.log.Debug("vcache: entry info", "name", e.entryPath)
	entry, ok := e.f.cache.Get(e.entryPath)
	if ok {
		e.f.log.Debug("vcache: cache hit", "name", e.entryPath)
		return virtual.New(entry).Stat()
	}
	e.f.log.Debug("vcache: cache miss", "name", e.entryPath)
	return e.de.Info()
}
