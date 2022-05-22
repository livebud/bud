package custom_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/bud"
	"github.com/livebud/bud/package/bud/custom"
	"github.com/matryer/is"
)

func TestHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd := &custom.Command{
		Dir:    dir,
		Stdout: stdout,
		Stderr: stderr,
		Env: bud.Env{
			"NO_COLOR": "1",
			"TMPDIR":   t.TempDir(),
		},
		Args: []string{"--help"},
	}
	err = cmd.Run(ctx)
	is.NoErr(err)
	is.Equal(stderr.String(), "")
	is.True(strings.Contains(stdout.String(), "cli"))
	is.True(strings.Contains(stdout.String(), "build command"))
	is.True(strings.Contains(stdout.String(), "run command"))
	is.True(strings.Contains(stdout.String(), "new scaffold"))
}
