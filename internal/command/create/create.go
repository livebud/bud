package create

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/command"
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
	eg, ctx2 := errgroup.WithContext(ctx)
	eg.Go(func() error { return c.generatePackageJSON(ctx2, tmpDir, filepath.Base(c.Dir)) })
	eg.Go(func() error { return c.generateGitIgnore(ctx2, tmpDir) })
	eg.Go(func() error { return c.generateGoMod(ctx2, tmpDir) })
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
	// TODO: clean this mess up.
	// It's breaking out of the packagejson.go file, but moving symlinks doesn't
	// seem to work.
	if c.Link {
		npm, err := exec.LookPath("npm")
		if err != nil {
			return err
		}
		cmd := exec.CommandContext(ctx, npm, "link", "--loglevel=error", "livebud")
		cmd.Dir = c.Dir
		cmd.Stderr = os.Stderr
		cmd.Env = []string{
			"HOME=" + os.Getenv("HOME"),
			"PATH=" + os.Getenv("PATH"),
			"NO_COLOR=1",
		}
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
