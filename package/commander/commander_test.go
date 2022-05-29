package commander_test

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/commander"
	"github.com/matthewmueller/diff"
)

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

func TestHelp(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cmd := commander.New("cli").Writer(actual)
	ctx := context.Background()
	err := cmd.Parse(ctx, []string{"-h"})
	is.NoErr(err)
	isEqual(t, actual.String(), `
  {bold}Usage:{reset}
    cli

`)
}

func TestHelpArgs(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cmd := commander.New("cp").Writer(actual)
	cmd.Arg("src").String(nil)
	cmd.Arg("dst").String(nil).Default(".")
	ctx := context.Background()
	err := cmd.Parse(ctx, []string{"-h"})
	is.NoErr(err)
	isEqual(t, actual.String(), `
  {bold}Usage:{reset}
    cp {dim}<src>{reset} {dim}<dst>{reset}

`)
}

func TestInvalid(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cmd := commander.New("cli").Writer(actual)
	ctx := context.Background()
	err := cmd.Parse(ctx, []string{"blargle"})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"--flag", "cool"})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"--flag", "10"})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"--flag"})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"--flag", "1", "--flag", "2"})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"--flag", "a:1 + 1", "--flag", "b:2"})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	cli.Arg("arg").StringMap(&args)
	// Can have only one arg
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"a:1 + 1"})
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
	cli.Arg("arg").StringMap(&args)
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	cli.Arg("arg").StringMap(&args).Default(map[string]string{
		"a": "1",
		"b": "2",
	})
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"build"})
	is.NoErr(err)
	is.Equal(len(trace), 1)
	is.Equal(trace[0], "build")
	isEqual(t, actual.String(), ``)
}

func TestSubHelp(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cli := commander.New("bud").Writer(actual)
	cli.Flag("log", "specify the logger").Bool(nil)
	cli.Command("run", "run your application")
	cli.Command("build", "build your application")
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"-h"})
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

func TestEmptyUsage(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cli := commander.New("bud").Writer(actual)
	cli.Flag("log", "").Bool(nil)
	cli.Command("run", "")
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"-h"})
	is.NoErr(err)
	isEqual(t, actual.String(), `
  {bold}Usage:{reset}
    bud {dim}[flags]{reset} {dim}[command]{reset}

  {bold}Flags:{reset}
    --log

  {bold}Commands:{reset}
    run

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
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"-h"})
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
	cli.Arg("arg").String(&arg)
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"cool"})
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
	cli.Arg("arg").String(&arg).Default("default")
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	cli.Arg("arg").String(&arg)
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
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
	cli.Arg("arg").String(&arg)
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"deploy"})
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
		// Ignore -test.count otherwise this will continue recursively
		var args []string
		for _, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, "-test.count=") {
				continue
			}
			args = append(args, arg)
		}
		cmd := exec.CommandContext(ctx, os.Args[0], append(args, "-test.v=true", "-test.run=^TestInterrupt$")...)
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
	ctx := context.Background()
	if err := cli.Parse(ctx, []string{}); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		is.NoErr(err)
	}
}

// TODO: example support

func TestArgsStrings(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var args []string
	cli.Command("build", "build command")
	cli.Command("run", "run command")
	cli.Args("custom").Strings(&args)
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"new", "view"})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(len(args), 2)
	is.Equal(args[0], "new")
	is.Equal(args[1], "view")
	isEqual(t, actual.String(), ``)
}

func TestUsageError(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return commander.Usage()
	})
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
	is.NoErr(err)
	isEqual(t, actual.String(), `
  {bold}Usage:{reset}
    cli

`)
}

func TestIdempotent(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cli := commander.New("cli").Writer(actual)
	var f1 string
	cmd := cli.Command("run", "run command")
	cmd.Flag("f1", "cli flag").Short('f').String(&f1)
	var f2 string
	cmd.Flag("f2", "cli flag").String(&f2)
	var f3 string
	cmd.Flag("f3", "cli flag").String(&f3)
	ctx := context.Background()
	args := []string{"run", "--f1=a", "--f2=b", "--f3", "c"}
	err := cli.Parse(ctx, args)
	is.NoErr(err)
	is.Equal(f1, "a")
	is.Equal(f2, "b")
	is.Equal(f3, "c")
	f1 = ""
	f2 = ""
	f3 = ""
	err = cli.Parse(ctx, args)
	is.NoErr(err)
	is.Equal(f1, "a")
	is.Equal(f2, "b")
	is.Equal(f3, "c")
}

func TestManualHelp(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cli := commander.New("cli").Writer(actual)
	var help bool
	var dir string
	cli.Flag("help", "help menu").Short('h').Bool(&help).Default(false)
	cli.Flag("chdir", "change directory").Short('C').String(&dir)
	called := 0
	cli.Run(func(ctx context.Context) error {
		is.Equal(help, true)
		is.Equal(dir, "somewhere")
		called++
		return nil
	})
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"--help", "--chdir", "somewhere"})
	is.NoErr(err)
	is.Equal(actual.String(), "")
	is.Equal(called, 1)
}

func TestManualHelpUsage(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cli := commander.New("cli").Writer(actual)
	var help bool
	var dir string
	cli.Flag("help", "help menu").Short('h').Bool(&help).Default(false)
	cli.Flag("chdir", "change directory").Short('C').String(&dir)
	called := 0
	cli.Run(func(ctx context.Context) error {
		is.Equal(help, true)
		is.Equal(dir, "somewhere")
		called++
		return commander.Usage()
	})
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"--help", "--chdir", "somewhere"})
	is.NoErr(err)
	is.Equal(called, 1)
	isEqual(t, actual.String(), `
  {bold}Usage:{reset}
    cli {dim}[flags]{reset}

  {bold}Flags:{reset}
    -C, --chdir  {dim}change directory{reset}
    -h, --help   {dim}help menu{reset}

`)
}

func TestAfterRun(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cli := commander.New("cli").Writer(actual)
	called := 0
	var ctx context.Context
	cli.Run(func(c context.Context) error {
		called++
		ctx = c
		return nil
	})
	err := cli.Parse(context.Background(), []string{})
	is.NoErr(err)
	is.Equal(called, 1)
	select {
	case <-ctx.Done():
		is.Fail() // Context shouldn't have been cancelled
	default:
	}
}

func TestArgsClearSlice(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	args := []string{"a", "b"}
	cli.Args("custom").Strings(&args)
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"c", "d"})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(len(args), 2)
	is.Equal(args[0], "c")
	is.Equal(args[1], "d")
	isEqual(t, actual.String(), ``)
}

func TestArgClearMap(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	args := map[string]string{"a": "a"}
	cli.Arg("custom").StringMap(&args)
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"b:b"})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(len(args), 1)
	is.Equal(args["b"], "b")
	isEqual(t, actual.String(), ``)
}

func TestFlagClearSlice(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	args := []string{"a", "b"}
	cli.Flag("f", "flag").Strings(&args)
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"-f", "c", "-f", "d"})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(len(args), 2)
	is.Equal(args[0], "c")
	is.Equal(args[1], "d")
	isEqual(t, actual.String(), ``)
}

func TestFlagClearMap(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	args := map[string]string{"a": "a"}
	cli.Flag("f", "flag").StringMap(&args)
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"-f", "b:b"})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(len(args), 1)
	is.Equal(args["b"], "b")
	isEqual(t, actual.String(), ``)
}

func TestFlagCustom(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	hot := ""
	cli.Flag("hot", "hot server").Custom(func(v string) error {
		hot = v
		return nil
	})
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"--hot=:35729"})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(hot, ":35729")
}

func TestFlagCustomError(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Flag("hot", "hot server").Custom(func(v string) error {
		return fmt.Errorf("unable to parse")
	})
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"--hot=:35729"})
	is.True(err != nil)
	is.Equal(err.Error(), `invalid value ":35729" for flag -hot: unable to parse`)
}

func TestFlagCustomMissing(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	cli.Flag("hot", "hot server").Custom(func(v string) error {
		return nil
	})
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
	is.True(err != nil)
	is.Equal(err.Error(), `missing --hot`)
}

func TestFlagCustomMissingDefault(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	hot := ""
	cli.Flag("hot", "hot server").Custom(func(v string) error {
		hot = v
		return nil
	}).Default(":35729")
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	ctx := context.Background()
	err := cli.Parse(ctx, []string{})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(hot, ":35729")
}

func TestFlagCustomDefault(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("cli").Writer(actual)
	hot := ""
	cli.Flag("hot", "hot server").Custom(func(v string) error {
		hot = v
		return nil
	}).Default(":35729")
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	ctx := context.Background()
	err := cli.Parse(ctx, []string{"--hot=false"})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(hot, "false")
}
