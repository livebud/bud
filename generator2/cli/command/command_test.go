package command_test

import (
	"os"
	"strings"
	"testing"

	"gitlab.com/mnm/bud/package/testapp"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/package/testcli"
)

// TODO: show help on empty
func TestEmpty(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	dir := t.TempDir()
	err := os.RemoveAll(dir)
	is.NoErr(err)
	td := testdir.New()
	err = td.Write(dir)
	is.NoErr(err)
	cli, err := testcli.Load(dir)
	is.NoErr(err)
	stdout, stderr, err := cli.Run("-h")
	is.NoErr(err)
	is.Equal(stdout, "")
	is.Equal(stderr, "")
}

func TestHelp(t *testing.T) {
	is := is.New(t)
	dir := "_tmp"
	err := os.RemoveAll(dir)
	is.NoErr(err)
	td := testdir.New()
	err = td.Write(dir)
	is.NoErr(err)
	cli, err := testcli.Load(dir)
	is.NoErr(err)
	stdout, stderr, err := cli.Run("-h")
	is.NoErr(err)
	is.Equal(stderr, "")
	is.True(strings.Contains(stdout, "Usage")) // should contain Usage
	is.True(strings.Contains(stdout, "build")) // should contain build
}

func TestBuild(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli, err := testcli.Load(dir)
	is.NoErr(err)
	stdout, stderr, err := cli.Run("build")
	is.NoErr(err)
	is.Equal(stderr, "")
	is.Equal(stdout, "")
	// Run the app
	app := testapp.New(dir)
	stdout, stderr, err = app.Run("-h")
	is.NoErr(err)
	is.True(strings.Contains(stdout, "Usage:")) // should contain Usage
	is.True(strings.Contains(stdout, "app"))    // should contain app
}
