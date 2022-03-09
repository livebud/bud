package clic_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"gitlab.com/mnm/bud/framework/clid"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/testdir"
)

// TODO: show help on empty
func TestEmpty(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	err := os.RemoveAll(dir)
	is.NoErr(err)
	td := testdir.New()
	err = td.Write(dir)
	is.NoErr(err)
	cli, err := clid.Compile(ctx, dir)
	is.NoErr(err)
	stdout, stderr, err := cli.Run("-h")
	is.NoErr(err)
	is.Equal(stdout, "")
	is.Equal(stderr, "")
}

func TestHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := "_tmp"
	err := os.RemoveAll(dir)
	is.NoErr(err)
	td := testdir.New()
	err = td.Write(dir)
	is.NoErr(err)
	cli, err := clid.Compile(ctx, dir)
	is.NoErr(err)
	stdout, stderr, err := cli.Run("-h")
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
	cli, err := clid.Compile(ctx, dir)
	is.NoErr(err)
	stdout, stderr, err := cli.Run("build")
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
