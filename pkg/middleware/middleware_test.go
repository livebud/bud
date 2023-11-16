package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestMiddleware(t *testing.T) {
	is := is.New(t)
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-A", "a")
			next.ServeHTTP(w, r)
			w.Header().Set("X-B", "b")
		})
	}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(w.Header().Get("X-A"), "a")
		is.Equal(w.Header().Get("X-B"), "")
		w.Write([]byte("hello"))
	})
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	middleware(h).ServeHTTP(rec, req)
	is.Equal(rec.Header().Get("X-A"), "a")
	is.Equal(rec.Header().Get("X-B"), "b")
	is.Equal(rec.Body.String(), "hello")
}
