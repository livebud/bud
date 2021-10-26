package commander_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
	"gitlab.com/mnm/bud/commander"
)

func TestHelp(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cmd := commander.New("say").Writer(actual)
	err := cmd.Parse([]string{"-h"})
	is.NoErr(err)
	isEqual(t, actual.String(), `
{bold}Usage:{reset}
  say
`)
}

func TestInvalid(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cmd := commander.New("say").Writer(actual)
	err := cmd.Parse([]string{"blargle"})
	is.Equal(err.Error(), "unexpected blargle")
	isEqual(t, actual.String(), ``)
}
func TestSimple(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	cli := commander.New("say").Writer(actual)
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
	cli := commander.New("say").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var message string
	cli.Flag("message", "say the message").String(&message)
	err := cli.Parse([]string{"--message", "cool"})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(message, "cool")
	isEqual(t, actual.String(), ``)
}
func TestFlagStringDefault(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("say").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var message string
	cli.Flag("message", "say the message").String(&message).Default("default")
	err := cli.Parse([]string{})
	is.NoErr(err)
	is.Equal(1, called)
	is.Equal(message, "default")
	isEqual(t, actual.String(), ``)
}

func TestFlagStringRequired(t *testing.T) {
	is := is.New(t)
	actual := new(bytes.Buffer)
	called := 0
	cli := commander.New("say").Writer(actual)
	cli.Run(func(ctx context.Context) error {
		called++
		return nil
	})
	var message string
	cli.Flag("message", "say the message").String(&message)
	err := cli.Parse([]string{})
	is.Equal(err.Error(), "missing message")
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
  bud

{bold}Commands:{reset}
  build  build your application
  run  run your application
`)
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
// 			cmd := commander.New("say", "same command").Writer(actual)
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
// 			cmd := commander.New("say", "same command").Writer(actual)
// 			err := cmd.Parse([]string{"blargle"})
// 			assert.EqualError(t, err, "unexpected blargle")
// 			equal(t, actual.String(), ``)
// 		},
// 	},
// 	{
// 		name: "simple",
// 		test: func(t testing.TB) {
// 			actual := new(bytes.Buffer)
// 			cmd := commander.New("say", "same command").Writer(actual)
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
// 			cmd := commander.New("say", "same command").Writer(actual)
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
// 			cmd := commander.New("say", "same command").Writer(actual)
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
// 			cmd := commander.New("say", "same command").Writer(actual)
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
// 			cmd := commander.New("say", "same command").Writer(actual)
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
