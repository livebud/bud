package budsvr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/rpc"

	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/budrpc"
	"github.com/livebud/bud/package/js"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/socket"
)

func New(bus pubsub.Client, fsys fs.FS, log log.Interface, vm js.VM) *Server {
	rpc := rpc.NewServer()
	rpc.RegisterName("bud", &Service{bus, fsys, log, vm})
	rpc.RegisterName("remotefs", remotefs.NewService(fsys))
	return &Server{rpc}
}

type Server struct {
	rpc *rpc.Server
}

const defaultPrefix = "BUD_RPC"

// ServeFrom serves the filesystem from a listener passed in by a parent process
func (s *Server) ServeFrom(prefix string) error {
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
	return s.Serve(ln)
}

// Serve the filesystem from a listener
func (s *Server) Serve(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		go s.rpc.ServeConn(conn)
	}
}

type Service struct {
	bus  pubsub.Client
	fsys fs.FS
	log  log.Interface
	vm   js.VM
}

func (s *Service) Render(req budrpc.RenderRequest, res *budrpc.RenderResponse) error {
	script, err := fs.ReadFile(s.fsys, "bud/view/_ssr.js")
	if err != nil {
		return err
	}
	expr := fmt.Sprintf(`%s; bud.render(%q, %s)`, script, req.Route, req.Props)
	result, err := s.vm.Eval("_ssr.js", expr)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(result), res); err != nil {
		return err
	}
	return nil
}

// Publish an event
func (s *Service) Publish(req budrpc.PublishRequest, res *struct{}) error {
	s.bus.Publish(req.Topic, req.Data)
	return nil
}
