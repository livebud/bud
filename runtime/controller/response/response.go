package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"gitlab.com/mnm/bud/runtime/controller/request"
)

// Format returns different responses depending on the Accepts request header
// TODO: tighten up the types. Maybe HTML(w, r) & JSON(w, r) interfaces
type Format struct {
	HTML http.Handler
	JSON http.Handler
}

var _ http.Handler = (*Format)(nil)

func (f *Format) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	acceptable := request.Accepts(r)
	switch {
	case f.HTML != nil && acceptable.Accepts("text/html"):
		f.HTML.ServeHTTP(w, r)
	case f.JSON != nil && acceptable.Accepts("application/json"):
		f.JSON.ServeHTTP(w, r)
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	}
}

// Response struct
type Response struct {
	status  int
	headers map[string]string
}

// Status of a response
func Status(code int) *Response {
	return &Response{
		status:  code,
		headers: map[string]string{},
	}
}

// Set the header
func (res *Response) Set(key, value string) *Response {
	res.headers[key] = value
	return res
}

// Redirect to url
func (res *Response) Redirect(path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Attach all preset headers
		header := w.Header()
		for key, value := range res.headers {
			header.Set(key, value)
		}
		// Default status is 302 Found
		if res.status == 0 {
			res.status = http.StatusFound
		}
		// Redirect the response
		http.Redirect(w, r, path, res.status)
	})
}

// JSON response
func JSON(props interface{}) http.Handler {
	response := &Response{
		headers: map[string]string{},
	}
	return response.JSON(props)
}

// JSON responds with a JSON response.
func (res *Response) JSON(props interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Attach all preset headers
		header := w.Header()
		for key, value := range res.headers {
			header.Set(key, value)
		}
		// Override any existing content types
		header.Set("Content-Type", "application/json")
		// Marshal the JSON response
		result, err := json.Marshal(props)
		if err != nil {
			w.WriteHeader(500)
			// TODO: standardize this
			w.Write([]byte(fmt.Sprintf(`{"error":{"message":%q}}`, err.Error())))
			return
		}
		// Default status is 200 OK
		if res.status == 0 {
			res.status = 200
		}
		// Write the response
		w.WriteHeader(res.status)
		w.Write(result)
	})
}

// HTML response
func HTML(body string) http.Handler {
	response := &Response{
		headers: map[string]string{},
	}
	return response.HTML(body)
}

// HTML responds with an HTML response.
func (res *Response) HTML(body string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Attach all preset headers
		header := w.Header()
		for key, value := range res.headers {
			header.Set(key, value)
		}
		// Override any existing content types
		header.Set("Content-Type", "text/html")
		// Default status is 200 OK
		if res.status == 0 {
			res.status = 200
		}
		// Write the response
		w.WriteHeader(res.status)
		w.Write([]byte(wrapHTML(body)))
	})
}

// TODO: make hot reload configurable
func wrapHTML(body string) string {
	return `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8"/>
		</head>
		<body>
			` + body + `
			<script>
				const sse = new EventSource("http://0.0.0.0:35729")
				sse.addEventListener("message", () => { location.reload() })
			</script>
		</body>
		</html>`
}

func (res *Response) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Attach all preset headers
	header := w.Header()
	for key, value := range res.headers {
		header.Set(key, value)
	}
	if res.status == 0 {
		res.status = 200
	}
	w.WriteHeader(res.status)
}

// RedirectPath returns the response path.
func RedirectPath(r *http.Request, subpath string) string {
	switch r.Method {
	case "POST":
		return path.Join(r.URL.Path, subpath)
	case "DELETE":
		dir := path.Dir(r.URL.Path)
		if dir == "." {
			return "/"
		}
		return dir
	default:
		return r.URL.Path
	}
}
