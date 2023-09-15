package session_test

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/session"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
)

func equal(t testing.TB, jar *cookiejar.Jar, h http.Handler, r *http.Request, expect string) {
	t.Helper()
	for _, cookie := range jar.Cookies(r.URL) {
		r.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	w := rec.Result()
	jar.SetCookies(r.URL, w.Cookies())
	dump, err := httputil.DumpResponse(w, true)
	if err != nil {
		if err.Error() != expect {
			t.Fatalf("unexpected error: %v", err)
		}
		return
	}
	diff.TestHTTP(t, expect, string(dump))
}

func TestSetGetCookie(t *testing.T) {
	is := is.New(t)
	jar, err := cookiejar.New(nil)
	is.NoErr(err)
	router := mux.New()
	router.Get("/set", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "cookie_name", Value: "cookie_value"})
	}))
	router.Get("/get", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("cookie_name")
		is.NoErr(err)
		http.SetCookie(w, cookie)
	}))
	req := httptest.NewRequest(http.MethodGet, "http://example.com/set", nil)
	equal(t, jar, router, req, `
		HTTP/1.1 200 OK
		Connection: close
		Set-Cookie: cookie_name=cookie_value
	`)
	req = httptest.NewRequest(http.MethodGet, "http://example.com/get", nil)
	equal(t, jar, router, req, `
		HTTP/1.1 200 OK
		Connection: close
		Set-Cookie: cookie_name=cookie_value
	`)
}

func TestSessionCounter(t *testing.T) {
	is := is.New(t)
	jar, err := cookiejar.New(nil)
	is.NoErr(err)
	type Session struct {
		Visits int `json:"visits"`
	}
	middleware := session.New()
	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := session.From[*Session](r.Context())
		is.NoErr(err)
		session.Visits++
	}))
	req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	equal(t, jar, handler, req, `
		HTTP/1.1 200 OK
		Connection: close
		Set-Cookie: sid=visits=1; HttpOnly
	`)
	req = httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	equal(t, jar, handler, req, `
		HTTP/1.1 200 OK
		Connection: close
		Set-Cookie: sid=visits=2; HttpOnly
	`)
	req = httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	equal(t, jar, handler, req, `
		HTTP/1.1 200 OK
		Connection: close
		Set-Cookie: sid=visits=3; HttpOnly
	`)
}

func TestSessionNested(t *testing.T) {
	is := is.New(t)
	jar, err := cookiejar.New(nil)
	is.NoErr(err)
	type User struct {
		ID int `json:"id"`
	}
	type Session struct {
		Visits int   `json:"visits"`
		User   *User `json:"user,omitempty"`
	}
	middleware := session.New()
	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := session.From[*Session](r.Context())
		is.NoErr(err)
		session.Visits++
		if session.Visits == 2 {
			session.User = &User{ID: 1}
		}
		if session.Visits == 4 {
			session.User = nil
		}
	}))
	req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	equal(t, jar, handler, req, `
		HTTP/1.1 200 OK
		Connection: close
		Set-Cookie: sid=visits=1; HttpOnly
	`)
	req = httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	equal(t, jar, handler, req, `
		HTTP/1.1 200 OK
		Connection: close
		Set-Cookie: sid=user.id=1&visits=2; HttpOnly
	`)
	req = httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	equal(t, jar, handler, req, `
		HTTP/1.1 200 OK
		Connection: close
		Set-Cookie: sid=user.id=1&visits=3; HttpOnly
	`)
	req = httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	equal(t, jar, handler, req, `
		HTTP/1.1 200 OK
		Connection: close
		Set-Cookie: sid=visits=4; HttpOnly
	`)
}
