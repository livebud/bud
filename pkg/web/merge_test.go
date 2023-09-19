package web_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/livebud/bud/pkg/web"
	"github.com/matryer/is"
)

func TestMerge(t *testing.T) {
	is := is.New(t)
	merged := web.Merge(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("a", "aa")
			w.Write([]byte("a"))
		}),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("b", "bb")
			w.Write([]byte("b"))
		}),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("b", "cc")
			w.Write([]byte("c"))
		}),
	)
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	merged.ServeHTTP(w, r)
	res := w.Result()
	is.Equal(res.StatusCode, http.StatusOK)
	is.Equal(res.Header.Get("a"), "aa")
	is.Equal(res.Header.Get("b"), "bb")
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "cba")
}
