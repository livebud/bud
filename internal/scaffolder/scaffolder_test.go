package scaffolder_test

import (
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/scaffolder"
)

func TestTemplate(t *testing.T) {
	is := is.New(t)
	scaffold, err := scaffolder.Load()
	is.NoErr(err)
	err = scaffold.Generate(
		scaffold.Template("hey.gotext", "hello {{.}}", "ok"),
		scaffold.JSON("hey.gotext", "ok"),
	)
	is.NoErr(err)
	err = scaffold.Command("npm", "install").Run()
	is.NoErr(err)
}

func TestExistingDir(t *testing.T) {
	t.SkipNow()
}

func TestExistingEmptyDir(t *testing.T) {
	t.SkipNow()
}

func TestSymlink(t *testing.T) {
	t.SkipNow()
}

func TestCommand(t *testing.T) {
	t.SkipNow()
}