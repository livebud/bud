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
	"runtime"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/command"
	"github.com/livebud/bud/internal/version"
	"github.com/livebud/bud/package/gomod"
)

type Command struct {
	Bud *command.Bud
	Dir string
}

func (c *Command) Run(ctx context.Context) error {
	dir := filepath.Join(c.Bud.Dir, c.Dir)
	if _, err := os.Stat(dir); nil == err {
		return fmt.Errorf("%s already exists", dir)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	tmpDir, err := ioutil.TempDir("", "bud-create-*")
	if err != nil {
		return err
	}
	eg, ctx2 := errgroup.WithContext(ctx)
	eg.Go(func() error { return c.generatePackageJSON(ctx2, tmpDir, filepath.Base(dir)) })
	eg.Go(func() error { return c.generateGitIgnore(ctx2, tmpDir) })
	eg.Go(func() error { return c.generateGoMod(ctx2, tmpDir) })
	if err := eg.Wait(); err != nil {
		return err
	}
	// Create the project directory
	if err := os.MkdirAll(filepath.Dir(dir), 0755); err != nil {
		return err
	}
	// Move the temporary build path to the project directory
	if err := os.Rename(tmpDir, dir); err != nil {
		return err
	}
	// TODO: clean this mess up.
	// It's breaking out of the packagejson.go file, but moving symlinks via
	// os.Rename doesn't seem to work.
	if version.Bud == "latest" {
		npm, err := exec.LookPath("npm")
		if err != nil {
			return err
		}
		currentDir, err := dirname()
		if err != nil {
			return err
		}
		budDir, err := gomod.Absolute(currentDir)
		if err != nil {
			return err
		}
		fmt.Println(dir, "link", "livebud", filepath.Join(budDir, "livebud"))
		cmd := exec.CommandContext(ctx, npm, "link", "--loglevel=error", "livebud", filepath.Join(budDir, "livebud"))
		cmd.Dir = dir
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

// dirname gets the directory of this file
func dirname() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filepath.Dir(filename), nil
}
