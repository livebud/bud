package cli_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/internal/tester"
)

// TODO: show help on empty
func TestEmpty(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := "_tmp"
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli, err := tester.Compile(ctx, dir)
	is.NoErr(err)
	stdout, stderr, err := cli.Command("-h")
	is.NoErr(err)
	is.Equal(stdout, "Usage")
	is.Equal(stderr, "")
}
