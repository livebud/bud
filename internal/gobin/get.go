package gobin

import (
	"context"
	"os"
	"os/exec"
)

// Get calls `go get -u importPath@version`
func Get(ctx context.Context, dir string, importPaths ...string) error {
	cmd := exec.CommandContext(ctx, "go", append([]string{"get", "-u"}, importPaths...)...)
	// GOPRIVATE=* to ensure we're not dealing with any caches
	cmd.Env = append(os.Environ(), "GOPRIVATE=*")
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
