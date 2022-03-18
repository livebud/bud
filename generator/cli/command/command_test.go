package command_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/budtest"
	"gitlab.com/mnm/bud/internal/testdir"
)

func TestEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	bud, err := budtest.Find(dir)
	is.NoErr(err)
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	stdout, stderr, err := project.Execute(ctx, "-h")
	is.NoErr(err)
	is.NoErr(stdout.Contains("Usage:"))
	is.NoErr(stdout.Contains("-C, --chdir"))
	is.NoErr(stdout.Contains("build"))
	is.NoErr(stdout.Contains("run"))
	is.NoErr(stderr.Expect(""))
}
