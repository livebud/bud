package main

import (
	"context"
	"fmt"
	"os"

	"gitlab.com/mnm/bud/commander"
	"gitlab.com/mnm/bud/log/console"
)

func main() {
	if err := run(); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	cli := commander.New("bud")
	cmd := new(Bud)
	cli.Run(cmd.Run)
	// cli.StringVar(&cmd.workDir, "chdir", ".", "change the working directory")
	return cli.Parse(os.Args[1:])
}

type Bud struct {
	workDir string
}

func (b *Bud) Run(ctx context.Context) error {
	fmt.Println("running in", b.workDir)
	return nil
}
