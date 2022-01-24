package singleflight_test

import (
	"io/fs"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/2/singleflight"
)

func TestSingleFlight(t *testing.T) {
	is := is.New(t)
	loader := new(singleflight.Loader)
	wg := new(sync.WaitGroup)
	var calls int32 = 0
	fsys := &osfs{&calls}
	wg.Add(2)
	go func() {
		defer wg.Done()
		file, err := loader.Load(fsys, "singleflight_test.go")
		is.NoErr(err)
		fi, err := file.Stat()
		is.NoErr(err)
		is.Equal(fi.Name(), "singleflight_test.go")
	}()
	go func() {
		defer wg.Done()
		file, err := loader.Load(fsys, "singleflight_test.go")
		is.NoErr(err)
		fi, err := file.Stat()
		is.NoErr(err)
		is.Equal(fi.Name(), "singleflight_test.go")
	}()
	wg.Wait()
	is.Equal(calls, int32(1))
}

type osfs struct {
	calls *int32
}

func (m *osfs) Open(name string) (fs.File, error) {
	atomic.AddInt32(m.calls, 1)
	time.Sleep(10 * time.Millisecond)
	return os.Open("singleflight_test.go")
}
