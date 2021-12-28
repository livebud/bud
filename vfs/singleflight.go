package vfs

import (
	"io"
	"io/fs"
	"io/ioutil"

	"golang.org/x/sync/singleflight"
)

func SingleFlight(fsys fs.FS) fs.FS {
	return &singleFlight{fsys: fsys, cache: Memory{}}
}

type singleFlight struct {
	loader singleflight.Group
	cache  Memory
	fsys   fs.FS
}

func (s *singleFlight) Open(name string) (fs.File, error) {
	// TODO: support concurrency
	if _, ok := s.cache[name]; ok {
		return s.cache.Open(name)
	}
	value, err, _ := s.loader.Do(name, func() (interface{}, error) {
		file, err := s.fsys.Open(name)
		if err != nil {
			return nil, err
		}
		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}
		if stat.IsDir() {
			s.cache[name] = &File{
				ModTime: stat.ModTime(),
				Mode:    stat.Mode(),
				Sys:     stat.Sys(),
			}
			return file, nil
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		s.cache[name] = &File{
			Data:    data,
			ModTime: stat.ModTime(),
			Mode:    stat.Mode(),
			Sys:     stat.Sys(),
		}
		// Seek back to the beginning
		if s, ok := file.(io.Seeker); ok {
			if _, err := s.Seek(0, io.SeekStart); err != nil {
				return nil, err
			}
		}
		return file, nil
	})
	if err != nil {
		return nil, err
	}
	return value.(fs.File), nil
}

// ReadDir implements the fs.ReadDirFS to pass the capability down
func (s *singleFlight) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(s.fsys, name)
}

// type cache struct {
// 	mu sync.RWMutex

// 	files map[string]*fstest.MapFile
// }

// func (c *cache) Set(key string, file fs.File) {
// 	c.mu.Lock()
// 	c.files[key] = &fstest.MapFile{}
// 	c.mu.Unlock()
// }
// func (c *cache) Get(key string) *fstest.MapFile {
// 	c.mu.RLock()
// 	defer c.mu.RUnlock()
// 	return c.files[key]
// }
