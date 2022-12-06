package once

import (
	"github.com/livebud/bud/internal/errs"
)

type Closer struct {
	closes []func() error
	once   Error
}

func (c *Closer) Add(fn func() error) {
	c.closes = append(c.closes, fn)
}

func (c *Closer) Close() (err error) {
	return c.once.Do(func() error {
		for i := len(c.closes) - 1; i >= 0; i-- {
			err = errs.Join(err, c.closes[i]())
		}
		return err
	})
}
