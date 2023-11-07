package res

import (
	"fmt"
	"net/http"

	"github.com/livebud/bud/log"
)

func New(log log.Log) *response {
	return &response{0, nil}
}

func Errorf(w http.ResponseWriter, msg string, v ...any) *response {
	return &response{0, fmt.Errorf(msg, v...)}
}

func JSON(v any) *response {
	return &response{0, nil}
}

func Html(w http.ResponseWriter, h string) *response {
	return &response{0, nil}
}

func Redirect(to string, v ...any) *response {
	return &response{0, nil}
}

func Redirects(w http.ResponseWriter, r *http.Request, to string, v ...any) *response {
	return &response{0, nil}
}

func BadRequest(w http.ResponseWriter, msg string, v ...any) *response {
	return &response{400, fmt.Errorf(msg, v...)}
}

func InternalServerError(w http.ResponseWriter, msg string, v ...any) *response {
	return &response{500, fmt.Errorf(msg, v...)}
}

func Unauthorized(w http.ResponseWriter, msg string, v ...any) *response {
	return &response{403, fmt.Errorf(msg, v...)}
}

func Status(status int) *response {
	return &response{status, nil}
}

type response struct {
	status int
	error  error
}

func (res *response) Errorf(w http.ResponseWriter, msg string, v ...any) *response {
	return &response{0, fmt.Errorf(msg, v...)}
}

func (res *response) JSON(v any) *response {
	return &response{0, nil}
}

func (res *response) Json(w http.ResponseWriter, v any) *response {
	return &response{0, nil}
}

func (res *response) Redirect(r *http.Request, to string, v ...any) *response {
	// return &response{0, fmt.Errorf(msg, v...)}
	return res
}

func (res *response) Respond(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(res.status)
}

func (res *response) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(res.status)
}
