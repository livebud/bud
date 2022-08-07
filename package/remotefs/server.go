package remotefs

import (
	"io"
	"io/fs"
	"net/rpc"
	"path"

	"github.com/livebud/bud/internal/virtual"
)

func Serve(fsys fs.FS, conn io.ReadWriteCloser) {
	server := rpc.NewServer()
	server.RegisterName("remotefs", &Server{fsys})
	server.ServeConn(conn)
}

type Server struct {
	fsys fs.FS
}

func (s *Server) Open(name string, vfile *fs.File) error {
	file, err := s.fsys.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.IsDir() {
		des, err := fs.ReadDir(s.fsys, name)
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
			Name:    path.Base(name),
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
		Name:    path.Base(name),
		Data:    data,
		ModTime: stat.ModTime(),
		Mode:    stat.Mode(),
	}
	return nil
}

func (s *Server) ReadDir(name string, vdes *[]fs.DirEntry) error {
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
			Sys:     stat.Sys(),
			Size:    stat.Size(),
		})
	}
	return nil
}
