package imports_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/imports"
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
