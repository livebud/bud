package config

import (
	"io"
	"os"
	"strconv"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/shell"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/levelfilter"
)

func New() *Config {
	return &Config{
		Log:       "info",
		Embed:     false,
		Minify:    false,
		Hot:       false,
		Stderr:    os.Stderr,
		Stdout:    os.Stdout,
		Stdin:     os.Stdin,
		Env:       os.Environ(),
		ListenDev: ":35729",
		ListenApp: ":3000",
		ListenAFS: ":0",
	}
}

type Config struct {
	Log string
	Dir string

	Embed  bool
	Minify bool
	Hot    bool

	Stderr io.Writer
	Stdout io.Writer
	Stdin  io.Reader
	Env    []string

	ListenDev string
	ListenApp string
	ListenAFS string
}

func (c *Config) Flag() *framework.Flag {
	return &framework.Flag{
		Embed:  c.Embed,
		Minify: c.Minify,
		Hot:    c.Hot,
	}
}

func (c *Config) Flags() []string {
	return []string{
		"--embed=" + strconv.FormatBool(c.Embed),
		"--minify=" + strconv.FormatBool(c.Minify),
		"--hot=" + strconv.FormatBool(c.Hot),
	}
}

func (c *Config) Command() *shell.Command {
	return &shell.Command{
		Dir:    c.Dir,
		Stderr: c.Stderr,
		Stdout: c.Stdout,
		Stdin:  c.Stdin,
		Env:    c.Env,
	}
}

func (c *Config) Module() (*gomod.Module, error) {
	return gomod.Find(c.Dir)
}

func (c *Config) Logger() (log.Log, error) {
	logLevel, err := log.ParseLevel(c.Log)
	if err != nil {
		return nil, err
	}
	log := log.New(levelfilter.New(console.New(c.Stderr), logLevel))
	return log, nil
}
