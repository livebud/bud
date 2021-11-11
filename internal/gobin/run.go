package gobin

import (
	"context"
	"os"
	"os/exec"
)

// Run calls `go run -mod=mod main.go ...`
func Run(ctx context.Context, dir, mainpath string, args ...string) error {
	cmd := exec.CommandContext(ctx, "go", append([]string{"run", "-mod=mod", mainpath}, args...)...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	// stderr := new(bytes.Buffer)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		if isCleanExit(err) {
			return nil
		}
		return err
	}
	return nil
}

// Interpret exit code 3 as an error-free exit. This allows interrupts like
// SIGINT to exit cleanly.
func isCleanExit(err error) bool {
	if e, ok := err.(*exec.ExitError); ok {
		if e.ExitCode() == 3 {
			return true
		}
	}
	return false
}
