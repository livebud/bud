package request_test

import (
	"net/http/httptest"
	"testing"

	"github.com/livebud/bud/pkg/request"
	"github.com/matryer/is"
)

func TestEmpty(t *testing.T) {
	is := is.New(t)
	r := httptest.NewRequest("GET", "/", nil)
	ctype, err := request.Negotiate(r, "text/html", "application/json")
	is.NoErr(err)
	is.Equal(ctype, "text/html")
	ctype, err = request.Negotiate(r, "application/json", "text/html")
	is.NoErr(err)
	is.Equal(ctype, "application/json")
}

func TestAcceptHTML(t *testing.T) {
	is := is.New(t)
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Accept", "text/html")
	ctype, err := request.Negotiate(r, "text/html", "application/json")
	is.NoErr(err)
	is.Equal(ctype, "text/html")
	ctype, err = request.Negotiate(r, "application/json", "text/html")
	is.NoErr(err)
	is.Equal(ctype, "text/html")
	ctype, err = request.Negotiate(r, "text/html")
	is.NoErr(err)
	is.Equal(ctype, "text/html")
	ctype, err = request.Negotiate(r, "application/json")
	is.NoErr(err)
	is.Equal(ctype, "")
}

func TestContentTypeHTML(t *testing.T) {
	is := is.New(t)
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "text/html")
	ctype, err := request.Negotiate(r, "text/html", "application/json")
	is.NoErr(err)
	is.Equal(ctype, "text/html")
	ctype, err = request.Negotiate(r, "application/json", "text/html")
	is.NoErr(err)
	is.Equal(ctype, "text/html")
	ctype, err = request.Negotiate(r, "text/html")
	is.NoErr(err)
	is.Equal(ctype, "text/html")
	ctype, err = request.Negotiate(r, "application/json")
	is.NoErr(err)
	is.Equal(ctype, "")
}

func TestAcceptJSON(t *testing.T) {
	is := is.New(t)
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Accept", "application/json")
	ctype, err := request.Negotiate(r, "text/html", "application/json")
	is.NoErr(err)
	is.Equal(ctype, "application/json")
	ctype, err = request.Negotiate(r, "application/json", "text/html")
	is.NoErr(err)
	is.Equal(ctype, "application/json")
	ctype, err = request.Negotiate(r, "text/html")
	is.NoErr(err)
	is.Equal(ctype, "")
	ctype, err = request.Negotiate(r, "application/json")
	is.NoErr(err)
	is.Equal(ctype, "application/json")
}

func TestContentTypeJSON(t *testing.T) {
	is := is.New(t)
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/json")
	ctype, err := request.Negotiate(r, "text/html", "application/json")
	is.NoErr(err)
	is.Equal(ctype, "application/json")
	ctype, err = request.Negotiate(r, "application/json", "text/html")
	is.NoErr(err)
	is.Equal(ctype, "application/json")
	ctype, err = request.Negotiate(r, "text/html")
	is.NoErr(err)
	is.Equal(ctype, "")
	ctype, err = request.Negotiate(r, "application/json")
	is.NoErr(err)
	is.Equal(ctype, "application/json")
}

func TestAccept(t *testing.T) {
	is := is.New(t)
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/json")
	ctype := request.Accept(r, "text/html", "application/json")
	is.Equal(ctype, "application/json")
	ctype = request.Accept(r, "application/json", "text/html")
	is.Equal(ctype, "application/json")
	ctype = request.Accept(r, "text/html")
	is.Equal(ctype, "")
	ctype = request.Accept(r, "application/json")
	is.Equal(ctype, "application/json")
}
