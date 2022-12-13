package once_test

import (
	"errors"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/once"
)

func TestCloserOk(t *testing.T) {
	is := is.New(t)
	var closer once.Closer
	a := func() error { return nil }
	b := func() error { return nil }
	closer.Add(a)
	closer.Add(b)
	err := closer.Close()
	is.NoErr(err)
}

func TestCloserErrors(t *testing.T) {
	is := is.New(t)
	e1 := errors.New("error 1")
	e2 := errors.New("error 2")
	var closer once.Closer
	a := func() error { return e1 }
	b := func() error { return e2 }
	closer.Add(a)
	closer.Add(b)
	err := closer.Close()
	is.True(err != nil)
	is.Equal(err.Error(), "error 2. error 1")
}
