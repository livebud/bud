package gobin

import (
	"context"
	"os"
	"os/exec"
)

// Test calls `go test -mod=mod <testpath>.go ...`
func Test(ctx context.Context, dir, testpath string, args ...string) error {
	cmd := exec.CommandContext(ctx, "go", append([]string{"test", "-mod=mod", testpath}, args...)...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
