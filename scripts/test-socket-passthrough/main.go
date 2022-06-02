package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/socket"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	listener, err := socket.Listen(":0")
	if err != nil {
		return err
	}
	// run `go build`
	childPath := filepath.Join(os.TempDir(), "test-socket-passthrough", "child")
	cmd := exe.Command(ctx, "go", "build", "-o", childPath, "scripts/test-socket-passthrough/child/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return err
	}
	// Start the web server
	cmd = exe.Command(ctx, childPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	appFile, err := listener.File()
	if err != nil {
		return err
	}
	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "APP", appFile)
	if err := cmd.Start(); err != nil {
		return err
	}
	res, err := http.DefaultClient.Get("http://" + listener.Addr().String())
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: " + res.Status)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println("Got body: " + string(body))
	if err := cmd.Restart(context.Background()); err != nil {
		return err
	}
	fmt.Println("restarted")
	res, err = http.DefaultClient.Get("http://" + listener.Addr().String())
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: " + res.Status)
	}
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println("Got body: " + string(body))
	if err := cmd.Close(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
