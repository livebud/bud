package once_test

import (
	"errors"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/once"
)

func TestStringNil(t *testing.T) {
	is := is.New(t)
	called := 0
	callOnce := once.String(func() (string, error) {
		called++
		return "ok", nil
	})
	res, err := callOnce()
	is.NoErr(err)
	is.Equal(called, 1)
	is.Equal(res, "ok")
	res, err = callOnce()
	is.NoErr(err)
	is.Equal(called, 1)
	is.Equal(res, "ok")
}

func TestStringError(t *testing.T) {
	is := is.New(t)
	called := 0
	callOnce := once.String(func() (string, error) {
		called++
		return "ok", errors.New("oh noz")
	})
	res, err := callOnce()
	is.True(err != nil)
	is.Equal(err.Error(), "oh noz")
	is.Equal(called, 1)
	is.Equal(res, "ok")
	res, err = callOnce()
	is.True(err != nil)
	is.Equal(err.Error(), "oh noz")
	is.Equal(called, 1)
	is.Equal(res, "ok")
}
