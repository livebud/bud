package cli_test

import (
	"context"
	"strings"
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
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli, err := tester.Compile(ctx, dir)
	is.NoErr(err)
	stdout, stderr, err := cli.Command("-h")
	is.NoErr(err)
	is.Equal(stdout, "")
	is.Equal(stderr, "")
}

func TestHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli, err := tester.Compile(ctx, dir)
	is.NoErr(err)
	stdout, stderr, err := cli.Command("-h")
	is.NoErr(err)
	is.Equal(stderr, "")
	is.True(strings.Contains(stdout, "Usage")) // should contain Usage
	is.True(strings.Contains(stdout, "build")) // should contain build
	is.True(strings.Contains(stdout, "run"))   // should contain run
}

func TestBuild(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli, err := tester.Compile(ctx, dir)
	is.NoErr(err)
	stdout, stderr, err := cli.Command("build")
	is.NoErr(err)
	is.Equal(stderr, "")
	is.Equal(stdout, "")
	// // Run the app
	// app := testapp.New(dir)
	// stdout, stderr, err = app.Run("-h")
	// is.NoErr(err)
	// is.Equal(stderr, "")
	// is.True(strings.Contains(stdout, "Usage:")) // should contain Usage
	// is.True(strings.Contains(stdout, "app"))    // should contain app
}
