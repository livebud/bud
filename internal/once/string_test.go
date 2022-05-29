package once_test

import (
	"errors"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/once"
)

func TestStringNil(t *testing.T) {
	is := is.New(t)
	var once once.String
	called := 0
	res, err := once.Do(func() (string, error) {
		called++
		return "1", nil
	})
	is.NoErr(err)
	is.Equal(called, 1)
	is.Equal(res, "1")
	res, err = once.Do(func() (string, error) {
		called++
		return "2", errors.New("oh noz")
	})
	is.NoErr(err)
	is.Equal(called, 1)
	is.Equal(res, "1")
}

func TestStringError(t *testing.T) {
	is := is.New(t)
	var once once.String
	called := 0
	res, err := once.Do(func() (string, error) {
		called++
		return "1", errors.New("oh noz")
	})
	is.True(err != nil)
	is.Equal(err.Error(), "oh noz")
	is.Equal(called, 1)
	is.Equal(res, "1")
	res, err = once.Do(func() (string, error) {
		called++
		return "2", nil
	})
	is.True(err != nil)
	is.Equal(err.Error(), "oh noz")
	is.Equal(called, 1)
	is.Equal(res, "1")
}
