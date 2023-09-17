package session_test

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/session"
	"github.com/livebud/bud/pkg/session/cookiestore"
	"github.com/livebud/bud/pkg/session/internal/cookies"
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

func TestSession(t *testing.T) {
	is := is.New(t)
	jar, err := cookiejar.New(nil)
	is.NoErr(err)
	// cookies := cookies.New()
	// store := cookiestore.New(cookies)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// store.Load()
		// cookie, err := cookies.Read(r, "sid")
		// is.NoErr(err)
		// is.Equal(cookie.Value, "testId")
		// session, err := store.Load(r.Context(), cookie.Value)
		// is.NoErr(err)
		// is.Equal(session.ID, "testId")
		// is.Equal(session.Data, map[string]interface{}{})
		// session.Data["foo"] = "bar"
	})
	req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	equal(t, jar, handler, req, ``)
}

func TestSessionCounter(t *testing.T) {
	is := is.New(t)
	jar, err := cookiejar.New(nil)
	is.NoErr(err)
	type Session struct {
		Visits int `json:"visits"`
	}
	cookies := cookies.New()
	store := cookiestore.New(cookies)
	store.GenerateID = func() (string, error) {
		return "testId", nil
	}
	middleware := session.New(cookies, store)
	lastValue := -1
	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := session.From(r.Context())
		is.NoErr(err)
		visits, ok := session.Get("visits").(int)
		if !ok {
			visits = 0
		}
		is.Equal(visits, lastValue+1)
		visits++
		lastValue++
		session.Set("visits", visits)
	}))
	req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	equal(t, jar, handler, req, `
		HTTP/1.1 200 OK
		Connection: close
		Set-Cookie: testId=CN7QIAIBARCGC5DBAH7YAAABBQARAAAACP7YAAABAZ3GS43JORZQG2LOOQCAEAAC
		Set-Cookie: sid=testId; HttpOnly
	`)
	req = httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	equal(t, jar, handler, req, `
		HTTP/1.1 200 OK
		Connection: close
		Set-Cookie: testId=CN7QIAIBARCGC5DBAH7YAAABBQARAAAACP7YAAABAZ3GS43JORZQG2LOOQCAEAAE
		Set-Cookie: sid=testId; HttpOnly
	`)
	req = httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	equal(t, jar, handler, req, `
		HTTP/1.1 200 OK
		Connection: close
		Set-Cookie: testId=CN7QIAIBARCGC5DBAH7YAAABBQARAAAACP7YAAABAZ3GS43JORZQG2LOOQCAEAAG
		Set-Cookie: sid=testId; HttpOnly
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
	cookies := cookies.New()
	middleware := session.New(cookies, cookiestore.New(cookies))
	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := session.From(r.Context())
		is.NoErr(err)
		fmt.Println("got visits", session.Get("visits"))
		// session.Visits++
		// if session.Visits == 2 {
		// 	session.User = &User{ID: 1}
		// }
		// if session.Visits == 4 {
		// 	session.User = nil
		// }
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
