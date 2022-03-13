package expand_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"gitlab.com/mnm/bud/pkg/bud"
	"gitlab.com/mnm/bud/pkg/buddy"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/buddy/expand"
	"gitlab.com/mnm/bud/internal/testdir"
)

func TestHelp(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	td := testdir.New()
	dir := "_tmp"
	is.NoErr(os.RemoveAll(dir))
	err := td.Write(dir)
	is.NoErr(err)
	kit, err := buddy.Load(dir)
	is.NoErr(err)
	ctx := context.Background()
	bud := bud.New(kit)
	err = bud.Expand(ctx, &expand.Input{})
	is.NoErr(err)
	cmd := exec.Command("./bud/cli", "-h")
	cmd.Dir = dir
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	is.NoErr(err)
	fmt.Println(stdout.String())
}
