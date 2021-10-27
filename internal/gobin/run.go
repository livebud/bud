package gobin

import (
	"context"
	"os"
	"os/exec"
)

// Run calls `go run -u importPath@version`
func Run(ctx context.Context, dir, mainpath string, args ...string) error {
	cmd := exec.CommandContext(ctx, "go", append([]string{"run", mainpath}, args...)...)
	// GOPRIVATE=* to ensure we're not dealing with any caches
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	// stderr := new(bytes.Buffer)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
