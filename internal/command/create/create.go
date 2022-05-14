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
	// Check if we can write into the directory
	if err := checkDir(dir); err != nil {
		return err
	}
	tmpDir, err := ioutil.TempDir("", "bud-create-*")
	if err != nil {
		return err
	}

	err = c.generateGoMod(ctx, tmpDir)
	if err != nil {
		return err
	}

	eg, ctx2 := errgroup.WithContext(ctx)
	eg.Go(func() error { return c.generateGitIgnore(ctx2, tmpDir) })
	eg.Go(func() error { return c.generatePackageJSON(ctx2, tmpDir, filepath.Base(dir)) })
	if err := eg.Wait(); err != nil {
		return err
	}
	// Create the project directory
	if err := os.MkdirAll(filepath.Dir(dir), 0755); err != nil {
		return err
	}
	// Try moving the temporary build path to the project directory
	if err := os.Rename(tmpDir, dir); err != nil {
		// Can't rename on top of an existing directory
		if !errors.Is(err, fs.ErrExist) {
			return err
		}
		// Move inner files over
		fis, err := os.ReadDir(tmpDir)
		if err != nil {
			return err
		}
		for _, fi := range fis {
			if err := os.Rename(filepath.Join(tmpDir, fi.Name()), filepath.Join(dir, fi.Name())); err != nil {
				return err
			}
		}
	}
	// TODO: clean this mess up.
	// It's breaking out of the packagejson.go file, but moving symlinks via
	// os.Rename doesn't seem to work.
	if version.Bud == "latest" {
		npm, err := exec.LookPath("npm")
		if err != nil {
			return err
		}
		budDir, err := findBudDir()
		if err != nil {
			return err
		}
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

func checkDir(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		// If it doesn't exist, treat it as empty
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		// All other errors should cause a failure
		return err
	}
	// Check if to see if the directory is empty
	fis, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(fis) > 0 {
		return fmt.Errorf("%q must be empty", dir)
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

func findBudDir() (string, error) {
	currentDir, err := dirname()
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(currentDir) {
		// Attempt to find it within $GOPATH/src
		if gopath := os.Getenv("GOPATH"); gopath != "" {
			currentDir = filepath.Join(gopath, "src", currentDir)
		}
	}
	return gomod.Absolute(currentDir)
}

func findBudModule() (*gomod.Module, error) {
	dir, err := findBudDir()
	if err != nil {
		return nil, err
	}
	return gomod.Find(dir)
}
