package tester

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gitlab.com/mnm/bud/generator/cli"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/socket"
)

// Compile the CLI driver
func Compile(ctx context.Context, dir string) (*CLI, error) {
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	compiler := cli.New(module)
	if err := compiler.Compile(ctx); err != nil {
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

func (c *CLI) Run(args ...string) (*Process, error) {
	cmd := exec.Command("./bud/cli", "run")
	cmd.Dir = c.dir

	cmd.Stderr = os.Stdout
	cmd.Stderr = os.Stderr

	// Load the environment
	env := make([]string, len(c.Env))
	for key, value := range c.Env {
		env = append(env, key+"="+value)
	}

	// Since we're starting the web server, initialize a unix domain socket
	// to listen and pass that socket to the application process
	// Start the unix domain socket
	name, err := ioutil.TempDir("", "bud-testapp-*")
	if err != nil {
		return nil, err
	}
	socketPath := filepath.Join(name, "tmp.sock")
	// Heads up: If you see the `bind: invalid argument` error, there's a chance
	// the path is too long. 103 characters appears to be the limit on OSX,
	// https://github.com/golang/go/issues/6895.
	if len(socketPath) > 103 {
		return nil, fmt.Errorf("socket name is too long")
	}
	ln, err := socket.Listen(socketPath)
	if err != nil {
		return nil, err
	}
	transport, err := socket.Transport(socketPath)
	if err != nil {
		return nil, err
	}
	// Add socket configuration to the command
	files, socketEnv, err := socket.Files(ln)
	if err != nil {
		return nil, err
	}
	cmd.ExtraFiles = append(cmd.ExtraFiles, files...)
	env = append(env, string(socketEnv))

	// Set the env
	cmd.Env = env

	// Start the webserver
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Process{
		cmd: cmd,
		ln:  ln,
		client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}, nil
}

func (c *CLI) Build() (*App, error) {
	cmd := exec.Command("./bud/cli", "build")
	cmd.Dir = c.dir
	env := make([]string, len(c.Env))
	for key, value := range c.Env {
		env = append(env, key+"="+value)
	}
	cmd.Env = env
	cmd.Stderr = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return &App{c.dir, c.Env}, nil
}

// Command runs a custom command
func (c *CLI) Command(args ...string) (stdout, stderr string, err error) {
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
