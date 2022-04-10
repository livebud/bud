package router

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"gitlab.com/mnm/bud/package/router/radix"
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
		// Redirect for trailing slashes or paths with uppercase letters
		urlPath := r.URL.Path
		redirect := false
		// Strip any trailing slash (e.g. /users/ => /users)
		if hasTrailingSlash(urlPath) {
			urlPath = strings.TrimRight(urlPath, "/")
			redirect = true
		}
		// Ensure that all paths are case-insensitive (e.g. /USERS => /users)
		if hasUpper(urlPath) {
			urlPath = strings.ToLower(urlPath)
			redirect = true
		}
		// Redirect all at once, instead of for each rule
		if redirect {
			http.Redirect(w, r, strings.ToLower(urlPath), http.StatusPermanentRedirect)
			return
		}
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

func hasTrailingSlash(path string) bool {
	return path != "/" && strings.HasSuffix(path, "/")
}

func hasUpper(path string) bool {
	for _, r := range path {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
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
