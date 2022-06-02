package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/js/v8client"
	"github.com/livebud/bud/package/socket"
)

func run(ctx context.Context) error {
	listener, err := socket.Listen(":4444")
	if err != nil {
		return err
	}
	defer listener.Close()
	fileListener, err := listener.File()
	if err != nil {
		return err
	}
	v8client, err := v8client.Load(ctx)
	if err != nil {
		return err
	}
	if err := v8client.Script("script.js", "const __svelte__ = 3+3"); err != nil {
		return err
	}
	fmt.Println("calling child")
	cmd := exe.Command(ctx, "go", "run", "scripts/test-v8-pipe/grandchild/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	extrafile.Forward(&cmd.ExtraFiles, &cmd.Env, "V8")
	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "APP", fileListener)
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
