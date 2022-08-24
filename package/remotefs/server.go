package remotefs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/rpc"

	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/package/socket"
)

// ServeFrom serves the filesystem from a listener passed in by a parent process
func ServeFrom(ctx context.Context, fsys fs.FS, prefix string) error {
	if prefix == "" {
		prefix = defaultPrefix
	}
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
	server.RegisterName("remotefs", NewService(fsys))
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
