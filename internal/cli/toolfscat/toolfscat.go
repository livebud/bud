package toolfscat

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/bfs"
	"github.com/livebud/bud/internal/cli/bud"
)

func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{
		bud: bud,
		in:  in,
		Flag: &framework.Flag{
			Env:    in.Env,
			Stderr: in.Stderr,
			Stdin:  in.Stdin,
			Stdout: in.Stdout,
		},
	}
}

type Command struct {
	bud  *bud.Command
	in   *bud.Input
	Flag *framework.Flag
	Path string
}

func (c *Command) Run(ctx context.Context) error {
	log, err := bud.Log(c.in.Stderr, c.bud.Log)
	if err != nil {
		return err
	}
	module, err := bud.Module(c.bud.Dir)
	if err != nil {
		return err
	}
	bfs, err := bfs.Load(c.Flag, log, module)
	if err != nil {
		return err
	}
	defer bfs.Close()
	code, err := fs.ReadFile(bfs, path.Clean(c.Path))
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, string(code))
	return nil
}
