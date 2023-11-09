package public_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/pkg/public"
	"github.com/matryer/is"
)

// Pulled from: https://github.com/mathiasbynens/small
// Built with: xxd -i small.ico
var favicon = []byte{
	0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00,
	0x18, 0x00, 0x30, 0x00, 0x00, 0x00, 0x16, 0x00, 0x00, 0x00, 0x28, 0x00,
	0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00,
	0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00,
}

// Small valid gif: https://github.com/mathiasbynens/small/blob/master/gif.gif
var gif = []byte{
	0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01,
	0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x3b,
}

func TestPublic(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"favicon.ico": &fstest.MapFile{
			Data: favicon,
		},
		"ga.js": &fstest.MapFile{
			Data: []byte(`function ga(track){}`),
		},
		"normalize/normalize.css": &fstest.MapFile{
			Data: []byte(`* { box-sizing: border-box; }`),
		},
		"lol.gif": &fstest.MapFile{
			Data: gif,
		},
	}
	public := public.New(fsys)
	req := httptest.NewRequest("GET", "/favicon.ico", nil)
	rec := httptest.NewRecorder()
	public.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(200, res.StatusCode)
	actual, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(actual, favicon)
	// Ubuntu CI reports a different MIME type than OSX
	is.True(strings.Contains(res.Header.Get("Content-Type"), "image/"))
	is.True(strings.Contains(res.Header.Get("Content-Type"), "icon"))

	// /ga.js
	req = httptest.NewRequest("GET", "/ga.js", nil)
	rec = httptest.NewRecorder()
	public.ServeHTTP(rec, req)
	res = rec.Result()
	is.Equal(200, res.StatusCode)
	actual, err = io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(actual, []byte(`function ga(track){}`))
	is.True(strings.Contains(res.Header.Get("Content-Type"), "/javascript"))

	// /normalize/normalize.css
	req = httptest.NewRequest("GET", "/normalize/normalize.css", nil)
	rec = httptest.NewRecorder()
	public.ServeHTTP(rec, req)
	res = rec.Result()
	is.Equal(200, res.StatusCode)
	actual, err = io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(actual, []byte(`* { box-sizing: border-box; }`))
	is.True(strings.Contains(res.Header.Get("Content-Type"), "css"))

	// /lol.gif
	req = httptest.NewRequest("GET", "/lol.gif", nil)
	rec = httptest.NewRecorder()
	public.ServeHTTP(rec, req)
	res = rec.Result()
	is.Equal(200, res.StatusCode)
	actual, err = io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(actual, gif)
	is.True(strings.Contains(res.Header.Get("Content-Type"), "image/"))
	is.True(strings.Contains(res.Header.Get("Content-Type"), "gif"))

	// 404.html
	req = httptest.NewRequest("GET", "/404.html", nil)
	rec = httptest.NewRecorder()
	public.ServeHTTP(rec, req)
	res = rec.Result()
	is.Equal(404, res.StatusCode)
	actual, err = io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(actual, []byte("404 page not found\n"))
}
