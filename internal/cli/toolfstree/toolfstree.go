package toolfstree

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/livebud/bud/internal/printfs"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/cli/bud"
)

func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{
		bud: bud,
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
	Flag *framework.Flag
	Dir  string
}

func (c *Command) Run(ctx context.Context) error {
	log, err := bud.Log(c.Flag.Stdout, c.bud.Log)
	if err != nil {
		return err
	}
	dir := path.Clean(c.Dir)
	module, err := bud.Module(path.Join(c.bud.Dir, dir))
	if err != nil {
		return err
	}
	bfs, err := bud.FileSystem(ctx, log, module, c.Flag)
	if err != nil {
		return err
	}
	defer bfs.Close()
	tree, err := printfs.Print(bfs, dir)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, tree)
	return nil
}
