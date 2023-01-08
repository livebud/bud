package afs

import (
	"context"
	"io/fs"
	"net"

	"github.com/livebud/bud/package/genfs"

	"github.com/livebud/bud/internal/cli/run2/config"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/js"
	"github.com/livebud/bud/package/log"
)

type Launcher struct {
	Cache    genfs.Cache
	Config   *config.Config
	Listener net.Listener
	Log      log.Log
	Module   *gomod.Module
	VM       js.VM
}

func (l *Launcher) Generate(ctx context.Context, dirs ...string) error {
	return nil
}

func (l *Launcher) Launch(ctx context.Context) (*Process, error) {
	return &Process{}, nil
}

type Process struct {
}

func (p *Process) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

func (p *Process) ReadDir(name string) ([]fs.DirEntry, error) {
	return nil, fs.ErrNotExist
}

func (p *Process) Close() error {
	return nil
}

// func New(cfg *config.Config, cache *dag.DB, dev dev.Client) *Server {
// 	return &Server{}
// }

// type Server struct {
// }

// func (a *Server) Generate(ctx context.Context, dirs ...string) error {
// 	return nil
// }

// func (s *Server) Serve(ln net.Listener) error {
// 	return nil
// }

// func (s *Server) Open(name string) (fs.File, error) {
// 	return nil, fs.ErrNotExist
// }

// func (s *Server) ReadDir(name string) ([]fs.DirEntry, error) {
// 	return nil, fs.ErrNotExist
// }
