package bud_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/package/gomod"

	"github.com/livebud/bud/internal/bud"
	"github.com/livebud/bud/internal/testdir"
	runtime_bud "github.com/livebud/bud/runtime/bud"
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
	dir := t.TempDir()
	fmt.Println("loading dir")
	td := testdir.New()
	is.NoErr(td.Write(dir))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	fmt.Println("loaded")
	fmt.Println("loading compiler")
	compiler, err := bud.Load(module)
	is.NoErr(err)
	fmt.Println("loaded compiler")
	ctx := context.Background()
	fmt.Println("compiling")
	project, err := compiler.Compile(ctx, &runtime_bud.Flag{})
	is.NoErr(err)
	is.NoErr(exists(filepath.Join(dir, "bud", "cli")))
	is.NoErr(notExists(filepath.Join(dir, "bud", "app")))
	_, err = project.Build(ctx)
	is.NoErr(err)
	is.NoErr(exists(filepath.Join(dir, "bud", "cli")))
	is.NoErr(exists(filepath.Join(dir, "bud", "app")))
}
