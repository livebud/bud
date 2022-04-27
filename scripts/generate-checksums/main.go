package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/livebud/bud/package/commander"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log/console"
)

func main() {
	if err := run(); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	cmd := new(Command)
	cli := commander.New("generate-checksums")
	cli.Run(cmd.Run)
	return cli.Parse(context.Background(), os.Args[1:])
}

type Command struct {
	Version string
}

type State struct {
	Notes string
}

func (c *Command) Run(ctx context.Context) error {
	module, err := gomod.Find(".")
	if err != nil {
		return err
	}
	paths, err := filepath.Glob(module.Directory("release", "*.tar.gz"))
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join("release", "checksums.txt"))
	if err != nil {
		return err
	}
	defer f.Close()
	for _, path := range paths {
		sha, err := computeSHA(path)
		if err != nil {
			return err
		}
		if _, err := f.Write([]byte(sha + "  " + filepath.Base(path) + "\n")); err != nil {
			return err
		}
	}
	return nil
}

func computeSHA(path string) (string, error) {
	h := sha256.New()
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
