package cli_test

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testcli"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/versions"
	"golang.org/x/mod/modfile"
)

func fileFirstLine(filePath string) string {
	file, _ := os.Open(filePath)
	defer file.Close()
	scanner := bufio.NewReader(file)
	line, _ := scanner.ReadString('\n')
	return line
}

func TestCreateOutsideGoPath(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(dir)
	is.NoErr(td.NotExists(".gitignore"))
	result, err := cli.Run(ctx, "create", dir)
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.In(result.Stderr(), "Ready")
	is.NoErr(td.Exists(".gitignore"))
	is.NoErr(td.Exists("go.sum"))
	is.NoErr(td.Exists("package.json"))
	is.NoErr(td.Exists("package-lock.json"))
	is.Equal(fileFirstLine(filepath.Join(dir, "go.mod")), "module change.me\n")
}

func TestCreateOutsideGoPathModulePath(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	cli := testcli.New(dir)
	is.NoErr(td.NotExists(".gitignore"))
	result, err := cli.Run(ctx, "create", "--module=github.com/my/app", dir)
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.In(result.Stderr(), "Ready")
	is.NoErr(td.Exists(".gitignore"))
	is.NoErr(td.Exists("go.sum"))
	is.NoErr(td.Exists("package.json"))
	is.NoErr(td.Exists("package-lock.json"))
	is.NoErr(td.Exists("public/favicon.ico"))
	is.Equal(fileFirstLine(filepath.Join(dir, "go.mod")), "module github.com/my/app\n")
}

func TestModfileAutoQuote(t *testing.T) {
	is := is.New(t)
	actual := modfile.AutoQuote(`github.com/livebud/bud`)
	is.Equal(actual, `github.com/livebud/bud`)
	actual = modfile.AutoQuote(`github.com/livebud/bud with spaces`)
	is.Equal(actual, `"github.com/livebud/bud with spaces"`)
}

func TestCreateSeesWelcome(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(dir)
	is.NoErr(td.NotExists(".gitignore"))
	result, err := cli.Run(ctx, "create", dir)
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.In(result.Stderr(), "Ready")

	// Start the app
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	// Test the index page
	res, err := app.Get("/")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	is.In(res.Body().String(), "Hey Bud")
	is.In(res.Body().String(), "Hey Bud") // should work multiple times
	// Test the client-side JS
	res, err = app.Get("/bud/view/_index.svelte.js")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	is.In(res.Body().String(), "Hey Bud")
	is.Equal(app.Stdout(), "")
	is.NoErr(td.Exists("bud/app"))
	is.NoErr(app.Close())
}

func TestCreateRemoveBudGraceful(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(dir)
	is.NoErr(td.NotExists(".gitignore"))
	result, err := cli.Run(ctx, "create", dir)
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.In(result.Stderr(), "Ready")

	// Start the app
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	// Test the index page
	res, err := app.Get("/")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	is.In(res.Body().String(), "Hey Bud")
	is.NoErr(td.Exists("bud/app"))
	is.NoErr(td.Exists("bud/afs"))
	is.NoErr(td.Exists("bud/bud.db"))

	// Remove the bud directory
	is.NoErr(os.RemoveAll(filepath.Join(dir, "bud")))
	is.NoErr(td.NotExists("bud/app"))
	is.NoErr(td.NotExists("bud/afs"))
	is.NoErr(td.NotExists("bud/bud.db"))

	// Wait for the app to shutdown gracefully
	is.NoErr(app.Wait())

	// New startup should be graceful
	app, err = cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Test the index page
	res, err = app.Get("/")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	is.In(res.Body().String(), "Hey Bud")
	is.NoErr(td.Exists("bud/app"))
	is.NoErr(td.Exists("bud/afs"))
	is.NoErr(td.Exists("bud/bud.db"))
	is.NoErr(app.Close())
}

func TestCreateRemovePublicGraceful(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(dir)
	is.NoErr(td.NotExists(".gitignore"))
	result, err := cli.Run(ctx, "create", dir)
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.In(result.Stderr(), "Ready")

	// Start the app
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)

	// Test the favicon
	res, err := app.Get("/favicon.ico")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	is.NoErr(td.Exists("public/favicon.ico"))

	// Remove the public directory
	is.NoErr(os.RemoveAll(filepath.Join(dir, "public")))
	is.NoErr(td.NotExists("bud/public"))

	// Wait for the app to recover
	readyCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	is.NoErr(app.Ready(readyCtx))
	cancel()

	res, err = app.Get("/favicon.ico")
	is.NoErr(err)
	is.Equal(res.Status(), 404)

	is.NoErr(app.Close())
}

func TestReleaseVersionOk(t *testing.T) {
	if versions.Bud == "latest" {
		t.Skip("Skipping release version test for latest")
	}
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(dir)
	is.NoErr(td.NotExists(".gitignore"))
	result, err := cli.Run(ctx, "create", "--dev=false", dir)
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.In(result.Stderr(), "Ready")
	is.NoErr(td.Exists(".gitignore"))
	is.NoErr(td.Exists("go.sum"))
	is.NoErr(td.Exists("package.json"))
	is.NoErr(td.Exists("package-lock.json"))
	is.Equal(fileFirstLine(filepath.Join(dir, "go.mod")), "module change.me\n")
	gomod, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	is.NoErr(err)
	is.In(string(gomod), "github.com/livebud/bud v"+versions.Bud)
}
