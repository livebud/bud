package cli_test

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/socket"
	"github.com/matryer/is"
)

func exists(t testing.TB, dir string, path string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(dir, path)); err != nil {
		t.Fatalf("%q should exist but doesn't: %s", path, err)
	}
}

func notExists(t testing.TB, dir string, path string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(dir, path)); nil == err {
		t.Fatalf("%q exists but shouldn't: %s", path, err)
	}
}

func listenUnix(socketPath string) (net.Listener, *http.Client, error) {
	transport, err := socket.Transport(socketPath)
	if err != nil {
		return nil, nil, err
	}
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	listener, err := socket.Listen(socketPath)
	if err != nil {
		return nil, nil, err
	}
	return listener, client, nil
}

type file interface {
	File() (*os.File, error)
}

func setupCLI(t testing.TB, cli *cli.CLI) (stdout, stderr *bytes.Buffer) {
	stdout = new(bytes.Buffer)
	stderr = new(bytes.Buffer)
	// Write to stdio as well so debugging doesn't become too confusing
	cli.Stdout = io.MultiWriter(stdout, os.Stdout)
	cli.Stderr = io.MultiWriter(stderr, os.Stderr)
	cli.Env["TMPDIR"] = t.TempDir()
	cli.Env["NO_COLOR"] = "1"
	return stdout, stderr
}

func injectListener(t testing.TB, cli *cli.CLI) (*http.Client, func()) {
	t.Helper()
	// Start listening on a unix domain socket
	socketPath := filepath.Join(t.TempDir(), "unix.sock")
	listener, client, err := listenUnix(socketPath)
	if err != nil {
		t.Fatalf("unable to listen on socket path %q: %s", socketPath, err)
	}
	// Pull the files and environment from the listener
	files, env, err := socket.Files(listener)
	if err != nil {
		t.Fatalf("unable to derive *os.File from net.Listener: %s", err)
	}
	// Inject into CLI
	cli.ExtraFiles = append(cli.ExtraFiles, files...)
	cli.Env[env.Key()] = env.Value()
	// Return the client and a way to shutdown the listener
	return client, func() {
		if err := listener.Close(); err != nil {
			t.Fatalf("unexpected error while closing listener: %s", err)
		}
	}
}

func startCLI(t testing.TB, cli *cli.CLI, args ...string) func() {
	t.Helper()
	process, err := cli.Start(context.Background(), args...)
	if err != nil {
		t.Fatalf("unable to start cli: %s", err)
	}
	return func() {
		if err := process.Close(); err != nil {
			t.Fatalf("error while closing cli: %s", err)
		}
	}
}

func TestBuildEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli := cli.New(dir)
	stdout, stderr := setupCLI(t, cli)
	err = cli.Run(ctx, "build")
	is.NoErr(err)
	exists(t, dir, "bud/cli")
	exists(t, dir, "bud/app")
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
}

func TestHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli := cli.New(dir)
	stdout, stderr := setupCLI(t, cli)
	err = cli.Run(ctx, "--help")
	is.NoErr(err)
	exists(t, dir, "bud/cli")
	notExists(t, dir, "bud/app")
	is.Equal(stderr.String(), "")
	is.True(strings.Contains(stdout.String(), "cli"))
	is.True(strings.Contains(stdout.String(), "build command"))
	is.True(strings.Contains(stdout.String(), "run command"))
	is.True(strings.Contains(stdout.String(), "new scaffold"))
}

func TestChdir(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli := cli.New(".")
	stdout, stderr := setupCLI(t, cli)
	err = cli.Run(ctx, "--chdir", dir, "build")
	is.NoErr(err)
	exists(t, dir, "bud/cli")
	exists(t, dir, "bud/app")
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
}

func TestChdirHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli := cli.New(".")
	stdout, stderr := setupCLI(t, cli)
	err = cli.Run(ctx, "--chdir", dir, "--help")
	is.NoErr(err)
	exists(t, dir, "bud/cli")
	notExists(t, dir, "bud/app")
	is.Equal(stderr.String(), "")
	is.True(strings.Contains(stdout.String(), "cli"))
	is.True(strings.Contains(stdout.String(), "build command"))
	is.True(strings.Contains(stdout.String(), "run command"))
	is.True(strings.Contains(stdout.String(), "new scaffold"))
}

func TestOutsideModuleUsage(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	cli := cli.New(dir)
	stdout, stderr := setupCLI(t, cli)
	err := cli.Run(ctx)
	is.NoErr(err)
	notExists(t, dir, "bud/cli")
	notExists(t, dir, "bud/app")
	is.Equal(stderr.String(), "")
	is.True(strings.Contains(stdout.String(), "bud"))
	is.True(strings.Contains(stdout.String(), "-C, --chdir"))
	is.True(strings.Contains(stdout.String(), "-h, --help"))
	is.True(strings.Contains(stdout.String(), "create"))
	is.True(strings.Contains(stdout.String(), "tool"))
	is.True(strings.Contains(stdout.String(), "version"))
}

func TestOutsideModuleHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	cli := cli.New(dir)
	stdout, stderr := setupCLI(t, cli)
	err := cli.Run(ctx, "--help")
	is.NoErr(err)
	notExists(t, dir, "bud/cli")
	notExists(t, dir, "bud/app")
	is.Equal(stderr.String(), "")
	is.True(strings.Contains(stdout.String(), "bud"))
	is.True(strings.Contains(stdout.String(), "-C, --chdir"))
	is.True(strings.Contains(stdout.String(), "-h, --help"))
	is.True(strings.Contains(stdout.String(), "create"))
	is.True(strings.Contains(stdout.String(), "tool"))
	is.True(strings.Contains(stdout.String(), "version"))
}

func TestToolV8(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	cli := cli.New(dir)
	stdout, stderr := setupCLI(t, cli)
	cli.Stdin = bytes.NewBufferString("2+2")
	err := cli.Run(ctx, "tool", "v8")
	is.NoErr(err)
	notExists(t, dir, "bud/cli")
	notExists(t, dir, "bud/app")
	is.Equal(stderr.String(), "")
	is.True(strings.Contains(stdout.String(), "4"))
}

func TestRunWelcome(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli := cli.New(dir)
	cli.Stdout = os.Stdout
	cli.Stderr = os.Stderr
	stdout, stderr := setupCLI(t, cli)
	client, close := injectListener(t, cli)
	defer close()
	cleanup := startCLI(t, cli, "run")
	defer cleanup()
	res, err := client.Get("http://host/")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(res.StatusCode, 200)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.True(strings.Contains(string(body), "Hey Bud"))
	is.Equal(stdout.String(), "")
	is.True(strings.Contains(stderr.String(), "info: Listening on "))
}

func TestRunController(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New()
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string { return "from index" }
	`
	err := td.Write(dir)
	is.NoErr(err)
	cli := cli.New(dir)
	cli.Stdout = os.Stdout
	cli.Stderr = os.Stderr
	stdout, stderr := setupCLI(t, cli)
	client, close := injectListener(t, cli)
	defer close()
	cleanup := startCLI(t, cli, "run")
	defer cleanup()
	res, err := client.Get("http://host/")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(res.StatusCode, 200)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.True(strings.Contains(string(body), "from index"))
	is.Equal(stdout.String(), "")
	is.True(strings.Contains(stderr.String(), "info: Listening on "))
}
