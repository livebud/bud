package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/livebud/bud/internal/shell"

	"github.com/livebud/bud/internal/extrafile"
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
	cmd := exec.CommandContext(ctx, "go", "build", "-o", childPath, "scripts/test-socket-passthrough/child/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return err
	}
	// Start the web server
	cmd = exec.CommandContext(ctx, childPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	appFile, err := listener.File()
	if err != nil {
		return err
	}
	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "APP", appFile)
	process, err := shell.Start(ctx, cmd)
	if err != nil {
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
	process, err = process.Restart(context.Background())
	if err != nil {
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
	if err := process.Close(); err != nil {
		return err
	}
	if err := process.Wait(); err != nil {
		return err
	}
	return nil
}
