package create_test

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"golang.org/x/mod/modfile"
)

func fileFirstLine(filePath string) string {
	file, _ := os.Open(filePath)
	defer file.Close()
	scanner := bufio.NewReader(file)
	line, _ := scanner.ReadString('\n')
	return line
}

func TestCreateOutsideGoPath(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(dir)
	is.NoErr(td.NotExists(".gitignore"))
	result, err := cli.Run(ctx, "create", dir)
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.Equal(fileFirstLine(filepath.Join(dir, "go.mod")), "module change.me\n")
	is.NoErr(td.Exists(".gitignore"))
	is.NoErr(td.Exists("go.sum"))
	is.NoErr(td.Exists("package.json"))
	is.NoErr(td.Exists("package-lock.json"))
}

func TestCreateOutsideGoPathModulePath(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	cli := testcli.New(dir)
	is.NoErr(td.NotExists(".gitignore"))
	result, err := cli.Run(ctx, "create", "--module=github.com/my/app", dir)
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.Equal(fileFirstLine(filepath.Join(dir, "go.mod")), "module github.com/my/app\n")
	is.NoErr(td.Exists(".gitignore"))
	is.NoErr(td.Exists("go.sum"))
	is.NoErr(td.Exists("package.json"))
	is.NoErr(td.Exists("package-lock.json"))
}

func TestAutoQuote(t *testing.T) {
	is := is.New(t)
	actual := modfile.AutoQuote(`github.com/livebud/bud`)
	is.Equal(actual, `github.com/livebud/bud`)
	actual = modfile.AutoQuote(`github.com/livebud/bud with spaces`)
	is.Equal(actual, `"github.com/livebud/bud with spaces"`)
}
