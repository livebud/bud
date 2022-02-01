package generate

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"gitlab.com/mnm/bud/internal/fscache"
	"gitlab.com/mnm/bud/pkg/socket"

	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/pkg/commander"
	"gitlab.com/mnm/bud/pkg/gen"
	"gitlab.com/mnm/bud/pkg/generator"
)

type Command struct {
	Dir        string
	Port       int
	Embed      bool
	Hot        bool
	Minify     bool
	Generators map[string]gen.Generator
	Args       []string
}

func (c *Command) Parse() error {
	cli := commander.New("bud")
	cli.Run(c.Run)
	return cli.Parse([]string{})
}

func (c *Command) Run(ctx context.Context) error {
	// TODO: Use the passed in generators (c.Generators)
	fsCache := fscache.New()
	generator, err := generator.Load(c.Dir, generator.WithFSCache(fsCache))
	if err != nil {
		return err
	}
	if err := generator.Generate(ctx); err != nil {
		return err
	}
	// Load the socket up, this should come from LISTENER_FDS
	listener, err := socket.Load(":3000")
	if err != nil {
		return err
	}
	mainPath := filepath.Join(c.Dir, "bud", "main.go")
	// Check to see if we generated a main.go
	if _, err := os.Stat(mainPath); err != nil {
		return err
	}
	// Build into bud/main
	binPath := filepath.Join(c.Dir, "bud", "main")
	if err := gobin.Build(ctx, c.Dir, mainPath, binPath); err != nil {
		return err
	}
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		// TODO: improve the welcome server
		return http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Welcome Server!\n"))
		}))
	}

	files, env, err := socket.Files(listener)
	if err != nil {
		return err
	}
	// Run the app
	cmd := exec.Command(binPath, c.Args...)
	cmd.Env = append(os.Environ(), string(env))
	cmd.ExtraFiles = files
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = c.Dir
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
