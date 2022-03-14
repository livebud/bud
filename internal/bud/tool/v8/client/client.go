package client

import (
	"context"
	"encoding/gob"
	"errors"
	"io"
	"os"

	"gitlab.com/mnm/bud/internal/bud"
	v8 "gitlab.com/mnm/bud/pkg/js/v8"
)

type Command struct {
	Bud *bud.Command
}

func (c *Command) Run(ctx context.Context) error {
	vm := v8.New()
	dec := gob.NewDecoder(os.Stdin)
	enc := gob.NewEncoder(os.Stdout)
	for {
		var expr string
		if err := dec.Decode(&expr); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		result, err := vm.Eval("<stdin>", string(expr))
		if err != nil {
			return err
		}
		if err := enc.Encode(result); err != nil {
			return err
		}
	}
}
