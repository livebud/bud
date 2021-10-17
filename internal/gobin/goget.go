package gobin

import (
	"context"
	"os"
	"os/exec"
)

// GoGet calls `go get -u importPath@version`
func GoGet(ctx context.Context, dir, importPath string) error {
	cmd := exec.CommandContext(ctx, "go", "get", "-u", importPath)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
