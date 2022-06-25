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

func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{
		bud:  bud,
		in:   in,
		Flag: new(framework.Flag),
	}
}

type Command struct {
	bud  *bud.Command
	in   *bud.Input
	Flag *framework.Flag
	Dir  string
}

func (c *Command) Run(ctx context.Context) error {
	log, err := bud.Log(c.in.Stderr, c.bud.Log)
	if err != nil {
		return err
	}
	module, err := bud.Module(c.Dir)
	if err != nil {
		return err
	}
	fsys, err := bud.FileSystem(log, module, c.Flag)
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
