package httplog

import (
	"fmt"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/middleware"
	"github.com/segmentio/ksuid"
)

// WithRequestId sets the request id function for generating a unique request id
// for each request
func WithRequestId(fn func(r *http.Request) string) func(*middlewareOption) {
	return func(opts *middlewareOption) {
		opts.requestId = fn
	}
}

type middlewareOption struct {
	requestId func(r *http.Request) string
}

// RequestId is a function for generating a unique request id
func defaultRequestId(r *http.Request) string {
	// Support an existing request id
	requestId := r.Header.Get("X-Request-Id")
	if requestId == "" {
		requestId = ksuid.New().String()
		// Set just in case we use it later
		r.Header.Set("X-Request-Id", requestId)
	}
	return requestId
}

// New uses the logger to log requests and responses
func New(log logs.Log, options ...func(*middlewareOption)) middleware.Middleware {
	opts := &middlewareOption{
		requestId: defaultRequestId,
	}
	for _, option := range options {
		option(opts)
	}
	return middleware.Func(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := log.Fields(logs.Fields{
				"url":         r.RequestURI,
				"method":      r.Method,
				"remote_addr": r.RemoteAddr,
				"request_id":  opts.requestId(r),
			})
			ctx := logs.ToContext(r.Context(), log)
			r = r.WithContext(ctx)
			log.Info("request")
			res := httpsnoop.CaptureMetrics(next, w, r)
			log = log.Fields(logs.Fields{
				"status":   res.Code,
				"duration": fmt.Sprintf("%dms", res.Duration.Milliseconds()),
				"size":     res.Written,
			})
			switch {
			case res.Code >= 500:
				log.Error("response")
			case res.Code >= 400:
				log.Warn("response")
			default:
				log.Info("response")
			}
		})
	})
}
