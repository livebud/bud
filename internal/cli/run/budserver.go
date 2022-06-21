package run

import (
	"context"
	"io/fs"
	"net"

	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/devserver"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/runtime/web"
)

type budServer struct {
	budln net.Listener
	bus   pubsub.Client
	fsys  fs.FS
	log   log.Interface
}

// Run the bud server
func (s *budServer) Run(ctx context.Context) error {
	vm, err := v8.Load()
	if err != nil {
		return err
	}
	devServer := devserver.New(s.fsys, s.bus, s.log, vm)
	err = web.Serve(ctx, s.budln, devServer)
	s.log.Debug("run: bud server closed", "err", err)
	return err
}
