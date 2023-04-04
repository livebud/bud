package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/livebud/bud/commands/internal/mainfile"
	"github.com/livebud/bud/internal/dag"

	"github.com/livebud/bud/commands"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/levelfilter"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/virtual"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args[1:]...); err != nil {
		console.Error(err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args ...string) error {
	flags := flag.NewFlagSet("commands", flag.ExitOnError)
	if err := flags.Parse(args); err != nil {
		return err
	}
	args = flags.Args()
	if len(args) < 1 {
		return fmt.Errorf("usage: run main.go <div>")
	}
	log := log.New(levelfilter.New(console.New(os.Stderr), log.InfoLevel))
	module, err := gomod.Find(args[0])
	if err != nil {
		return err
	}
	genfs := genfs.New(dag.Discard, module, log)
	parser := parser.New(genfs, module)
	genfs.FileGenerator("bud/command/command.go", commands.New(module, parser))
	injector := di.New(genfs, log, module, parser)
	// Load the mainfile
	genfs.FileGenerator("bud/main.go", mainfile.New(injector, module))
	// Sync the code
	if err := virtual.Sync(log, genfs, module, "bud"); err != nil {
		return err
	}
	// Build the mainfile
	cmd := exec.Command("go", "build", "-mod=mod", "-o=bud/main", "bud/main.go")
	cmd.Dir = module.Directory()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	// Run the mainfile
	stdout, err := execute(module.Directory(), args[1:]...)
	if err != nil {
		return err
	}
	fmt.Print(stdout)
	return nil
}

func execute(dir string, args ...string) (string, error) {
	cmd := exec.Command("./bud/main", args...)
	cmd.Dir = dir
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s: %w", stderr.String(), err)
	}
	return stdout.String(), nil
}
