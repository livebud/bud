package remotefs

import (
	"io"
	"io/fs"

	"github.com/livebud/bud/internal/virtual"
)

func NewService(fsys fs.FS) *Service {
	return &Service{fsys}
}

type Service struct {
	fsys fs.FS
}

func (s *Service) Open(path string, vfile *fs.File) error {
	file, err := s.fsys.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.IsDir() {
		des, err := fs.ReadDir(s.fsys, path)
		if err != nil {
			return err
		}
		entries := make([]fs.DirEntry, len(des))
		for i, de := range des {
			entries[i] = &virtual.DirEntry{
				Path:    de.Name(),
				ModeDir: de.IsDir(),
			}
		}
		// Return a directory
		*vfile = &virtual.Dir{
			Path:    path,
			Entries: entries,
		}
		return nil
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	*vfile = &virtual.File{
		Path: path,
		Data: data,
	}
	return nil
}

func (s *Service) ReadDir(name string, vdes *[]fs.DirEntry) error {
	des, err := fs.ReadDir(s.fsys, name)
	if err != nil {
		return err
	}
	for _, de := range des {
		*vdes = append(*vdes, &virtual.DirEntry{
			Path:    de.Name(),
			ModeDir: de.IsDir(),
		})
	}
	return nil
}
