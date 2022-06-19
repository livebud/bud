package toolfsls

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"

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
	dir := path.Clean(c.Dir)
	des, err := fs.ReadDir(fsys, dir)
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
