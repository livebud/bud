package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/livebud/bud/package/router/radix"
)

// New router
func New() *Router {
	return &Router{
		methods: map[string]radix.Tree{},
	}
}

// Router struct
type Router struct {
	methods map[string]radix.Tree
}

var _ http.Handler = (*Router)(nil)

// Add a handler to a route
func (rt *Router) Add(method, route string, handler http.Handler) error {
	if !isMethod(method) {
		return fmt.Errorf("router: %q is not a valid HTTP method", method)
	}
	return rt.add(method, route, handler)
}

func (rt *Router) add(method, route string, handler http.Handler) error {
	if route == "/" {
		return rt.insert(method, route, handler)
	}
	// Trim any trailing slash and lowercase the route
	route = strings.TrimRight(strings.ToLower(route), "/")
	return rt.insert(method, route, handler)
}

// Insert the route into the method's radix tree
func (rt *Router) insert(method, route string, handler http.Handler) error {
	if _, ok := rt.methods[method]; !ok {
		rt.methods[method] = radix.New()
	}
	return rt.methods[method].Insert(route, handler)
}

// Get route
func (rt *Router) Get(route string, handler http.Handler) error {
	return rt.add(http.MethodGet, route, handler)
}

// Post route
func (rt *Router) Post(route string, handler http.Handler) error {
	return rt.add(http.MethodPost, route, handler)
}

// Put route
func (rt *Router) Put(route string, handler http.Handler) error {
	return rt.add(http.MethodPut, route, handler)
}

// Patch route
func (rt *Router) Patch(route string, handler http.Handler) error {
	return rt.add(http.MethodPatch, route, handler)
}

// Delete route
func (rt *Router) Delete(route string, handler http.Handler) error {
	return rt.add(http.MethodDelete, route, handler)
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := rt.Middleware(http.NotFoundHandler())
	handler.ServeHTTP(w, r)
}

// Middleware implements the router middleware
func (rt *Router) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tree, ok := rt.methods[r.Method]
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		// Strip any trailing slash (e.g. /users/ => /users)
		urlPath := trimTrailingSlash(r.URL.Path)
		// Match the path
		match, ok := tree.Match(urlPath)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		// Add the slots
		if len(match.Slots) > 0 {
			query := r.URL.Query()
			for _, slot := range match.Slots {
				query.Set(slot.Key, slot.Value)
			}
			r.URL.RawQuery = query.Encode()
		}
		// Call the handler
		match.Handler.ServeHTTP(w, r)
	})
}

func trimTrailingSlash(path string) string {
	if path == "/" {
		return path
	}
	return strings.TrimRight(path, "/")
}

// isMethod returns true if method is a valid HTTP method
func isMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPost,
		http.MethodPut, http.MethodPatch, http.MethodDelete,
		http.MethodConnect, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}
