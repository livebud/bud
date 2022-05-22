package build_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/livebud/bud/package/bud"
	"github.com/livebud/bud/package/exe"

	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/bud/build"
	"github.com/matryer/is"
)

func TestWelcome(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	// Build the binary
	builder := &build.Command{
		Dir:    dir,
		Stdout: stdout,
		Stderr: stderr,
		Env: bud.Env{
			"TMPDIR": t.TempDir(),
		},
	}
	err = builder.Build(ctx)
	is.NoErr(err)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
	// Start the built binary
	cmd := exe.Command(ctx, "bud/app")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err = cmd.Start()
	is.NoErr(err)

	// defer func() { is.NoErr(process.Close()) }()
	// res, err := client.Get("http://host/")
	// is.NoErr(err)
	// defer res.Body.Close()
	// is.Equal(res.StatusCode, 200)
	// body, err := io.ReadAll(res.Body)
	// is.NoErr(err)
	// is.True(strings.Contains(string(body), "Hey Bud"))
}
