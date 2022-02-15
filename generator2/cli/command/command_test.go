package command_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/package/testcmd"
)

func TestEmpty(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	err := os.RemoveAll(dir)
	is.NoErr(err)
	td := testdir.New()
	err = td.Write(dir)
	is.NoErr(err)
	stdout, err := testcmd.Expand(dir)
	is.NoErr(err)
	fmt.Println(stdout)
}
