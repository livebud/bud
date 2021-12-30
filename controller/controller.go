package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/mnm/bud/view"
)

func New(w http.ResponseWriter, r *http.Request, v view.Renderer) *Context {
	return &Context{w, r, v}
}

type Context struct {
	w http.ResponseWriter
	r *http.Request
	v view.Renderer
}

func (c *Context) Unmarshal(v interface{}) error {
	return Unmarshal(c.r, v)
}

func (c *Context) Status(code int) *Response {
	return &Response{
		writer: c.w,
		header: c.w.Header(),
		view:   c.v,
		status: code,
	}
}

// Response struct
type Response struct {
	writer http.ResponseWriter
	view   view.Renderer
	header http.Header
	status int
}

// Set the header
func (res *Response) Set(key, value string) *Response {
	res.header.Set(key, value)
	return res
}

// Redirect to url
func (res *Response) Redirect(path string) {
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	// Attach all preset headers
	// 	header := w.Header()
	// 	for key, value := range res.headers {
	// 		header.Set(key, value)
	// 	}
	// 	// Default status is 302 Found
	// 	if res.status == 0 {
	// 		res.status = http.StatusFound
	// 	}
	// 	// Redirect the response
	// 	http.Redirect(w, r, path, res.status)
	// })
}

func (res *Response) RenderError(err error) {
	res.JSON(map[string]string{
		"error": err.Error(),
	})
}

func (res *Response) Render(path string, props interface{}) {
	fmt.Println("rendering", path, props)
	response, err := res.view.Render(path, props)
	if err != nil {
		res.RenderError(err)
		return
	}
	for key, val := range response.Headers {
		res.header.Set(key, val)
	}
	// Write the response
	res.writer.WriteHeader(response.Status)
	res.writer.Write([]byte(response.Body))
}

func (res *Response) JSON(props interface{}) {
	// Override any existing content types
	res.header.Set("Content-Type", "application/json")
	// Marshal the JSON response
	result, err := json.Marshal(props)
	if err != nil {
		res.writer.WriteHeader(500)
		// TODO: standardize this
		res.writer.Write([]byte(fmt.Sprintf(`{"error":{"message":%q}}`, err.Error())))
		return
	}
	// Write the response
	res.writer.WriteHeader(res.status)
	res.writer.Write(result)
}

// func (res *Response) Text(text []byte) {
// 	// Override any existing content types
// 	res.header.Set("Content-Type", "text/plain")
// 	// Write the response
// 	res.writer.WriteHeader(res.status)
// 	res.writer.Write(text)
// }
