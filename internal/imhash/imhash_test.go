package imhash_test

import (
	"context"
	"testing"

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
	cli := testcli.New(dir)
	result, err := cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists("bud/internal/app/main.go"))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	hash1, err := imhash.Hash(module, "bud/internal/app")
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
	result, err = cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists("bud/internal/app/main.go"))
	hash2, err := imhash.Hash(module, "bud/internal/app")
	is.NoErr(err)
	is.Equal(len(hash2), 11)
	is.True(hash1 != hash2)
}

func TestAppHashEmbeds(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		import "embed"
		//go:embed *.sql css/*
		//go:embed html/index.html
		var Migrations embed.FS
		type Controller struct {}
		func (c *Controller) Index() string {
			return "Hello Users!!"
		}
	`
	td.Files["controller/1.up.sql"] = `CREATE UNIVERSE;`
	td.Files["controller/1.down.sql"] = `DROP UNIVERSE;`
	td.Files["controller/html/index.html"] = `<blink>hi</blink>`
	td.Files["controller/css/a.css"] = `body { font-family: Comic Sans; }`
	td.Files["controller/css/b.css"] = `body { color: red; }`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	result, err := cli.Run(ctx, "build")
	is.NoErr(td.Exists("bud/internal/app/main.go"))
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	module, err := gomod.Find(dir)
	is.NoErr(err)
	hash1, err := imhash.Hash(module, "bud/internal/app")
	is.NoErr(err)
	is.Equal(len(hash1), 11)

	// Add more migrations
	td.Files["controller/2.up.sql"] = `ALTER UNIVERSE;`
	td.Files["controller/2.down.sql"] = `UNALTER UNIVERSE;`
	is.NoErr(td.Write(ctx))
	result, err = cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.NoErr(td.Exists("bud/internal/app/main.go"))
	hash2, err := imhash.Hash(module, "bud/internal/app")
	is.NoErr(err)
	is.Equal(len(hash2), 11)
	is.True(hash1 != hash2)

	// Update a CSS file
	td.Files["controller/css/b.css"] = `body { color: green; }`
	is.NoErr(td.Write(ctx))
	result, err = cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.NoErr(td.Exists("bud/internal/app/main.go"))
	hash3, err := imhash.Hash(module, "bud/internal/app")
	is.NoErr(err)
	is.Equal(len(hash3), 11)
	is.True(hash2 != hash3)

	// Update the HTML file
	td.Files["controller/html/index.html"] = `<strong>hi</strong>`
	is.NoErr(td.Write(ctx))
	result, err = cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.NoErr(td.Exists("bud/internal/app/main.go"))
	hash4, err := imhash.Hash(module, "bud/internal/app")
	is.NoErr(err)
	is.Equal(len(hash3), 11)
	is.True(hash3 != hash4)

	// Update an unembedded HTML file
	td.Files["controller/html/other.html"] = `doesnt matter`
	is.NoErr(td.Write(ctx))
	result, err = cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.NoErr(td.Exists("bud/internal/app/main.go"))
	hash5, err := imhash.Hash(module, "bud/internal/app")
	is.NoErr(err)
	is.Equal(len(hash3), 11)
	is.Equal(hash4, hash5)
}
