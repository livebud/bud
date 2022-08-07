package remotefs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/rpc"
	"path"

	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/internal/virtual"
	"github.com/livebud/bud/package/socket"
)

// ServeFrom serves the filesystem from a listener passed in by a parent process
func ServeFrom(ctx context.Context, fsys fs.FS, prefix string) error {
	files := extrafile.Load(prefix)
	if len(files) == 0 {
		return fmt.Errorf("remotefs: no extra files passed into the process")
	}
	ln, err := socket.From(files[0])
	if err != nil {
		return fmt.Errorf("remotefs: unable to turn extra file into listener. %w", err)
	}
	defer ln.Close()
	go Serve(fsys, ln)
	<-ctx.Done()
	return nil
}

// Serve the filesystem from a listener
func Serve(fsys fs.FS, ln net.Listener) error {
	server := rpc.NewServer()
	server.RegisterName("remotefs", &Server{fsys})
	return accept(server, ln)
}

// Accept connections from the listener. This will block until the listener is
// closed
func accept(server *rpc.Server, ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		go server.ServeConn(conn)
	}
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
