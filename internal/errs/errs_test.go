package errs_test

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/errs"
	"github.com/livebud/bud/internal/is"
)

func TestNil(t *testing.T) {
	is := is.New(t)
	is.NoErr(errs.Join(nil))
}

func TestOne(t *testing.T) {
	is := is.New(t)
	err := errs.Join(errors.New("one"), nil)
	is.Equal(err.Error(), "one")
}

func TestOneIs(t *testing.T) {
	is := is.New(t)
	err := fmt.Errorf("one: %w", fs.ErrNotExist)
	is.True(errors.Is(errs.Join(nil, err, nil), fs.ErrNotExist))
}

// Can't assume multiple non-nested errors are a certain type of error.
func TestMultiIsNot(t *testing.T) {
	is := is.New(t)
	e1 := fmt.Errorf("one: %w", fs.ErrNotExist)
	e2 := fmt.Errorf("two: %w", fs.ErrNotExist)
	is.True(!errors.Is(errs.Join(e1, e2), fs.ErrNotExist))
}

func TestMulti(t *testing.T) {
	is := is.New(t)
	e1 := errors.New("one")
	e2 := errors.New("two")
	is.Equal(errs.Join(nil, e1, e2).Error(), "one. two")
}

func TestCompose(t *testing.T) {
	is := is.New(t)
	e1 := errors.New("one")
	e2 := errors.New("two")
	e3 := errs.Join(e1, e2)
	e4 := errors.New("four")
	e5 := errs.Join(e3, nil, e4)
	is.Equal(e5.Error(), "one. two. four")
}

func TestFuncNil(t *testing.T) {
	is := is.New(t)
	fn := func() (err error) {
		return errs.Join(err)
	}
	is.NoErr(fn())
}

func TestFuncLoop(t *testing.T) {
	is := is.New(t)
	closers := []func() error{
		func() error { return nil },
		func() error { return errors.New("b") },
		func() error { return errors.New("c") },
	}
	fn := func() (err error) {
		for _, close := range closers {
			err = errs.Join(err, close())
		}
		return err
	}
	err := fn()
	is.True(err != nil)
	is.Equal(err.Error(), "b. c")
}
