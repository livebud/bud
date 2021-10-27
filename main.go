package main

import (
	"context"
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/log/console"
)

func main() {
	if err := run(); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	bud := flag.NewFlagSet("bud", flag.ContinueOnError)
	bud.SetOutput(ioutil.Discard)
	chdir := bud.String("chdir", ".", "change the working directory")
	if err := bud.Parse(os.Args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			return err
		}
		// TODO: handle bud -h when not in an app directory
	}
	abspath, err := filepath.Abs(*chdir)
	if err != nil {
		return err
	}
	mainPath := filepath.Join(abspath, "bud", "command", "main.go")
	// TODO: generate command/main.go
	// TODO: better context
	return gobin.Run(context.Background(), abspath, mainPath, bud.Args()...)
}
