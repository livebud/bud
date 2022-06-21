package toolfscat

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/cli/bud"
)

func New(bud *bud.Command) *Command {
	return &Command{
		bud:  bud,
		Flag: new(framework.Flag),
	}
}

type Command struct {
	bud  *bud.Command
	Flag *framework.Flag
	Path string
}

func (c *Command) Run(ctx context.Context) error {
	log, err := c.bud.Logger()
	if err != nil {
		return err
	}
	module, err := c.bud.Module()
	if err != nil {
		return err
	}
	fsys, err := c.bud.FileSystem(log, module, c.Flag)
	if err != nil {
		return err
	}
	code, err := fs.ReadFile(fsys, path.Clean(c.Path))
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, string(code))
	return nil
}
