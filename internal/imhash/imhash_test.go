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

func TestAppHashEmbed(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		import "embed"
		//go:embed *.sql
		var Migrations embed.FS
		type Controller struct {}
		func (c *Controller) Index() string {
			return "Hello Users!!"
		}
	`
	td.Files["controller/1.up.sql"] = `CREATE UNIVERSE;`
	td.Files["controller/1.down.sql"] = `DROP UNIVERSE;`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	stdout, stderr, err := cli.Run(ctx, "build")
	is.NoErr(td.Exists(
		"bud/.cli/main.go",
		"bud/.app/main.go",
	))
	is.NoErr(err)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
	module, err := gomod.Find(dir)
	is.NoErr(err)
	hash1, err := imhash.Hash(module, "bud")
	is.NoErr(err)
	is.Equal(len(hash1), 11)
	// Update
	td.Files["controller/2.up.sql"] = `ALTER UNIVERSE;`
	td.Files["controller/2.down.sql"] = `UNALTER UNIVERSE;`

	is.NoErr(td.Write(ctx))
	stdout, _, err = cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
	is.NoErr(td.Exists(
		"bud/.cli/main.go",
		"bud/.app/main.go",
	))
	hash2, err := imhash.Hash(module, "bud")
	is.NoErr(err)
	is.Equal(len(hash2), 11)
	is.True(hash1 != hash2)
}
