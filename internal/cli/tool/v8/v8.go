package v8

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	v8 "github.com/livebud/bud/package/js/v8"
)

type Command struct {
	Stdin  io.Reader
	Stdout io.Writer
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
	fmt.Fprintln(c.Stdout, result)
	return nil
}

func (c *Command) getScript() (string, error) {
	code, err := ioutil.ReadAll(c.Stdin)
	if err != nil {
		return "", err
	}
	script := string(code)
	if script == "" {
		return "", errors.New("missing script to evaluate")
	}
	return script, nil
}
