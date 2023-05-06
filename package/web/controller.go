package web

import (
	"fmt"
	"net/http"

	"github.com/livebud/bud/framework/controller/controllerrt/request"
)

type Request[Body any] struct {
	Body Body
}

type Response[Body any] interface {
}

type response[Body any] struct {
}

func Handler[Req any, Res any](fn func(req *Request[Req], res Response[Res]) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body Req
		if err := request.Unmarshal(r, &body); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		req := &Request[Req]{
			Body: body,
		}
		res := &response[Res]{}
		fn(req, res)
	})
}
