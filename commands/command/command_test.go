package command_test

import (
	"context"
	"flag"
	"fmt"
	"io"
	"testing"

	"github.com/livebud/bud/internal/is"
)

func TestFlagSet(t *testing.T) {
	is := is.New(t)
	fset := flag.NewFlagSet("bud", flag.ContinueOnError)
	cmd := &Command{
		Global: &Global{},
	}
	fset.Var(&stringValue{&String{&cmd.App, nil}, false}, "app", "")
	// fset.Var(&stringValue{&String{cmd.Remote, nil}, false}, "remote", "")
	fset.Var(&stringValue{&String{&cmd.As, nil}, false}, "as", "")
	// fset.StringVar(&cmd.App, "app", "", "")
	// // fset.StringVar(cmd.Remote, "remote", "", "")
	// fset.StringVar(&cmd.As, "as", "", "")
	fset.SetOutput(io.Discard)
	args := []string{}
	all := []string{
		"--app=myapp",
		// "--remote", "heroku",
		"addons:attach",
		"heroku-postgresql:hobby-dev",
		"--app=myappz",
		"-as", "mydb",
	}
	for {
		err := fset.Parse(all)
		is.NoErr(err)
		if fset.NArg() == 0 {
			break
		}
		args = append(args, fset.Arg(0))
		all = fset.Args()[1:]
	}

	// err = fset.Parse(fset.Args()[1:])
	// is.NoErr(err)
	// err = fset.Parse(fset.Args()[1:])
	// is.NoErr(err)
	// fmt.Println(fset.Args())
	// fmt.Println(fset.NFlag())
	// fmt.Println(fset.Args())
	fmt.Println("got args", args)
	cmd.Run(context.Background())
}

type Global struct {
	App    string
	Remote *string
}

type Command struct {
	*Global
	As string
}

func (c *Command) Run(ctx context.Context) error {
	fmt.Println("app", c.App)
	fmt.Println("remote", c.Remote)
	fmt.Println("as", c.As)
	return nil
}

type String struct {
	target *string
	defval *string // default value
}

func (v *String) Default(value string) {
	v.defval = &value
}

func (v *String) Optional() {
	v.defval = new(string)
}

type stringValue struct {
	inner *String
	set   bool
}

// func (v *stringValue) verify(displayName string) error {
// 	if v.set {
// 		return nil
// 	} else if v.inner.defval != nil {
// 		*v.inner.target = *v.inner.defval
// 		return nil
// 	}
// 	return fmt.Errorf("missing %s", displayName)
// }

func (v *stringValue) Set(val string) error {
	if v.inner.target == nil {
		v.inner.target = new(string)
	}
	*v.inner.target = val
	v.set = true
	return nil
}

func (v *stringValue) String() string {
	if v.inner == nil {
		return ""
	} else if v.set {
		return *v.inner.target
	} else if v.inner.defval != nil {
		return *v.inner.defval
	}
	return ""
}
