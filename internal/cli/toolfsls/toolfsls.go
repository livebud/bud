package toolfsls

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"

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
	Dir  string
}

func (c *Command) Run(ctx context.Context) error {
	log, err := bud.Log(c.in.Stdout, c.bud.Log)
	if err != nil {
		return err
	}
	dir := path.Clean(c.Dir)
	module, err := bud.Module(path.Join(c.bud.Dir, dir))
	if err != nil {
		return err
	}
	bfs, err := bfs.Load(c.Flag, log, module)
	if err != nil {
		return err
	}
	defer bfs.Close()
	des, err := fs.ReadDir(bfs, dir)
	if err != nil {
		return err
	}
	// Directories come first
	sort.Slice(des, func(i, j int) bool {
		if des[i].IsDir() && !des[j].IsDir() {
			return true
		} else if !des[i].IsDir() && des[j].IsDir() {
			return false
		}
		return des[i].Name() < des[j].Name()
	})
	// Print out list
	for _, de := range des {
		name := de.Name()
		if de.IsDir() {
			name += "/"
		}
		fmt.Fprintln(os.Stdout, name)
	}
	return nil
}
