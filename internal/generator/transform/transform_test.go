package transform_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
)

func TestNoProject(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	cli := testcli.New(cli.New(dir))
	stdout, stderr, err := cli.Run(ctx)
	is.NoErr(err)
	is.In(stdout.String(), "bud")
	is.Equal(stderr.String(), "")
	is.NoErr(td.NotExists("bud/.cli/transform/transform.go"))
}

func TestEmptyProject(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	stdout, stderr, err := cli.Run(ctx)
	is.NoErr(err)
	is.In(stdout.String(), "cli")
	is.Equal(stderr.String(), "")
	is.NoErr(td.Exists("bud/.cli/transform/transform.go"))
}

func TestMarkdownPlugin(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Modules["github.com/livebud/bud-test-plugin"] = "v0.0.6"
	td.Files["view/index.md"] = `# hello`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	is.In(res.Body().String(), `<h1>hello</h1>`)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
}
