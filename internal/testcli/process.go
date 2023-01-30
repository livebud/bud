package testcli

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/livebud/bud/internal/extrafile"

	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/internal/shell"
)

// StartApp starts bud/app. It's meant to be used after running `bud build`.
// TODO: integrate better with testcli
func (c *CLI) StartApp(ctx context.Context, args ...string) (*Process, error) {
	webLn, webc, err := listen(":0")
	if err != nil {
		return nil, err
	}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd := shell.Command{
		Dir:    c.dir,
		Env:    c.Env.List(),
		Stdin:  c.Stdin,
		Stdout: io.MultiWriter(os.Stdout, stdout),
		Stderr: io.MultiWriter(os.Stderr, stderr),
	}
	closer := new(once.Closer)
	webFile, err := webLn.File()
	if err != nil {
		return nil, err
	}
	closer.Add(webFile.Close)
	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "WEB", webFile)
	args = append(args, "--listen", webLn.Addr().String())
	process, err := cmd.Start(ctx, filepath.Join("bud", "app"), prependFlags(args)...)
	if err != nil {
		return nil, err
	}
	closer.Add(process.Close)
	return &Process{
		closer: closer,
		stdout: stdout,
		stderr: stderr,
		webc:   webc,
	}, nil
}

type Process struct {
	closer *once.Closer
	stdout *bytes.Buffer
	stderr *bytes.Buffer
	webc   *http.Client
}

func (p *Process) Close() error {
	return p.closer.Close()
}

func (p *Process) Get(path string) (*Response, error) {
	req, err := getRequest(path)
	if err != nil {
		return nil, err
	}
	return do(p.webc, req)
}
