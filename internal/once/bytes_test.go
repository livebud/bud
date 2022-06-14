package once_test

import (
	"errors"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/once"
)

func TestBytesNil(t *testing.T) {
	is := is.New(t)
	var once once.Bytes
	called := 0
	res, err := once.Do(func() ([]byte, error) {
		called++
		return []byte("1"), nil
	})
	is.NoErr(err)
	is.Equal(called, 1)
	is.Equal(res, []byte("1"))
	res, err = once.Do(func() ([]byte, error) {
		called++
		return []byte("2"), errors.New("oh noz")
	})
	is.NoErr(err)
	is.Equal(called, 1)
	is.Equal(res, []byte("1"))
}

func TestBytesError(t *testing.T) {
	is := is.New(t)
	var once once.Bytes
	called := 0
	res, err := once.Do(func() ([]byte, error) {
		called++
		return []byte("1"), errors.New("oh noz")
	})
	is.True(err != nil)
	is.Equal(err.Error(), "oh noz")
	is.Equal(called, 1)
	is.Equal(res, []byte("1"))
	res, err = once.Do(func() ([]byte, error) {
		called++
		return []byte("2"), nil
	})
	is.True(err != nil)
	is.Equal(err.Error(), "oh noz")
	is.Equal(called, 1)
	is.Equal(res, []byte("1"))
}
