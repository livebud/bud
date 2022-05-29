package once_test

import (
	"errors"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/once"
)

func TestErrorNil(t *testing.T) {
	is := is.New(t)
	var once once.Error
	called := 0
	err := once.Do(func() error {
		called++
		return nil
	})
	is.NoErr(err)
	is.Equal(called, 1)
	err = once.Do(func() error {
		called++
		return errors.New("oh noz")
	})
	is.NoErr(err)
	is.Equal(called, 1)
}

func TestError(t *testing.T) {
	is := is.New(t)
	var once once.Error
	called := 0
	err := once.Do(func() error {
		called++
		return errors.New("oh noz")
	})
	is.True(err != nil)
	is.Equal(err.Error(), "oh noz")
	is.Equal(called, 1)
	err = once.Do(func() error {
		called++
		return err
	})
	is.True(err != nil)
	is.Equal(err.Error(), "oh noz")
	is.Equal(called, 1)
}
