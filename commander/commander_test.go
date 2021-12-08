package commander_test

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
	"gitlab.com/mnm/bud/commander"
)

func TestHelp(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cmd := commander.New("cli").Writer(actual)
	err := cmd.Parse([]string{"-h"})
	is.NoErr(err)
	isEqual(t, actual.String(), `
  {bold}Usage:{reset}
    cli

`)
}

func TestInvalid(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cmd := commander.New("cli").Writer(actual)
	err := cmd.Parse([]string{"blargle"})
	is.Equal(err.Error(), "unexpected blargle")
	isEqual(t, actual.String(), ``)
}
func TestSimple(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cli := commander.New("cli").Writer(actual)
	called := 0
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	err := cli.Parse([]string{})
	is.NoErr(err)
	is.Equal(1, called)
	isEqual(t, actual.String(), ``)
}
func TestFlagString(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flag string
	cli.Flag("flag", "cli flag").String(&flag)
	err := cli.Parse([]string{"--flag", "cool"})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(flag, "cool")
	isEqual(t, actual.String(), ``)
}
func TestFlagStringDefault(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flag string
	cli.Flag("flag", "cli flag").String(&flag).Default("default")
	err := cli.Parse([]string{})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(flag, "default")
	isEqual(t, actual.String(), ``)
}

func TestFlagStringRequired(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flag string
	cli.Flag("flag", "cli flag").String(&flag)
	err := cli.Parse([]string{})
	is.Equal(err.Error(), "missing --flag")
}
func TestFlagInt(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flag int
	cli.Flag("flag", "cli flag").Int(&flag)
	err := cli.Parse([]string{"--flag", "10"})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(flag, 10)
	isEqual(t, actual.String(), ``)
}
func TestFlagIntDefault(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flag int
	cli.Flag("flag", "cli flag").Int(&flag).Default(10)
	err := cli.Parse([]string{})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(flag, 10)
	isEqual(t, actual.String(), ``)
}

func TestFlagIntRequired(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flag int
	cli.Flag("flag", "cli flag").Int(&flag)
	err := cli.Parse([]string{})
	is.Equal(err.Error(), "missing --flag")
}
func TestFlagBool(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flag bool
	cli.Flag("flag", "cli flag").Bool(&flag)
	err := cli.Parse([]string{"--flag"})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(flag, true)
	isEqual(t, actual.String(), ``)
}
func TestFlagBoolDefault(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flag bool
	cli.Flag("flag", "cli flag").Bool(&flag).Default(true)
	err := cli.Parse([]string{})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(flag, true)
	isEqual(t, actual.String(), ``)
}

func TestFlagBoolRequired(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flag bool
	cli.Flag("flag", "cli flag").Bool(&flag)
	err := cli.Parse([]string{})
	is.Equal(err.Error(), "missing --flag")
}

func TestFlagStrings(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flags []string
	cli.Flag("flag", "cli flag").Strings(&flags)
	err := cli.Parse([]string{"--flag", "1", "--flag", "2"})
	is.NoErr(err)
	is.Equal(len(flags), 2)
	is.Equal(flags[0], "1")
	is.Equal(flags[1], "2")
}

func TestFlagStringsRequired(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flags []string
	cli.Flag("flag", "cli flag").Strings(&flags)
	err := cli.Parse([]string{})
	is.Equal(err.Error(), "missing --flag")
}

func TestFlagStringsDefault(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flags []string
	cli.Flag("flag", "cli flag").Strings(&flags).Default("a", "b")
	err := cli.Parse([]string{})
	is.NoErr(err)
	is.Equal(len(flags), 2)
	is.Equal(flags[0], "a")
	is.Equal(flags[1], "b")
}

func TestFlagStringMap(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flags map[string]string
	cli.Flag("flag", "cli flag").StringMap(&flags)
	err := cli.Parse([]string{"--flag", "a:1 + 1", "--flag", "b:2"})
	is.NoErr(err)
	is.Equal(len(flags), 2)
	is.Equal(flags["a"], "1 + 1")
	is.Equal(flags["b"], "2")
}

func TestFlagStringMapRequired(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flags map[string]string
	cli.Flag("flag", "cli flag").StringMap(&flags)
	err := cli.Parse([]string{})
	is.Equal(err.Error(), "missing --flag")
}

func TestFlagStringMapDefault(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var flags map[string]string
	cli.Flag("flag", "cli flag").StringMap(&flags).Default(map[string]string{
		"a": "1",
		"b": "2",
	})
	err := cli.Parse([]string{})
	is.NoErr(err)
	is.Equal(len(flags), 2)
	is.Equal(flags["a"], "1")
	is.Equal(flags["b"], "2")
}

func TestArgStringMap(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var args map[string]string
	cli.Arg("arg", "cli arg").StringMap(&args)
	// Can have only one arg
	err := cli.Parse([]string{"a:1 + 1"})
	is.NoErr(err)
	is.Equal(len(args), 1)
	is.Equal(args["a"], "1 + 1")
}

func TestArgStringMapRequired(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var args map[string]string
	cli.Arg("arg", "cli arg").StringMap(&args)
	err := cli.Parse([]string{})
	is.Equal(err.Error(), "missing arg")
}

func TestArgStringMapDefault(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var args map[string]string
	cli.Arg("arg", "cli arg").StringMap(&args).Default(map[string]string{
		"a": "1",
		"b": "2",
	})
	err := cli.Parse([]string{})
	is.NoErr(err)
	is.Equal(len(args), 2)
	is.Equal(args["a"], "1")
	is.Equal(args["b"], "2")
}

func TestSub(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cli := commander.New("bud").Writer(actual)
	var trace []string
	cli.Run(func(ctx context.Context) error {
		trace = append(trace, "bud")
		return nil
	})
	{
		sub := cli.Command("run", "run your application")
		sub.Run(func(ctx context.Context) error {
			trace = append(trace, "run")
			return nil
		})
	}
	{
		sub := cli.Command("build", "build your application")
		sub.Run(func(ctx context.Context) error {
			trace = append(trace, "build")
			return nil
		})
	}
	err := cli.Parse([]string{"build"})
	is.NoErr(err)
	is.Equal(len(trace), 1)
	is.Equal(trace[0], "build")
	isEqual(t, actual.String(), ``)
}

func TestSubHelp(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cli := commander.New("bud").Writer(actual)
	cli.Flag("log", "specify the logger").Bool(nil).Default(false)
	var trace []string
	cli.Run(func(ctx context.Context) error {
		trace = append(trace, "bud")
		return nil
	})
	{
		sub := cli.Command("run", "run your application")
		sub.Run(func(ctx context.Context) error {
			trace = append(trace, "run")
			return nil
		})
	}
	{
		sub := cli.Command("build", "build your application")
		sub.Run(func(ctx context.Context) error {
			trace = append(trace, "build")
			return nil
		})
	}
	err := cli.Parse([]string{"-h"})
	is.NoErr(err)
	isEqual(t, actual.String(), `
  {bold}Usage:{reset}
    bud {dim}[flags]{reset} {dim}[command]{reset}

  {bold}Flags:{reset}
    --log  {dim}specify the logger{reset}

  {bold}Commands:{reset}
    build  {dim}build your application{reset}
    run    {dim}run your application{reset}

`)
}
func TestSubHelpShort(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cli := commander.New("bud").Writer(actual)
	cli.Flag("log", "specify the logger").Short('L').Bool(nil).Default(false)
	cli.Flag("debug", "set the debugger").Bool(nil).Default(true)
	var trace []string
	cli.Run(func(ctx context.Context) error {
		trace = append(trace, "bud")
		return nil
	})
	{
		sub := cli.Command("run", "run your application")
		sub.Run(func(ctx context.Context) error {
			trace = append(trace, "run")
			return nil
		})
	}
	{
		sub := cli.Command("build", "build your application")
		sub.Run(func(ctx context.Context) error {
			trace = append(trace, "build")
			return nil
		})
	}
	err := cli.Parse([]string{"-h"})
	is.NoErr(err)
	isEqual(t, actual.String(), `
  {bold}Usage:{reset}
    bud {dim}[flags]{reset} {dim}[command]{reset}

  {bold}Flags:{reset}
    -L, --log  {dim}specify the logger{reset}
    --debug    {dim}set the debugger{reset}

  {bold}Commands:{reset}
    build  {dim}build your application{reset}
    run    {dim}run your application{reset}

`)
}

func TestArgString(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var arg string
	cli.Arg("arg", "cli arg").String(&arg)
	err := cli.Parse([]string{"cool"})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(arg, "cool")
	isEqual(t, actual.String(), ``)
}

func TestArgStringDefault(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var arg string
	cli.Arg("arg", "cli arg").String(&arg).Default("default")
	err := cli.Parse([]string{})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(arg, "default")
	isEqual(t, actual.String(), ``)
}

func TestArgStringRequired(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var arg string
	cli.Arg("arg", "cli arg").String(&arg)
	err := cli.Parse([]string{})
	is.Equal(err.Error(), "missing arg")
}

func TestSubArgString(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var arg string
	cli.Command("build", "build command")
	cli.Command("run", "run command")
	cli.Arg("arg", "cli arg").String(&arg)
	err := cli.Parse([]string{"deploy"})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(arg, "deploy")
	isEqual(t, actual.String(), ``)
}

// TestInterrupt tests interrupts canceling context. It spawns a copy of itself
// to run a subcommand. I learned this trick from Mitchell Hashimoto's excellent
// "Advanced Testing with Go" talk. We use stdout to synchronize between the
// process and subprocess.
func TestInterrupt(t *testing.T) {
	is := is.New(t)
	if value := os.Getenv("TEST_INTERRUPT"); value == "" {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cmd := exec.CommandContext(ctx, os.Args[0], append(os.Args[1:], "-test.v=true")...)
		cmd.Env = append(os.Environ(), "TEST_INTERRUPT=1")
		stdout, err := cmd.StdoutPipe()
		is.NoErr(err)
		cmd.Stderr = os.Stderr
		is.NoErr(cmd.Start())
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "ready" {
				break
			}
		}
		cmd.Process.Signal(os.Interrupt)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "cancelled" {
				break
			}
		}
		if err := cmd.Wait(); err != nil {
			is.True(errors.Is(err, context.Canceled))
		}
		return
	}
	cli := commander.New("cli")
	cli.Run(func(ctx context.Context) error {
		os.Stdout.Write([]byte("ready\n"))
		<-ctx.Done()
		os.Stdout.Write([]byte("cancelled\n"))
		return nil
	})
	if err := cli.Parse([]string{}); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		is.NoErr(err)
	}
}

// // TODO: more tests
// var tests = []struct {
// 	name string
// 	test func(t testing.TB)
// }{
// 	{
// 		name: "help",
// 		test: func(t testing.TB) {
// 			actual := new(bytes.Buffer)
// 			cmd := commander.New("cli", "same command").Writer(actual)
// 			err := cmd.Parse([]string{"-h"})
// 			assert.NoError(t, err)
// 			equal(t, actual.String(), `
// 				Usage:

// 					say [<flags>]

// 				Flags:

// 					-h, --help	Output usage information.
// 			`)
// 		},
// 	},
// 	{
// 		name: "invalid",
// 		test: func(t testing.TB) {
// 			actual := new(bytes.Buffer)
// 			cmd := commander.New("cli", "same command").Writer(actual)
// 			err := cmd.Parse([]string{"blargle"})
// 			assert.EqualError(t, err, "unexpected blargle")
// 			equal(t, actual.String(), ``)
// 		},
// 	},
// 	{
// 		name: "simple",
// 		test: func(t testing.TB) {
// 			actual := new(bytes.Buffer)
// 			cmd := commander.New("cli", "same command").Writer(actual)
// 			called := 0
// 			cmd.Run(func() error {
// 				called++
// 				return nil
// 			})
// 			err := cmd.Parse([]string{})
// 			assert.NoError(t, err)
// 			assert.Equal(t, 1, called)
// 			equal(t, actual.String(), ``)
// 		},
// 	},
// 	{
// 		name: "run error",
// 		test: func(t testing.TB) {
// 			actual := new(bytes.Buffer)
// 			cmd := commander.New("cli", "same command").Writer(actual)
// 			called := 0
// 			cmd.Run(func() error {
// 				called++
// 				return errors.New("oh noz")
// 			})
// 			err := cmd.Parse([]string{})
// 			assert.EqualError(t, err, "oh noz")
// 			assert.Equal(t, 1, called)
// 			equal(t, actual.String(), ``)
// 		},
// 	},
// 	{
// 		name: "help with example",
// 		test: func(t testing.TB) {
// 			actual := new(bytes.Buffer)
// 			cmd := commander.New("cli", "same command").Writer(actual)
// 			cmd.Example("say <something>", "say something")
// 			cmd.Example("say <something> [else]", "say something else")
// 			err := cmd.Parse([]string{"-h"})
// 			assert.NoError(t, err)
// 			equal(t, actual.String(), `
// 				Usage:

// 					say [<flags>]

// 				Flags:

// 					-h, --help	Output usage information.

// 				Examples:

// 					say something
// 					$ say <something>

// 					say something else
// 					$ say <something> [else]
// 			`)
// 		},
// 	},
// 	{
// 		name: "subcommand help with example",
// 		test: func(t testing.TB) {
// 			actual := new(bytes.Buffer)
// 			cmd := commander.New("cli", "same command").Writer(actual)
// 			en := cmd.Command("en", "say in english")
// 			en.Example("say en <something>", "say something")
// 			en.Example("say en <something> [else]", "say something else")
// 			err := cmd.Parse([]string{"en", "-h"})
// 			assert.NoError(t, err)
// 			equal(t, actual.String(), `
// 				say in english

// 				Usage:

// 					say en

// 				Flags:

// 					-h, --help	Output usage information.

// 				Examples:

// 					say something
// 					$ say en <something>

// 					say something else
// 					$ say en <something> [else]
// 			`)
// 		},
// 	},
// 	{
// 		name: "before function",
// 		test: func(t testing.TB) {
// 			actual := new(bytes.Buffer)
// 			cmd := commander.New("cli", "same command").Writer(actual)
// 			called := 0
// 			cmd.Before(func() error {
// 				called++
// 				return nil
// 			})
// 			en := cmd.Command("en", "say in english")
// 			en.Run(func() error {
// 				called++
// 				return nil
// 			})
// 			err := cmd.Parse([]string{"en"})
// 			assert.NoError(t, err)
// 			assert.Equal(t, 2, called)
// 			equal(t, actual.String(), ``)
// 		},
// 	},
// }

// func Test(t *testing.T) {
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			test.test(t)
// 		})
// 	}
// }

func isEqual(t testing.TB, actual, expected string) {
	t.Helper()
	equal(t, expected, replaceEscapeCodes(actual))
}

func replaceEscapeCodes(str string) string {
	// TODO: this is needlessly slow
	str = strings.ReplaceAll(str, "\033[0m", `{reset}`)
	str = strings.ReplaceAll(str, "\033[1m", `{bold}`)
	str = strings.ReplaceAll(str, "\033[37m", `{dim}`)
	str = strings.ReplaceAll(str, "\033[4m", `{underline}`)
	str = strings.ReplaceAll(str, "\033[36m", `{teal}`)
	str = strings.ReplaceAll(str, "\033[34m", `{blue}`)
	str = strings.ReplaceAll(str, "\033[33m", `{yellow}`)
	str = strings.ReplaceAll(str, "\033[31m", `{red}`)
	str = strings.ReplaceAll(str, "\033[32m", `{green}`)
	return str
}

// is checks if expect and actual are equal
func equal(t testing.TB, expect, actual string) {
	t.Helper()
	if expect == actual {
		return
	}
	var b bytes.Buffer
	b.WriteString("\n\x1b[4mExpect\x1b[0m:\n")
	b.WriteString(expect)
	b.WriteString("\n\n")
	b.WriteString("\x1b[4mActual\x1b[0m: \n")
	b.WriteString(actual)
	b.WriteString("\n\n")
	b.WriteString("\x1b[4mDifference\x1b[0m: \n")
	b.WriteString(diff.String(expect, actual))
	b.WriteString("\n")
	t.Fatal(b.String())
}
