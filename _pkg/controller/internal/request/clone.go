package request

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

func Clone(ctx context.Context, r *http.Request, body []byte) *http.Request {
	req := r.Clone(ctx)
	req.Body = io.NopCloser(bytes.NewReader(body))
	return req
}
