package remotefs

import (
	"io"
	"io/fs"
	"net/rpc"

	"github.com/livebud/bud/internal/virtual"
)

func NewServer(fsys fs.FS, conn io.ReadWriteCloser) *Server {
	server := rpc.NewServer()
	server.RegisterName(name, &Server{fsys})
	server.ServeConn(conn)
	return &Server{fsys}
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
