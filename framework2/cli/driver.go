package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gitlab.com/mnm/bud/framework/clic"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/pkg/gomod"
)

// Compile the CLI driver
func Compile(ctx context.Context, dir string) (*Driver, error) {
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	compiler := clic.New(module)
	if err := compiler.Compile(ctx); err != nil {
		return nil, err
	}
	return &Driver{
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

type Driver struct {
	dir string
	Env map[string]string
}

func (d *Driver) Run(args ...string) (stdout, stderr string, err error) {
	cmd := exec.Command("./bud/cli", args...)
	cmd.Dir = d.dir
	env := make([]string, len(d.Env))
	for key, value := range d.Env {
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
