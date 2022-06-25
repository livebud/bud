package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/livebud/bud/internal/current"
	"github.com/livebud/bud/internal/sig"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	ctx = sig.Trap(ctx, os.Interrupt)
	dirname, err := current.Directory()
	if err != nil {
		return err
	}
	// Build child
	cmd := exec.Command("go", "build", "-o", "child/main", "child/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dirname
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return err
	}
	// Run child
	cmd = exec.Command("child/main")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Dir = dirname
	if err := cmd.Start(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		fmt.Println("parent: interrupted!")
		if err := cmd.Wait(); err != nil {
			return err
		}
		fmt.Println("parent: exiting")
	}
	return nil
}
