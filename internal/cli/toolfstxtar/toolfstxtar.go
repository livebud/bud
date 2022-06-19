package toolfstxtar

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"

	"golang.org/x/tools/txtar"

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
	Dir  string
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
	ar := new(txtar.Archive)
	dir := path.Clean(c.Dir)
	err = fs.WalkDir(fsys, dir, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if de.IsDir() {
			return nil
		}
		code, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		ar.Files = append(ar.Files, txtar.File{
			Name: path,
			Data: code,
		})
		return nil
	})
	if err != nil {
		return err
	}
	// Print the archive to stdout
	fmt.Fprintln(os.Stdout, string(txtar.Format(ar)))
	return nil
}
