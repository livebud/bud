package remotefs

import (
	"io"
	"io/fs"
	"net/rpc"

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

func (s *Server) Open(name string, vfile **virtual.File) error {
	file, err := s.fsys.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	*vfile = &virtual.File{
		Name:    name,
		Data:    data,
		ModTime: stat.ModTime(),
		Mode:    stat.Mode(),
		Sys:     stat.Sys(),
	}
	return nil
}

func (s *Server) ReadDir(name string, vdes *[]*virtual.DirEntry) error {
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
			Base:    de.Name(),
			Mode:    stat.Mode(),
			ModTime: stat.ModTime(),
			Sys:     stat.Sys(),
			Size:    stat.Size(),
		})
	}
	return nil
}
