package cli_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/cli/clitest"
	"github.com/livebud/bud/internal/testdir"
	"github.com/matryer/is"
)

func TestBuildEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli := clitest.New(t, cli.New(dir))
	result := cli.Run(ctx, "build").NoErr()
	cli.File("bud/cli").Exists()
	cli.File("bud/app").Exists()
	result.Stdout.Equals("")
	result.Stderr.Equals("")
}

func TestBuildTwice(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli := clitest.New(t, cli.New(dir))
	result := cli.Run(ctx, "build").NoErr()
	cli.File("bud/cli").Exists()
	cli.File("bud/app").Exists()
	result.Stdout.Equals("")
	result.Stderr.Equals("")
	result = cli.Run(ctx, "build").NoErr()
	cli.File("bud/cli").Exists()
	cli.File("bud/app").Exists()
	result.Stdout.Equals("")
	result.Stderr.Equals("")
}

func TestHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli := clitest.New(t, cli.New(dir))
	result := cli.Run(ctx, "--help").NoErr()
	cli.File("bud/cli").Exists()
	cli.File("bud/app").NotExists()
	result.Stderr.Equals("")
	result.Stdout.Contains("cli")
	result.Stdout.Contains("build command")
	result.Stdout.Contains("run command")
	result.Stdout.Contains("new scaffold")
}

func TestChdir(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli := clitest.New(t, cli.New("."))
	result := cli.Run(ctx, "--chdir", dir, "build").NoErr()
	cli.File("bud/cli").Exists()
	cli.File("bud/app").Exists()
	result.Stderr.Equals("")
	result.Stdout.Equals("")
}

func TestChdirHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli := clitest.New(t, cli.New("."))
	res := cli.Run(ctx, "--chdir", dir, "--help").NoErr()
	cli.File("bud/cli").Exists()
	cli.File("bud/app").NotExists()
	res.Stderr.Equals("")
	res.Stdout.Contains("cli")
	res.Stdout.Contains("build command")
	res.Stdout.Contains("run command")
	res.Stdout.Contains("new scaffold")
}

// func TestOutsideModuleUsage(t *testing.T) {
// 	is := is.New(t)
// 	ctx := context.Background()
// 	dir := t.TempDir()
// 	cli := cli.New(dir)
// 	stdout, stderr := setupCLI(t, cli)
// 	err := cli.Run(ctx)
// 	is.NoErr(err)
// 	notExists(t, dir, "bud/cli")
// 	notExists(t, dir, "bud/app")
// 	is.Equal(stderr.String(), "")
// 	is.True(strings.Contains(stdout.String(), "bud"))
// 	is.True(strings.Contains(stdout.String(), "-C, --chdir"))
// 	is.True(strings.Contains(stdout.String(), "-h, --help"))
// 	is.True(strings.Contains(stdout.String(), "create"))
// 	is.True(strings.Contains(stdout.String(), "tool"))
// 	is.True(strings.Contains(stdout.String(), "version"))
// }

// func TestOutsideModuleHelp(t *testing.T) {
// 	is := is.New(t)
// 	ctx := context.Background()
// 	dir := t.TempDir()
// 	cli := cli.New(dir)
// 	stdout, stderr := setupCLI(t, cli)
// 	err := cli.Run(ctx, "--help")
// 	is.NoErr(err)
// 	notExists(t, dir, "bud/cli")
// 	notExists(t, dir, "bud/app")
// 	is.Equal(stderr.String(), "")
// 	is.True(strings.Contains(stdout.String(), "bud"))
// 	is.True(strings.Contains(stdout.String(), "-C, --chdir"))
// 	is.True(strings.Contains(stdout.String(), "-h, --help"))
// 	is.True(strings.Contains(stdout.String(), "create"))
// 	is.True(strings.Contains(stdout.String(), "tool"))
// 	is.True(strings.Contains(stdout.String(), "version"))
// }

// func TestToolV8(t *testing.T) {
// 	is := is.New(t)
// 	ctx := context.Background()
// 	dir := t.TempDir()
// 	cli := cli.New(dir)
// 	stdout, stderr := setupCLI(t, cli)
// 	cli.Stdin = bytes.NewBufferString("2+2")
// 	err := cli.Run(ctx, "tool", "v8")
// 	is.NoErr(err)
// 	notExists(t, dir, "bud/cli")
// 	notExists(t, dir, "bud/app")
// 	is.Equal(stderr.String(), "")
// 	is.True(strings.Contains(stdout.String(), "4"))
// }

func TestRunWelcome(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli := clitest.New(t, cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	res.Status(200)
	_ = res
	// cli := cli.New(dir)
	// cli.Stdout = os.Stdout
	// cli.Stderr = os.Stderr
	// stdout, stderr := setupCLI(t, cli)
	// client, close := injectListener(t, cli)
	// defer close()
	// cleanup := startCLI(t, cli, "run")
	// defer cleanup()
	// res, err := client.Get("http://host/")
	// is.NoErr(err)
	// defer res.Body.Close()
	// is.Equal(res.StatusCode, 200)
	// body, err := io.ReadAll(res.Body)
	// is.NoErr(err)
	// is.True(strings.Contains(string(body), "Hey Bud"))
	// is.Equal(stdout.String(), "")
	// is.True(strings.Contains(stderr.String(), "info: Listening on "))
}

// func TestRunController(t *testing.T) {
// 	is := is.New(t)
// 	dir := t.TempDir()
// 	td := testdir.New()
// 	td.Files["controller/controller.go"] = `
// 		package controller
// 		type Controller struct {}
// 		func (c *Controller) Index() string { return "from index" }
// 	`
// 	err := td.Write(dir)
// 	is.NoErr(err)
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
