package v8server_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/js/v8client"
	"github.com/livebud/bud/package/js/v8server"
	"golang.org/x/sync/errgroup"
)

func TestClient(t *testing.T) {
	is := is.New(t)
	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()
	client := v8client.New(r1, w2)
	server := v8server.New(r2, w1)
	eg := new(errgroup.Group)
	eg.Go(func() error { return server.Serve() })
	val, err := client.Eval("input.js", "2+2")
	is.NoErr(err)
	is.Equal(val, "4")
	w2.Close() // done writing
	is.NoErr(eg.Wait())
}

func TestMultiClient(t *testing.T) {
	is := is.New(t)

	cmd := exec.Command("go", "run", "test-server.go")
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	stdin, err := cmd.StdinPipe()
	is.NoErr(err)
	stdout, err := cmd.StdoutPipe()
	is.NoErr(err)
	is.NoErr(cmd.Start())
	client := v8client.New(stdout, stdin)
	// server := &v8server.Server{r2, w1}
	// eg := new(errgroup.Group)
	// eg.Go(func() error { return server.Serve() })
	err = client.Script("fib.js", `
		function fib(num) {
			if (num <= 1) return 1;
			return fib(num - 1) + fib(num - 2);
		}
	`)
	is.NoErr(err)
	// Run many times
	eg := new(errgroup.Group)
	for i := 0; i < 100; i++ {
		i := i
		// Create a bunch of goroutines alternating between 2+2 and 3+3
		eg.Go(func() error {
			isEven := i%2 == 0
			expr := "3+3"
			if isEven {
				expr = "2+2"
			}
			val, err := client.Eval("input.js", expr)
			if err != nil {
				return err
			}
			if isEven && val != "4" {
				return fmt.Errorf("%s != %s", expr, val)
			} else if !isEven && val != "6" {
				return fmt.Errorf("%s != %s", expr, val)
			}
			return nil
		})
	}
	is.NoErr(eg.Wait())
	is.NoErr(stdin.Close())
	is.NoErr(cmd.Wait())
}
