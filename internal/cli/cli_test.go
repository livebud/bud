package cli_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/testdir"
	"github.com/matryer/is"
)

func exists(t testing.TB, dir string, path string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(dir, path)); err != nil {
		t.Fatalf("%q should exist but doesn't: %s", path, err)
	}
}

func notExists(t testing.TB, dir string, path string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(dir, path)); nil == err {
		t.Fatalf("%q exists but shouldn't: %s", path, err)
	}
}

func TestBuildEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli := cli.New(dir)
	cli.Env["TMPDIR"] = t.TempDir()
	err = cli.Parse(ctx, "build")
	is.NoErr(err)
	exists(t, dir, "bud/cli")
	exists(t, dir, "bud/app")
}

func TestHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli := cli.New(dir)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cli.Stdout = stdout
	cli.Stderr = stderr
	cli.Env["NO_COLOR"] = "1"
	cli.Env["TMPDIR"] = t.TempDir()
	err = cli.Parse(ctx, "--help")
	is.NoErr(err)
	is.Equal(stderr.String(), "")
	is.True(strings.Contains(stdout.String(), "cli"))
	is.True(strings.Contains(stdout.String(), "build command"))
	is.True(strings.Contains(stdout.String(), "run command"))
	is.True(strings.Contains(stdout.String(), "new scaffold"))
}

func TestBuildChdir(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli := cli.New(".")
	cli.Env["TMPDIR"] = t.TempDir()
	err = cli.Parse(ctx, "--chdir", dir, "build")
	is.NoErr(err)
}
