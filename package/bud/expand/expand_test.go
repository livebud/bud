package expand_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/bud"
	"github.com/livebud/bud/package/bud/expand"
	"github.com/livebud/bud/package/gomod"
	"github.com/matryer/is"
)

func exists(paths ...string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}
	return nil
}

func notExists(paths ...string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); nil == err {
			return fmt.Errorf("%s exists but shouldn't", path)
		}
	}
	return nil
}

func TestEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	err = (&expand.Command{
		Module: module,
		Flag:   &bud.Flag{},
	}).Expand(ctx)
	is.NoErr(err)
	is.NoErr(exists(module.Directory("bud/cli")))
	is.NoErr(notExists(module.Directory("bud/.app")))
}
