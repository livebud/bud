package create

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	"gitlab.com/mnm/bud/internal/command"
)

type Command struct {
	Bud  *command.Bud
	Dir  string
	Link bool
}

func (c *Command) Run(ctx context.Context) error {
	if _, err := os.Stat(c.Dir); nil == err {
		return fmt.Errorf("%s already exists", c.Dir)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	tmpDir, err := ioutil.TempDir("", "bud-create-*")
	if err != nil {
		return err
	}
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { return c.generatePackageJSON(ctx, tmpDir, filepath.Base(c.Dir)) })
	eg.Go(func() error { return c.generateGitIgnore(ctx, tmpDir) })
	eg.Go(func() error { return c.generateGoMod(ctx, tmpDir) })
	if err := eg.Wait(); err != nil {
		return err
	}
	// Create the project directory
	if err := os.MkdirAll(filepath.Dir(c.Dir), 0755); err != nil {
		return err
	}
	// Move the temporary build path to the project directory
	if err := os.Rename(tmpDir, c.Dir); err != nil {
		return err
	}
	return nil
}
