package appserver

import (
	"context"
	"io/fs"

	"github.com/livebud/bud/internal/exe"
	"github.com/livebud/bud/internal/gobuild"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/log"
)

type Server struct {
	Builder *gobuild.Builder
	Starter *exe.Command
	Bus     pubsub.Client
	Log     log.Interface
	FS      fs.FS
}

func (s *Server) Serve(ctx context.Context) error {
}
