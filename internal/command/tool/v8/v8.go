package v8

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/livebud/bud/internal/command"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/mattn/go-isatty"
)

type Command struct {
	Bud *command.Bud
}

func (c *Command) Run(ctx context.Context) error {
	script, err := c.getScript()
	if err != nil {
		return err
	}
	vm, err := v8.Load()
	if err != nil {
		return err
	}
	result, err := vm.Eval("script.js", script)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, result)
	return nil
}

func (c *Command) getScript() (string, error) {
	code, err := ioutil.ReadAll(stdin())
	if err != nil {
		return "", err
	}
	script := string(code)
	if script == "" {
		return "", errors.New("missing script to evaluate")
	}
	return script, nil
}

// input from stdin or empty object by default.
func stdin() io.Reader {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		return strings.NewReader("")
	}
	return os.Stdin
}
