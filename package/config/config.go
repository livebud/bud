package config

import (
	"io"
	"net"
	"os"
	"strconv"

	"github.com/livebud/bud/internal/pubsub"
)

func Default() *Config {
	return &Config{
		".",
		"info",
		false,
		false,
		false,
		":0",
		":35729",
		":3000",
		os.Stdin,
		os.Stdout,
		os.Stderr,
		os.Environ(),
		pubsub.New(),
		nil,
		nil,
		nil,
	}
}

type Config struct {
	Dir string
	Log string

	Embed  bool
	Hot    bool
	Minify bool

	ListenAFS string
	ListenDev string
	ListenWeb string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    []string

	// Used for testing
	Bus         pubsub.Client
	AFSListener net.Listener
	DevListener net.Listener
	WebListener net.Listener
}

func (c *Config) Flags() []string {
	return []string{
		"--log=" + strconv.Quote(c.Log),
		"--embed=" + strconv.FormatBool(c.Embed),
		"--minify=" + strconv.FormatBool(c.Minify),
		"--hot=" + strconv.FormatBool(c.Hot),
	}
}
