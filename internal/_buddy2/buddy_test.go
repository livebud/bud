package buddy_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/pkg/gomod"
)

func TestCompile(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	module, err := gomod.Find(".")
	is.NoErr(err)
	bud := driver.New(module)
	cli, err := bud.Compile(ctx, driver.WithEmbed(true))
	is.NoErr(err)
	cmd := cli.Command(ctx, "-h")
	process, err := cli.Run(ctx, driver.WithPort(":3000"))
	is.NoErr(err)
	process.Close()
}
