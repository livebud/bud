package gobin

import (
	"context"
	"os"
	"os/exec"

	"gitlab.com/mnm/bud/package/gomod"
)

type Builder interface {
	Build(ctx context.Context, module *gomod.Module, mainPath string, outPath string, flags ...string) error
}

// Build calls `go build -mod=mod -o main [flags...] main.go`
func Build(ctx context.Context, module *gomod.Module, mainPath string, outPath string, flags ...string) error {
	// Compile the args
	args := append([]string{
		"build",
		"-mod=mod",
		"-o=" + outPath,
	}, flags...)
	args = append(args, mainPath)
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = module.Directory()
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
