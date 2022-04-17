package imhash_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/package/gomod"

	"github.com/livebud/bud/internal/imhash"

	"github.com/matryer/is"
	"github.com/livebud/bud/internal/budtest"
)

func TestAppHash(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	bud.Files["controller/controller.go"] = `
		package controller

		type Controller struct {
		}

		func (c *Controller) Index() string {
			return "Hello Users!"
		}
	`
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	is.NoErr(project.Exists("bud/.cli/main.go"))

	module, err := gomod.Find(dir)
	is.NoErr(err)

	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/main.go"))

	hash1, err := imhash.Hash(module, "bud/.app")
	is.NoErr(err)
	is.Equal(len(hash1), 11)

	// Update
	project.Files["controller/controller.go"] = `
		package controller

		type Controller struct {
		}

		func (c *Controller) Index() string {
			return "Hello Users!!"
		}
	`
	is.NoErr(project.Rewrite())
	app, err = project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/main.go"))

	hash2, err := imhash.Hash(module, "bud/.app")
	is.NoErr(err)
	is.Equal(len(hash2), 11)
	is.True(hash1 != hash2)
}
