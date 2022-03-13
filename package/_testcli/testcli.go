package testcli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/package/command/expand"
)

func Load(ctx context.Context, dir string) (*CLI, error) {
	expander, err := expand.Load(ctx, dir)
	if err != nil {
		return nil, err
	}
	if err := expander.Run(context.Background()); err != nil {
		return nil, err
	}
	return &CLI{
		dir: dir,
		Env: map[string]string{
			"HOME":       os.Getenv("HOME"),
			"PATH":       os.Getenv("PATH"),
			"GOPATH":     os.Getenv("GOPATH"),
			"GOCACHE":    os.Getenv("GOCACHE"),
			"GOMODCACHE": testdir.ModCache(dir).Directory(),
			"NO_COLOR":   "1",
			// TODO: remove once we can write a sum file to the modcache
			"GOPRIVATE": "*",
		},
	}, nil
}

type CLI struct {
	dir string
	Env map[string]string
}

func (c *CLI) Run(args ...string) (stdout, stderr string, err error) {
	cmd := exec.Command("./bud/cli", args...)
	cmd.Dir = c.dir
	env := make([]string, len(c.Env))
	for key, value := range c.Env {
		env = append(env, key+"="+value)
	}
	cmd.Env = env
	stdo := new(bytes.Buffer)
	cmd.Stdout = stdo
	stde := new(bytes.Buffer)
	cmd.Stderr = stde
	err = cmd.Run()
	if isExitStatus(err) {
		err = fmt.Errorf("%w: %s", err, stde.String())
	}
	return stdo.String(), stde.String(), err
}

func isExitStatus(err error) bool {
	return err != nil && strings.Contains(err.Error(), "exit status ")
}
