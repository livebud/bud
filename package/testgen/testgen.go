package testgen

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/livebud/bud/package/gotemplate"

	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/virtual"

	"github.com/livebud/bud/internal/envs"
	"github.com/livebud/bud/internal/shell"
	"github.com/livebud/bud/package/gomod"
)

func New(genfs fs.FS, module *gomod.Module) *Command {
	log := testlog.New()
	return &Command{
		Env: envs.Map{
			"NO_COLOR": "1",
			"HOME":     os.Getenv("HOME"),
			"PATH":     os.Getenv("PATH"),
			"GOPATH":   os.Getenv("GOPATH"),
			"TMPDIR":   os.TempDir(),
		},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		genfs:  genfs,
		log:    log,
		module: module,
	}
}

type Command struct {
	Env    envs.Map
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	genfs  fs.FS
	log    log.Log
	module *gomod.Module
}

// copy generator files over
func (c *Command) copy() error {
	return virtual.Copy(c.log, c.genfs, c.module)
}

// copy generator files over
func (c *Command) build(ctx context.Context, mainFile, mainBin string) error {
	cmd := c.command()
	if err := cmd.Run(ctx, "go", "build", "-mod", "mod", "-o", mainBin, mainFile); err != nil {
		return err
	}
	return nil
}

func (c *Command) command() *shell.Command {
	return &shell.Command{
		Dir:    c.module.Directory(),
		Env:    c.Env.List(),
		Stdin:  c.Stdin,
		Stdout: c.Stdout,
		Stderr: c.Stderr,
	}
}

// Run a command
func (c *Command) Run(ctx context.Context, mainFile string, args ...string) (stdout *bytes.Buffer, stderr *bytes.Buffer, err error) {
	if err := c.copy(); err != nil {
		return nil, nil, err
	}
	const mainBin = "./main"
	if err := c.build(ctx, mainFile, mainBin); err != nil {
		return nil, nil, err
	}
	stdout = new(bytes.Buffer)
	stderr = new(bytes.Buffer)
	cmd := &shell.Command{
		Dir:    c.module.Directory(),
		Env:    c.Env.List(),
		Stdin:  c.Stdin,
		Stdout: io.MultiWriter(c.Stdout, stdout),
		Stderr: io.MultiWriter(c.Stderr, stderr),
	}
	if err := cmd.Run(ctx, mainBin, args...); err != nil {
		return nil, nil, err
	}
	return stdout, stderr, nil
}

// Start a process
func (c *Command) Start(ctx context.Context, mainFile string, args ...string) (*shell.Process, error) {
	if err := c.copy(); err != nil {
		return nil, err
	}
	mainBin := filepath.Join(c.module.Directory(), "main")
	if err := c.build(ctx, mainFile, mainBin); err != nil {
		return nil, err
	}
	cmd := &shell.Command{
		Dir:    c.module.Directory(),
		Env:    c.Env.List(),
		Stdin:  c.Stdin,
		Stdout: c.Stdout,
		Stderr: c.Stderr,
	}
	return cmd.Start(ctx, mainBin, args...)
}

func Main(template string, state interface{}) *mainGen {
	return &mainGen{template: template, state: state}
}

type mainGen struct {
	template string
	state    interface{}
}

func (m *mainGen) GenerateFile(fsys genfs.FS, file *genfs.File) error {
	generator, err := gotemplate.Parse("main.gotext", m.template)
	if err != nil {
		return err
	}
	code, err := generator.Generate(m.state)
	if err != nil {
		return err
	}
	file.Data = []byte(code)
	return nil
}
