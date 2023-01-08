package dev

import (
	"context"
	"net"

	"github.com/livebud/bud/internal/cli/run2/config"
	"github.com/livebud/bud/internal/pubsub"

	"github.com/livebud/bud/package/js"
)

func New(cfg *config.Config, ps pubsub.Client) *Server {
	return &Server{cfg, ps}
}

type Server struct {
	cfg *config.Config
	ps  pubsub.Client
}

type Launcher struct {
	Config *config.Config
	Pubsub pubsub.Client

	Listener net.Listener
	VM       js.VM
}

func (l *Launcher) Launch(ctx context.Context) (*Process, error) {
	return &Process{}, nil
}

type Process struct {
}

var _ js.VM = (*Process)(nil)

func (p *Process) Script(path, script string) error {
	return nil
}

func (p *Process) Eval(path, expr string) (string, error) {
	return "", nil
}

func (p *Process) Close() error {
	return nil
}
