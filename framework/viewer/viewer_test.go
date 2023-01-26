package viewer_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testcli"
	"github.com/livebud/bud/internal/testdir"
)

func TestNoViews(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run", "--embed=false")
	is.NoErr(err)
	defer app.Close()
	// HTML response
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 404 Not Found
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff
	`))
	is.NoErr(app.Close())
}
