package controller_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
)

// // TODO: this test didn't discover the bud that "bud new" was no longer exposed
// // top-level, because project.Execute uses the inner CLI
func TestNewController(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	stdout, stderr, err := cli.Run(ctx, "new", "controller", "/", "index", "show")
	is.NoErr(err)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
	is.NoErr(td.Exists(
		"controller/controller.go",
		"view/index.svelte",
		"view/show.svelte",
	))
}
