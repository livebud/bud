package once_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/livebud/bud/pkg/once"
	"github.com/matryer/is"
)

func TestBytesNil(t *testing.T) {
	is := is.New(t)
	var o once.Bytes
	called := 0
	res, err := o.Do(func() ([]byte, error) {
		called++
		return []byte("1"), nil
	})
	is.NoErr(err)
	is.Equal(called, 1)
	is.Equal(res, []byte("1"))
	res, err = o.Do(func() ([]byte, error) {
		called++
		return []byte("2"), errors.New("oh noz")
	})
	is.NoErr(err)
	is.Equal(called, 1)
	is.Equal(res, []byte("1"))
}

func TestBytesError(t *testing.T) {
	is := is.New(t)
	var o once.Bytes
	called := 0
	res, err := o.Do(func() ([]byte, error) {
		called++
		return []byte("1"), errors.New("oh noz")
	})
	is.True(err != nil)
	is.Equal(err.Error(), "oh noz")
	is.Equal(called, 1)
	is.Equal(res, []byte("1"))
	res, err = o.Do(func() ([]byte, error) {
		called++
		return []byte("2"), nil
	})
	is.True(err != nil)
	is.Equal(err.Error(), "oh noz")
	is.Equal(called, 1)
	is.Equal(res, []byte("1"))
}

func ExampleFunc() {
	count := 0
	fn := once.Func(func() (int, error) {
		count++
		return count, nil
	})
	res, err := fn()
	fmt.Println(res, err)
	res, err = fn()
	fmt.Println(res, err)
	// Output:
	// 1 <nil>
	// 1 <nil>
}
