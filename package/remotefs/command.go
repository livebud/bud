package remotefs

import (
	"context"
	"io/fs"

	"github.com/livebud/bud/internal/errs"
	"github.com/livebud/bud/internal/exe"
	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/package/socket"
)

const defaultPrefix = "BUD_REMOTEFS"

// Start a remote server and connect to it
func Start(ctx context.Context, exec *exe.Template, name string, args ...string) (*Process, error) {
	cmd := &Command{
		Dir:        exec.Dir,
		Stdin:      exec.Stdin,
		Stdout:     exec.Stdout,
		Stderr:     exec.Stderr,
		Env:        exec.Env,
		ExtraFiles: exec.ExtraFiles,
	}
	return cmd.Start(ctx, name, args...)
}

// Command helps you launch a remotefs server and connect to it with the
// remotefs client
// TODO: remove this approach
type Command exe.Command

func (c *Command) Start(ctx context.Context, name string, args ...string) (*Process, error) {
	var closer once.Closer
	// Listen on any available TCP port
	// TODO: support other ways to start the server
	// TODO: use an os.Pipe?
	ln, err := socket.Listen(":0")
	if err != nil {
		return nil, err
	}
	closer.Add(ln.Close)
	// Turn the listener into a file to be passed to the subprocess
	file, err := ln.File()
	if err != nil {
		err = errs.Join(err, closer.Close())
		return nil, err
	}
	closer.Add(file.Close)
	// Inject the file listener into the subprocess
	extrafile.Inject(&c.ExtraFiles, &c.Env, defaultPrefix, file)
	// Start the subprocess
	process, err := (*exe.Command)(c).Start(ctx, name, args...)
	if err != nil {
		err = errs.Join(err, closer.Close())
		return nil, err
	}
	closer.Add(process.Close)
	// Dial the subprocess and return a client
	addr := ln.Addr().String()
	client, err := Dial(ctx, addr)
	if err != nil {
		err = errs.Join(err, closer.Close())
		return nil, err
	}
	closer.Add(client.Close)
	// Return the process
	return &Process{client, &closer, process, addr}, nil
}

type Process struct {
	client  *Client
	closer  *once.Closer
	process *exe.Process
	addr    string
}

var _ fs.FS = (*Process)(nil)
var _ fs.ReadDirFS = (*Process)(nil)

func (p *Process) URL() string {
	return p.addr
}

func (p *Process) Open(name string) (fs.File, error) {
	return p.client.Open(name)
}

func (p *Process) ReadDir(name string) (des []fs.DirEntry, err error) {
	return p.client.ReadDir(name)
}

func (p *Process) Change(path ...string) error {
	return p.client.Change(path...)
}

func (p *Process) Close() error {
	if err := p.closer.Close(); err != nil {
		return err
	}
	return p.process.Wait()
}

// Restart the process
func (p *Process) Restart(ctx context.Context) error {
	process, err := p.process.Restart(ctx)
	if err != nil {
		return err
	}
	p.process = process
	return nil
}
