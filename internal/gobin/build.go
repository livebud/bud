package gobin

import (
	"context"
	"os"
	"os/exec"
)

// Build calls `go build -mod=mod -o main main.go`
func Build(ctx context.Context, dir, mainpath string, outpath string) error {
	cmd := exec.CommandContext(ctx, "go", "build", "-mod=mod", "-o", outpath, mainpath)
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
