package once

import "github.com/livebud/bud/internal/errs"

type Closer struct {
	Closes []func() error
	once   Error
}

func (c *Closer) Close(reasons ...error) error {
	return c.once.Do(func() error {
		err := errs.Join(reasons...)
		for i := len(c.Closes) - 1; i >= 0; i-- {
			err = errs.Join(err, c.Closes[i]())
		}
		return err
	})
}
