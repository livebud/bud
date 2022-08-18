package once

import "github.com/livebud/bud/internal/errs"

type Closer struct {
	closes []func() error
	once   Error
}

func (c *Closer) Add(closes ...func() error) {
	c.closes = append(c.closes, closes...)
}

func (c *Closer) Close(reasons ...error) error {
	return c.once.Do(func() error {
		err := errs.Join(reasons...)
		for i := len(c.closes) - 1; i >= 0; i-- {
			err = errs.Join(err, c.closes[i]())
		}
		return err
	})
}
