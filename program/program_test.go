package program_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/livebud/bud/cli"
	"github.com/livebud/bud/di"
	"github.com/livebud/bud/program"
	"github.com/matryer/is"
)

func TestParse(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	in := di.New()
	di.Loader[cli.Parser](in, func(in di.Injector) (cli.Parser, error) {
		return cli.New("cli", "some cli"), nil
	})
	err := program.Parse(ctx, in, "-h")
	is.NoErr(err)
}

func TestParseSubcommand(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	in := di.New()
	di.Loader[cli.Parser](in, func(in di.Injector) (cli.Parser, error) {
		cli := cli.New("cli", "some cli")
		cmd := cli.Command("cmd", "some command")
		cmd.Run(func(ctx context.Context) error {
			return nil
		})
		return cli, nil
	})
	err := program.Parse(ctx, in, "cmd")
	is.NoErr(err)
}

func TestParseError(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	in := di.New()
	expected := errors.New("some error")
	di.Loader[cli.Parser](in, func(in di.Injector) (cli.Parser, error) {
		cli := cli.New("cli", "some cli")
		cmd := cli.Command("cmd", "some command")
		cmd.Run(func(ctx context.Context) error {
			return expected
		})
		return cli, nil
	})
	err := program.Parse(ctx, in, "cmd")
	is.True(errors.Is(err, expected))
}

func TestLoad(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	in := di.New()
	di.Loader[cli.Parser](in, func(in di.Injector) (cli.Parser, error) {
		cli := cli.New("cli", "some cli")
		cmd := cli.Command("cmd", "some command")
		cmd.Run(func(ctx context.Context) error {
			return nil
		})
		return cli, nil
	})
	cli, err := program.Load(in)
	is.NoErr(err)
	err = cli.Parse(ctx, "cmd")
	is.NoErr(err)
}

func ExampleRun() {
	loadCLI := func(in di.Injector) (cli.Parser, error) {
		cli := cli.New("cli", "some cli")
		cmd := cli.Command("say", "say hello")
		cmd.Run(func(ctx context.Context) error {
			fmt.Println("hello")
			return nil
		})
		return cli, nil
	}
	in := di.New()
	di.Loader[cli.Parser](in, loadCLI)
	ctx := context.Background()
	program.Parse(ctx, in, "say")
	// Output:
	// hello
}
