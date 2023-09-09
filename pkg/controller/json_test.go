package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"

	"github.com/livebud/bud/pkg/di"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
)

func jsonEqual(t testing.TB, h http.Handler, request, expect string) {
	t.Helper()
	rec := httptest.NewRecorder()
	parts := strings.Split(request, " ")
	if len(parts) != 2 {
		t.Fatalf("invalid request: %s", request)
	}
	u, err := url.Parse(parts[1])
	if err != nil {
		t.Fatalf("invalid request: %s", request)
	}
	req := httptest.NewRequest(parts[0], u.Path, nil)
	req.URL.RawQuery = u.RawQuery
	h.ServeHTTP(rec, req)
	res := rec.Result()
	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		if err.Error() != expect {
			t.Fatalf("unexpected error: %v", err)
		}
		return
	}
	diff.TestHTTP(t, expect, string(dump))
}

func TestJsonHandlerFunc(t *testing.T) {
	is := is.New(t)
	h, err := wrapJSON(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		hello
	`)
}

func TestJsonHandler(t *testing.T) {
	is := is.New(t)
	h, err := wrapJSON(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))
	is.NoErr(err)
	jsonEqual(t, h, "GET /", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		hello
	`)
}

func TestJsonInvalid(t *testing.T) {
	is := is.New(t)
	h, err := wrapJSON(1)
	is.True(err != nil)
	is.True(errors.Is(err, ErrInvalidHandler))
	is.Equal(err.Error(), `rpc: "int" is an invalid handler type`)
	is.Equal(h, nil)
}

func TestJsonInvalidFunc(t *testing.T) {
	is := is.New(t)
	h, err := wrapJSON(func(a chan int) error {
		return nil
	})
	is.True(err != nil)
	is.True(errors.Is(err, ErrInvalidHandler))
	is.Equal(err.Error(), `rpc: "func(chan int) error" is an invalid handler type`)
	is.Equal(h, nil)
}

func TestJsonErrorOk(t *testing.T) {
	is := is.New(t)
	h, err := wrapJSON(func() error {
		return nil
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /", `
		HTTP/1.1 204 No Content
		Connection: close
	`)
}

func TestJsonErrorNotOk(t *testing.T) {
	is := is.New(t)
	h, err := wrapJSON(func() error {
		return errors.New("oh noz")
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /", `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"oh noz"}
	`)
}

func TestJsonContext(t *testing.T) {
	is := is.New(t)
	called := false
	h, err := wrapJSON(func(ctx context.Context) {
		is.True(ctx != nil)
		called = true
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /", `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	is.True(called)
}

func TestJsonContextErrorOk(t *testing.T) {
	is := is.New(t)
	called := false
	h, err := wrapJSON(func(ctx context.Context) error {
		is.True(ctx != nil)
		called = true
		return nil
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /", `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	is.True(called)
}

func TestJsonContextErrorNotOk(t *testing.T) {
	is := is.New(t)
	called := false
	h, err := wrapJSON(func(ctx context.Context) error {
		is.True(ctx != nil)
		called = true
		return fmt.Errorf("oh noz")
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /", `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"oh noz"}
	`)
	is.True(called)
}

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func TestJsonContextErrorCustomError(t *testing.T) {
	is := is.New(t)
	called := false

	h, err := wrapJSON(func(ctx context.Context) error {
		is.True(ctx != nil)
		called = true
		return &customError{"oh noz"}
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /", `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"oh noz"}
	`)
	is.True(called)
}

type statusError struct {
	msg string
}

func (s *statusError) Status() int {
	return http.StatusBadRequest
}

func (s *statusError) Error() string {
	return s.msg
}

func TestJsonContextErrorCustomStatus(t *testing.T) {
	is := is.New(t)
	called := false
	h, err := wrapJSON(func(ctx context.Context) error {
		is.True(ctx != nil)
		called = true
		return &statusError{"oh noz"}
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /", `
		HTTP/1.1 400 Bad Request
		Connection: close
		Content-Type: application/json

		{"error":"oh noz"}
	`)
	is.True(called)
}

func TestJsonStructIn(t *testing.T) {
	is := is.New(t)
	called := false
	type In struct {
		Message string `json:"message"`
	}
	h, err := wrapJSON(func(in In) {
		is.Equal(in.Message, "hi")
		called = true
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /?message=hi", `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	is.True(called)
}

func TestJsonStructPtrIn(t *testing.T) {
	is := is.New(t)
	called := false
	type In struct {
		Message string `json:"message"`
	}
	h, err := wrapJSON(func(in *In) {
		is.Equal(in.Message, "hi")
		called = true
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /?message=hi", `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	is.True(called)
}

func TestJsonInjectedStruct(t *testing.T) {
	is := is.New(t)
	called := false
	in := di.New()
	type Env struct {
		DatabaseURL string
	}
	type Log struct {
		Env *Env
	}
	di.Provide[*Env](in, func() (*Env, error) {
		return &Env{"http://localhost:5432"}, nil
	})
	di.Provide[*Log](in, func(env *Env) (*Log, error) {
		return &Log{env}, nil
	})
	type Context struct {
		context.Context
		Log *Log
	}
	h, err := wrapJSON(func(ctx *Context) {
		is.True(ctx != nil)
		is.True(ctx.Log != nil)
		is.Equal(ctx.Log.Env.DatabaseURL, "http://localhost:5432")
		called = true
	})
	is.NoErr(err)
	jsonEqual(t, di.Middleware(in)(h), "GET /?message=hi", `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	is.True(called)
}

func TestJsonContextResponseWriter(t *testing.T) {
	is := is.New(t)
	called := false
	h, err := wrapJSON(func(ctx context.Context, w http.ResponseWriter) {
		is.True(ctx != nil)
		is.True(w != nil)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("hi"))
		called = true
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain

		hi
	`)
	is.True(called)
}

func TestJsonRequestResponseWriter(t *testing.T) {
	is := is.New(t)
	h, err := wrapJSON(func(r *http.Request, w http.ResponseWriter) {
	})
	is.True(err != nil)
	is.True(errors.Is(err, ErrInvalidHandler))
	is.True(h == nil)
}

func TestJsonCtxStructIn(t *testing.T) {
	is := is.New(t)
	called := false
	type In struct {
		Message string `json:"message"`
	}
	h, err := wrapJSON(func(ctx context.Context, in In) {
		is.Equal(in.Message, "hi")
		called = true
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /?message=hi", `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	is.True(called)
}

func TestJsonCtxStructPtrIn(t *testing.T) {
	is := is.New(t)
	called := false
	type In struct {
		Message string `json:"message"`
	}
	h, err := wrapJSON(func(ctx context.Context, in *In) {
		is.Equal(in.Message, "hi")
		called = true
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /?message=hi", `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	is.True(called)
}

func TestJsonOutStruct(t *testing.T) {
	is := is.New(t)
	called := false
	type Out struct {
		Message string `json:"message"`
	}
	h, err := wrapJSON(func() Out {
		called = true
		return Out{"hi"}
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /?message=hi", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"message":"hi"}
	`)
	is.True(called)
}

func TestJsonOutPtrStruct(t *testing.T) {
	is := is.New(t)
	called := false
	type Out struct {
		Message string `json:"message"`
	}
	h, err := wrapJSON(func() *Out {
		called = true
		return &Out{"hi"}
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /?message=hi", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"message":"hi"}
	`)
	is.True(called)
}

func TestJsonOutStructError(t *testing.T) {
	is := is.New(t)
	called := 0
	type Out struct {
		Message string `json:"message"`
	}
	h, err := wrapJSON(func() (Out, error) {
		called++
		if called > 1 {
			return Out{}, errors.New("oh noz")
		}
		return Out{"hi"}, nil
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /?message=hi", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"message":"hi"}
	`)
	jsonEqual(t, h, "GET /?message=hi", `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"oh noz"}
	`)
	is.Equal(called, 2)
}

func TestJsonOutPtrStructError(t *testing.T) {
	is := is.New(t)
	called := 0
	type Out struct {
		Message string `json:"message"`
	}
	h, err := wrapJSON(func() (*Out, error) {
		called++
		if called > 1 {
			return nil, errors.New("oh noz")
		}
		return &Out{"hi"}, nil
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /?message=hi", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"message":"hi"}
	`)
	jsonEqual(t, h, "GET /?message=hi", `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"oh noz"}
	`)
	is.Equal(called, 2)
}

func TestJsonInError(t *testing.T) {
	is := is.New(t)
	called := 0
	type In struct {
		Message string `json:"message"`
	}
	h, err := wrapJSON(func(in In) error {
		called++
		if in.Message != "hi" {
			return errors.New("oh noz")
		}
		return nil
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /?message=hi", `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	jsonEqual(t, h, "GET /?message=cool", `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"oh noz"}
	`)
	is.Equal(called, 2)
}

func TestJsonPtrInError(t *testing.T) {
	is := is.New(t)
	called := 0
	type In struct {
		Message string `json:"message"`
	}
	h, err := wrapJSON(func(in *In) error {
		called++
		if in.Message != "hi" {
			return errors.New("oh noz")
		}
		return nil
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /?message=hi", `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	jsonEqual(t, h, "GET /?message=cool", `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"oh noz"}
	`)
	is.Equal(called, 2)
}

func TestJsonListIn(t *testing.T) {
	is := is.New(t)
	called := 0
	type In struct {
		Messages []string `json:"messages"`
	}
	h, err := wrapJSON(func(in *In) *In {
		called++
		return in
	})
	is.NoErr(err)
	jsonEqual(t, h, "GET /?messages.0=a&messages.1=b&messages.2=c", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"messages":["a","b","c"]}
	`)
	is.Equal(called, 1)
}
