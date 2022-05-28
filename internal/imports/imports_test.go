package imports_test

import (
	"testing"

	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/internal/is"
)

func TestAdd(t *testing.T) {
	is := is.New(t)
	im := imports.New()
	is.Equal(im.Add("net/http"), "http")
	is.Equal(im.Add("net/http"), "http")
	is.Equal(im.Add("hop/http"), "http1")
}

func TestAddNamed(t *testing.T) {
	is := is.New(t)
	im := imports.New()
	is.Equal(im.AddNamed("www", "net/http"), "www")
	is.Equal(im.AddNamed("www", "net/http"), "www")
	is.Equal(im.AddNamed("www", "hop/http"), "www1")
	is.Equal(im.AddNamed("v8", "app.com/js/v8"), "v8")
}

func TestReserveBefore(t *testing.T) {
	is := is.New(t)
	im := imports.New()
	is.Equal(im.Reserve("web"), "web")
	is.Equal(len(im.List()), 0)
	is.Equal(im.Add("web"), "web")
	is.Equal(len(im.List()), 1)
}
func TestReserveAfter(t *testing.T) {
	is := is.New(t)
	im := imports.New()
	is.Equal(im.Add("web"), "web")
	is.Equal(len(im.List()), 1)
	is.Equal(im.Reserve("web"), "web")
	is.Equal(len(im.List()), 1)
	is.Equal(im.Reserve("duo/web"), "web1")
	is.Equal(len(im.List()), 1)
	is.Equal(im.Add("duo/web"), "web1")
	is.Equal(len(im.List()), 2)
}

func TestAddStd(t *testing.T) {
	is := is.New(t)
	im := imports.New()
	im.AddStd("os", "fmt", "net/http")
	is.Equal(len(im.List()), 3)
	is.Equal(im.List()[0].Name, "fmt")
	is.Equal(im.List()[0].Path, "fmt")
	is.Equal(im.List()[1].Name, "http")
	is.Equal(im.List()[1].Path, "net/http")
	is.Equal(im.List()[2].Name, "os")
	is.Equal(im.List()[2].Path, "os")
}
