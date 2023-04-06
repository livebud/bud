package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testcli"
	"github.com/livebud/bud/package/testdir"
)

func TestToolV8(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	cli := testcli.New(td.Directory())
	cli.Stdin = bytes.NewBufferString("2+2")
	result, err := cli.Run(ctx, "v8")
	is.NoErr(err)
	is.Equal(result.Stderr(), "")
	is.Equal(strings.TrimSpace(result.Stdout()), "4")
	is.NoErr(td.NotExists(
		"bud/cmd/app",
		"bud/app",
	))
}
