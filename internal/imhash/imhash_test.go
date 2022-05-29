package imhash_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/imhash"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/gomod"

	"github.com/livebud/bud/internal/is"
)

func TestAppHash(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string {
			return "Hello Users!"
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	stdout, stderr, err := cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
	is.NoErr(td.Exists(
		"bud/.cli/main.go",
		"bud/.app/main.go",
	))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	hash1, err := imhash.Hash(module, "bud/.app")
	is.NoErr(err)
	is.Equal(len(hash1), 11)
	// Update
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string {
			return "Hello Users!!"
		}
	`
	is.NoErr(td.Write(ctx))
	stdout, stderr, err = cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
	is.NoErr(td.Exists(
		"bud/.cli/main.go",
		"bud/.app/main.go",
	))
	hash2, err := imhash.Hash(module, "bud/.app")
	is.NoErr(err)
	is.Equal(len(hash2), 11)
	is.True(hash1 != hash2)
}
