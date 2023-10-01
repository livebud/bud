package slots_test

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/livebud/bud/view/slots"
	"github.com/matryer/is"
	"golang.org/x/sync/errgroup"
)

func TestSingle(t *testing.T) {
	is := is.New(t)
	w := new(bytes.Buffer)
	layout := slots.New(w)
	layout.Write([]byte("layout"))
	is.NoErr(layout.Close(nil))
	is.Equal(w.String(), "layout")
}

func wrap(left string, middle []byte, right string) []byte {
	return append(append([]byte(left), middle...), []byte(right)...)
}

func TestChain(t *testing.T) {
	is := is.New(t)
	w := new(bytes.Buffer)
	eg := new(errgroup.Group)
	mu := sync.Mutex{}
	called := []string{}
	layout := slots.New(w)
	eg.Go(func() (err error) {
		defer layout.Close(&err)
		data, err := io.ReadAll(layout)
		if err != nil {
			return err
		}
		wrapped := wrap("<l>", data, "</l>")
		n, err := layout.Write(wrapped)
		if err != nil {
			return err
		}
		is.Equal(n, len(wrapped))
		mu.Lock()
		called = append(called, "layout")
		mu.Unlock()
		return nil
	})
	frame1 := layout.New()
	eg.Go(func() (err error) {
		defer frame1.Close(&err)
		data, err := io.ReadAll(frame1)
		if err != nil {
			return err
		}
		wrapped := wrap("<f1>", data, "</f1>")
		n, err := frame1.Write(wrapped)
		if err != nil {
			return err
		}
		is.Equal(n, len(wrapped))
		mu.Lock()
		called = append(called, "frame1")
		mu.Unlock()
		return nil
	})
	frame2 := frame1.New()
	eg.Go(func() (err error) {
		defer frame2.Close(&err)
		data, err := io.ReadAll(frame2)
		if err != nil {
			return err
		}
		wrapped := wrap("<f2>", data, "</f2>")
		n, err := frame2.Write(wrapped)
		if err != nil {
			return err
		}
		is.Equal(n, len(wrapped))
		mu.Lock()
		called = append(called, "frame2")
		mu.Unlock()
		return nil
	})
	view := frame2.New()
	eg.Go(func() (err error) {
		defer view.Close(&err)
		data, err := io.ReadAll(view)
		if err != nil {
			return err
		}
		is.Equal(string(data), "")
		wrapped := wrap("<v>", data, "</v>")
		n, err := view.Write(wrapped)
		if err != nil {
			return err
		}
		is.Equal(n, len(wrapped))
		mu.Lock()
		called = append(called, "view")
		mu.Unlock()
		return nil
	})
	is.NoErr(eg.Wait())
	is.Equal(called, []string{"view", "frame2", "frame1", "layout"})
	is.Equal(w.String(), "<l><f1><f2><v></v></f2></f1></l>")
}

func TestChainWithError(t *testing.T) {
	is := is.New(t)
	w := new(bytes.Buffer)
	eg := new(errgroup.Group)
	called := []string{}
	layout := slots.New(w)
	eg.Go(func() (err error) {
		defer layout.Close(&err)
		data, err := io.ReadAll(layout)
		if err != nil {
			return err
		}
		wrapped := wrap("<l>", data, "</l>")
		n, err := layout.Write(wrapped)
		if err != nil {
			return err
		}
		is.Equal(n, len(wrapped))
		called = append(called, "layout")
		return nil
	})
	view := layout.New()
	eg.Go(func() (err error) {
		defer view.Close(&err)
		data, err := io.ReadAll(view)
		if err != nil {
			return err
		}
		is.Equal(string(data), "")
		called = append(called, "view")
		return fmt.Errorf("Oh noz!")
	})
	err := eg.Wait()
	is.True(err != nil)
	is.Equal(err.Error(), "Oh noz!")
	is.Equal(called, []string{"view"})
}

func TestFrameWithError(t *testing.T) {
	is := is.New(t)
	w := new(bytes.Buffer)
	eg := new(errgroup.Group)
	called := []string{}
	layout := slots.New(w)
	eg.Go(func() (err error) {
		defer layout.Close(&err)
		data, err := io.ReadAll(layout)
		if err != nil {
			return err
		}
		wrapped := wrap("<l>", data, "</l>")
		n, err := layout.Write(wrapped)
		if err != nil {
			return err
		}
		is.Equal(n, len(wrapped))
		called = append(called, "layout")
		return nil
	})
	frame := layout.New()
	eg.Go(func() (err error) {
		defer frame.Close(&err)
		data, err := io.ReadAll(frame)
		if err != nil {
			return err
		}
		is.Equal(string(data), "<v></v>")
		called = append(called, "frame")
		return fmt.Errorf("Oh noz!")
	})
	view := frame.New()
	eg.Go(func() (err error) {
		defer view.Close(&err)
		time.Sleep(100 * time.Millisecond)
		data, err := io.ReadAll(view)
		if err != nil {
			return err
		}
		is.Equal(string(data), "")
		view.Write([]byte("<v></v>"))
		called = append(called, "view")
		return nil
	})
	err := eg.Wait()
	is.True(err != nil)
	is.Equal(err.Error(), "Oh noz!")
	is.Equal(called, []string{"view", "frame"})
}

func TestFrameErrorNoRead(t *testing.T) {
	is := is.New(t)
	w := new(bytes.Buffer)
	eg := new(errgroup.Group)
	called := []string{}
	layout := slots.New(w)
	eg.Go(func() (err error) {
		defer layout.Close(&err)
		data, err := io.ReadAll(layout)
		if err != nil {
			return err
		}
		wrapped := wrap("<l>", data, "</l>")
		n, err := layout.Write(wrapped)
		if err != nil {
			return err
		}
		is.Equal(n, len(wrapped))
		called = append(called, "layout")
		return nil
	})
	frame := layout.New()
	eg.Go(func() (err error) {
		defer frame.Close(&err)
		return fmt.Errorf("Oh noz!")
	})
	view := frame.New()
	eg.Go(func() (err error) {
		defer view.Close(&err)
		time.Sleep(100 * time.Millisecond)
		data, err := io.ReadAll(view)
		if err != nil {
			return err
		}
		is.Equal(string(data), "")
		view.Write([]byte("<v></v>"))
		return nil
	})
	err := eg.Wait()
	is.True(err != nil)
	is.Equal(err.Error(), "Oh noz!")
}
