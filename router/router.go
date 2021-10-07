package router

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

// New router
func New() *Router {
	return &Router{
		paths: map[string]string{},
	}
}

// Route struct
type Route struct {
	router *Router
	key    string
}

// Key returns the render key
func (r *Route) Key() string {
	return r.key
}

// Keyed is an optional interface to support
// looking up routes by this key
type Keyed interface {
	Key() string
}

// Middleware type alias
type Middleware = func(http.Handler) http.Handler

// Router struct
type Router struct {
	stack []Middleware
	paths map[string]string

	once    sync.Once
	handler http.Handler
}

var _ http.Handler = (*Router)(nil)

// Use middleware
func (rt *Router) Use(middleware Middleware) *Router {
	rt.stack = append(rt.stack, middleware)
	return rt
}

// Add a path
func (rt *Router) Add(method, path string, handler http.Handler) {
	rt.stack = append(rt.stack, rt.wrap(method, path, handler))
}

// Get route
func (rt *Router) Get(path string, handler http.Handler) {
	rt.stack = append(rt.stack, rt.wrap(http.MethodGet, path, handler))
}

// Post route
func (rt *Router) Post(path string, handler http.Handler) {
	rt.stack = append(rt.stack, rt.wrap(http.MethodPost, path, handler))
}

// Put route
func (rt *Router) Put(path string, handler http.Handler) {
	rt.stack = append(rt.stack, rt.wrap(http.MethodPut, path, handler))
}

// Patch route
func (rt *Router) Patch(path string, handler http.Handler) {
	rt.stack = append(rt.stack, rt.wrap(http.MethodPatch, path, handler))
}

// Delete route
func (rt *Router) Delete(path string, handler http.Handler) {
	rt.stack = append(rt.stack, rt.wrap(http.MethodDelete, path, handler))
}

// ServeFS serves a filesystem
func (rt *Router) ServeFS(fs http.FileSystem) {
	rt.stack = append(rt.stack, rt.serveFS(fs))
}

// Compose middleware into one
func compose(stack ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		if len(stack) == 0 {
			return h
		}
		for i := len(stack) - 1; i >= 0; i-- {
			h = stack[i](h)
		}
		return h
	}
}

func (rt *Router) compose() {
	wrap := compose(rt.stack...)
	bottom := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	rt.handler = wrap(bottom)
}

// ServeHTTP function
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.once.Do(rt.compose)
	rt.handler.ServeHTTP(w, r)
}

// Convenience function for mapping a URL path to map of parameters.
// TODO: decide if I should keep this method
func Match(path, route string) map[string]string {
	regexp, err := Parse(route)
	if err != nil {
		return map[string]string{}
	}
	matcher := &matcher{regexp}
	return matcher.match(path)
}

type matcher struct {
	re *regexp.Regexp
}

func (m *matcher) match(path string) (params map[string]string) {
	match := m.re.FindStringSubmatch(path)
	if match == nil {
		return nil
	}
	params = map[string]string{}
	for i, name := range m.re.SubexpNames() {
		if i != 0 && name != "" {
			params[name] = match[i]
		}
	}
	return params
}

func (rt *Router) wrap(method, routePath string, handler http.Handler) Middleware {
	// Parse the path
	regexp, err := Parse(routePath)
	if err != nil {
		panic(fmt.Errorf("router is unable to parse %q route", routePath))
	}
	// Create the route
	route := &Route{router: rt}
	if keyed, ok := handler.(Keyed); ok {
		key := keyed.Key()
		route.key = key
	}
	matcher := &matcher{regexp}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				next.ServeHTTP(w, r)
				return
			}
			urlPath := r.URL.Path
			// Always redirect to non-trailing slash
			// TODO: figure out how to persist flashes across this redirect
			if urlPath != "/" && strings.HasSuffix(urlPath, "/") {
				r.URL.Path = strings.TrimSuffix(urlPath, "/")
				http.Redirect(w, r, r.URL.String(), http.StatusPermanentRedirect)
				return
			}
			// Match params with our URL
			params := matcher.match(urlPath)
			if params == nil {
				next.ServeHTTP(w, r)
				return
			}
			// Set params as query parameters.
			// Overrides conflicting query parameters.
			values := r.URL.Query()
			for key, value := range params {
				values.Set(key, value)
			}
			r.URL.RawQuery = values.Encode()
			// Add the router to the context
			// r = r.WithContext(context.WithValue(r.Context(), contextkey.Router, route))
			// Handle the request
			handler.ServeHTTP(w, r)
		})
	}
}

// serveFS serves the filesystem
func (rt *Router) serveFS(fs http.FileSystem) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlPath := r.URL.Path
			if r.Method != http.MethodGet || path.Ext(urlPath) == "" {
				next.ServeHTTP(w, r)
				return
			}
			file, err := fs.Open(urlPath)
			if err != nil {
				if os.IsNotExist(err) {
					next.ServeHTTP(w, r)
					return
				}
				fmt.Println("Static open error", err)
				return
			}
			stat, err := file.Stat()
			if err != nil {
				fmt.Println("Static stat error", err)
				return
			}
			if stat.IsDir() {
				next.ServeHTTP(w, r)
				return
			}
			http.ServeContent(w, r, urlPath, stat.ModTime(), file)
		})
	}
}

// // From pulls the route from context
// func From(ctx context.Context) *Route {
// 	v := ctx.Value(contextkey.Router)
// 	if v == nil {
// 		return nil
// 	}
// 	route, ok := v.(*Route)
// 	if !ok {
// 		return nil
// 	}
// 	return route
// }

// FileRouter alias
type FileRouter = Router
