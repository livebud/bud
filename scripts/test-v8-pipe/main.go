package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/livebud/bud/internal/extrafile"

	"github.com/livebud/bud/package/exe"

	"github.com/livebud/bud/package/js/v8server"
)

func run(ctx context.Context) error {
	r1, w2, err := os.Pipe()
	if err != nil {
		return err
	}
	// defer w2.Close()
	r2, w1, err := os.Pipe()
	if err != nil {
		return err
	}
	// defer w1.Close()
	v8Server := v8server.New(r1, w1)
	go func() {
		fmt.Println("err serving", v8Server.Serve())
	}()

	cmd := exe.Command(ctx, "go", "run", "scripts/test-v8-pipe/child/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "V8", r2, w2)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
