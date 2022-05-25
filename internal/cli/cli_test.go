package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
)

func TestBuildEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(cli.New(dir))
	stdout, stderr, err := cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
	is.NoErr(td.Exists(
		"bud/cli",
		"bud/app",
	))
}

func TestBuildTwice(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	stdout, stderr, err := cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
	is.NoErr(td.Exists(
		"bud/cli",
		"bud/app",
	))
	stdout, stderr, err = cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
	is.NoErr(td.Exists(
		"bud/cli",
		"bud/app",
	))
}

func TestHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	stdout, stderr, err := cli.Run(ctx, "--help")
	is.NoErr(err)
	is.Equal(stderr.String(), "")
	is.In(stdout.String(), "cli")
	is.In(stdout.String(), "build command")
	is.In(stdout.String(), "run command")
	is.In(stdout.String(), "new scaffold")
	is.NoErr(td.Exists(
		"bud/cli",
	))
	is.NoErr(td.NotExists(
		"bud/app",
	))
}

func TestChdir(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New("."))
	stdout, stderr, err := cli.Run(ctx, "--chdir", dir, "build")
	is.NoErr(err)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
	is.NoErr(td.Exists(
		"bud/cli",
		"bud/app",
	))
}

func TestChdirHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New("."))
	stdout, stderr, err := cli.Run(ctx, "--chdir", dir, "--help")
	is.NoErr(err)
	is.Equal(stderr.String(), "")
	is.In(stdout.String(), "cli")
	is.In(stdout.String(), "build command")
	is.In(stdout.String(), "run command")
	is.In(stdout.String(), "new scaffold")
	is.NoErr(td.Exists(
		"bud/cli",
	))
	is.NoErr(td.NotExists(
		"bud/app",
	))
}

func TestOutsideModule(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	cli := testcli.New(cli.New(dir))
	stdout, stderr, err := cli.Run(ctx)
	is.NoErr(err)
	is.Equal(stderr.String(), "")
	is.In(stdout.String(), "bud")
	is.In(stdout.String(), "-C, --chdir")
	is.In(stdout.String(), "-h, --help")
	is.In(stdout.String(), "create")
	is.In(stdout.String(), "tool")
	is.In(stdout.String(), "version")
	is.NoErr(td.NotExists(
		"bud/cli",
		"bud/app",
	))
}

func TestOutsideModuleHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	cli := testcli.New(cli.New(dir))
	stdout, stderr, err := cli.Run(ctx, "--help")
	is.NoErr(err)
	is.Equal(stderr.String(), "")
	is.In(stdout.String(), "bud")
	is.In(stdout.String(), "-C, --chdir")
	is.In(stdout.String(), "-h, --help")
	is.In(stdout.String(), "create")
	is.In(stdout.String(), "tool")
	is.In(stdout.String(), "version")
	is.NoErr(td.NotExists(
		"bud/cli",
		"bud/app",
	))
}

func TestToolV8(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	cli := testcli.New(cli.New(dir))
	cli.Stdin(bytes.NewBufferString("2+2"))
	stdout, stderr, err := cli.Run(ctx, "tool", "v8")
	is.NoErr(err)
	is.Equal(stderr.String(), "")
	is.Equal(strings.TrimSpace(stdout.String()), "4")
	is.NoErr(td.NotExists(
		"bud/cli",
		"bud/app",
	))
}

func TestRunWelcome(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	is.In(res.Body().String(), "Hey Bud")
	is.In(res.Body().String(), "Hey Bud") // should work multiple times
	is.Equal(stdout.String(), "")
	is.In(stderr.String(), "Listening on")
}

// func TestRunController(t *testing.T) {
// 	is := is.New(t)
// 	dir := t.TempDir()
// 	td := testdir.New(dir)
// 	td.Files["controller/controller.go"] = `
// 		package controller
// 		type Controller struct {}
// 		func (c *Controller) Index() string { return "from index" }
// 	`
// 	is.NoErr(td.Write(ctx))
// 	cli := cli.New(dir)
// 	cli.Stdout = os.Stdout
// 	cli.Stderr = os.Stderr
// 	stdout, stderr := setupCLI(t, cli)
// 	client, close := injectListener(t, cli)
// 	defer close()
// 	cleanup := startCLI(t, cli, "run")
// 	defer cleanup()
// 	res, err := client.Get("http://host/")
// 	is.NoErr(err)
// 	defer res.Body.Close()
// 	is.Equal(res.StatusCode, 200)
// 	body, err := io.ReadAll(res.Body)
// 	is.NoErr(err)
// 	is.True(strings.Contains(string(body), "from index"))
// 	is.Equal(stdout.String(), "")
// 	is.True(strings.Contains(stderr.String(), "info: Listening on "))
// }
