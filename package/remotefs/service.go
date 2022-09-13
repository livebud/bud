package remotefs

import (
	"io"
	"io/fs"

	"github.com/livebud/bud/package/virtual"
)

func NewService(fsys fs.FS) *Service {
	return &Service{fsys}
}

type Service struct {
	fsys fs.FS
}

func (s *Service) Open(path string, vfile *virtual.Entry) error {
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
			fi, err := de.Info()
			if err != nil {
				return err
			}
			entries[i] = &virtual.DirEntry{
				Path:    de.Name(),
				Mode:    de.Type(),
				ModTime: fi.ModTime(),
				Size:    fi.Size(),
			}
		}
		// Return a directory
		*vfile = &virtual.Dir{
			Path:    path,
			ModTime: stat.ModTime(),
			Mode:    stat.Mode(),
			Entries: entries,
		}
		return nil
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	*vfile = &virtual.File{
		Path:    path,
		Data:    data,
		ModTime: stat.ModTime(),
		Mode:    stat.Mode(),
	}
	return nil
}

func (s *Service) ReadDir(name string, vdes *[]fs.DirEntry) error {
	des, err := fs.ReadDir(s.fsys, name)
	if err != nil {
		return err
	}
	for _, de := range des {
		stat, err := de.Info()
		if err != nil {
			return err
		}
		*vdes = append(*vdes, &virtual.DirEntry{
			Path:    de.Name(),
			Mode:    stat.Mode(),
			ModTime: stat.ModTime(),
			Size:    stat.Size(),
		})
	}
	return nil
}
