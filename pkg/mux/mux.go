package mux

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"sort"
	"strings"

	"github.com/livebud/bud/pkg/middleware"
	"github.com/livebud/bud/pkg/mux/ast"
	"github.com/livebud/bud/pkg/mux/internal/radix"
	"github.com/livebud/bud/pkg/slots"
)

var (
	ErrDuplicate = radix.ErrDuplicate
	ErrNoMatch   = radix.ErrNoMatch
)

type Routes interface {
	// With(func(next http.Handler) http.Handler) Router
	Get(route string, fn http.HandlerFunc) error
	Post(route string, fn http.HandlerFunc) error
	Put(route string, fn http.HandlerFunc) error
	Patch(route string, fn http.HandlerFunc) error
	Delete(route string, fn http.HandlerFunc) error
	Layout(route string, fn http.HandlerFunc) error
	Error(route string, fn http.HandlerFunc) error
	Group(route string) Routes
}

type Match struct {
	Method  string
	Route   *ast.Route
	Path    string
	Slots   []*radix.Slot
	Handler http.Handler
	Layout  *LayoutHandler
	Error   *ErrorHandler
}

// WithBatch provides a custom batching function for mux. This is useful when a
// route has layout and frame handlers.
func WithBatch(batch func(...http.Handler) http.Handler) func(*Router) {
	return func(rt *Router) {
		rt.batch = batch
	}
}

func New(options ...func(*Router)) *Router {
	return &Router{
		base:    "",
		methods: map[string]*radix.Tree{},
		layouts: radix.New(),
		errors:  radix.New(),
		batch:   slots.Chain,
	}
}

type Router struct {
	base    string
	stack   []func(http.Handler) http.Handler
	methods map[string]*radix.Tree
	layouts *radix.Tree
	errors  *radix.Tree
	batch   func(...http.Handler) http.Handler
}

var _ http.Handler = (*Router)(nil)
var _ Routes = (*Router)(nil)

// Use appends middleware to the bottom of the stack
func (rt *Router) Use(middleware func(next http.Handler) http.Handler) {
	rt.stack = append(rt.stack, middleware)
}

// Get route
func (rt *Router) Get(route string, fn http.HandlerFunc) error {
	return rt.set(http.MethodGet, route, fn)
}

// Post route
func (rt *Router) Post(route string, fn http.HandlerFunc) error {
	return rt.set(http.MethodPost, route, fn)
}

// Put route
func (rt *Router) Put(route string, fn http.HandlerFunc) error {
	return rt.set(http.MethodPut, route, fn)
}

// Patch route
func (rt *Router) Patch(route string, fn http.HandlerFunc) error {
	return rt.set(http.MethodPatch, route, fn)
}

// Delete route
func (rt *Router) Delete(route string, fn http.HandlerFunc) error {
	return rt.set(http.MethodDelete, route, fn)
}

// Layout attaches a layout handler to a given route
func (rt *Router) Layout(route string, fn http.HandlerFunc) error {
	return rt.layouts.Insert(route, fn)
}

// Error attaches an error handler to a given route
func (rt *Router) Error(route string, fn http.HandlerFunc) error {
	return rt.errors.Insert(route, fn)
}

// Subrouter extends the router
type Subrouter interface {
	Routes(rt Routes)
}

// Add routes
func (rt *Router) Add(sr Subrouter) {
	sr.Routes(rt)
}

// Group routes within a route
func (rt *Router) Group(route string) Routes {
	return &Router{
		base:    strings.TrimSuffix(path.Join(rt.base, route), "/"),
		methods: rt.methods,
		layouts: rt.layouts,
		errors:  rt.errors,
	}
}

// ServeHTTP implements http.Handler
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := rt.Middleware(http.NotFoundHandler())
	handler.ServeHTTP(w, r)
}

// Middleware will return next on no match
func (rt *Router) Middleware(next http.Handler) http.Handler {
	return middleware.Compose(rt.stack...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Match the path
		match, err := rt.Match(r.Method, r.URL.Path)
		if err != nil {
			if errors.Is(err, radix.ErrNoMatch) {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		// Avoid calling other handlers if there's a path extension
		if path.Ext(r.URL.Path) != "" {
			match.Handler.ServeHTTP(w, r)
			return
		}
		// Chain the handlers together
		var handlers []http.Handler
		handlers = append(handlers, match.Handler)
		if match.Layout != nil {
			handlers = append(handlers, match.Layout.Handler)
		}
		// Batch the handlers together
		handler := rt.batch(handlers...)
		// Call the handler
		handler.ServeHTTP(w, r)
	}))
}

// Set a handler manually
func (rt *Router) Set(method string, route string, fn http.HandlerFunc) error {
	if !isMethod(method) {
		return fmt.Errorf("router: %q is not a valid HTTP method", method)
	}
	return rt.set(method, route, fn)
}

type Route struct {
	Method  string
	Route   string
	Handler http.Handler
	Layout  *LayoutHandler
	Error   *ErrorHandler
}

type LayoutHandler struct {
	Route   string
	Handler http.Handler
}

type ErrorHandler struct {
	Route   string
	Handler http.Handler
}

func (r *Route) String() string {
	return fmt.Sprintf("%s %s", r.Method, r.Route)
}

func (rt *Router) Find(method, route string) (*Route, error) {
	tree, ok := rt.methods[method]
	if !ok {
		return nil, fmt.Errorf("router: %w found for %s %s", ErrNoMatch, method, route)
	}
	node, err := tree.Find(route)
	if err != nil {
		return nil, err
	}
	r := &Route{
		Method:  method,
		Route:   node.Route.String(),
		Handler: node.Handler,
	}
	// Find the layout handler
	r.Layout, err = rt.findLayout(method, node.Route.String())
	if err != nil {
		return nil, err
	}
	// Find the error handler
	r.Error, err = rt.findError(method, node.Route.String())
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Routes lists all the routes
func (rt *Router) Routes() (routes []*Route) {
	for method, tree := range rt.methods {
		tree.Each(func(node *radix.Node) bool {
			if node.Route == nil {
				return true
			}
			routes = append(routes, &Route{
				Method:  method,
				Route:   node.Route.String(),
				Handler: node.Handler,
			})
			return true
		})
	}
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Method != routes[j].Method {
			return methodSort[routes[i].Method] < methodSort[routes[j].Method]
		}
		return routes[i].Route < routes[j].Route
	})
	return routes
}

// Match a route from a method and path
func (rt *Router) Match(method, path string) (*Match, error) {
	tree, ok := rt.methods[method]
	if !ok {
		return nil, fmt.Errorf("router: %w found for %s %s", ErrNoMatch, method, path)
	}
	m, err := tree.Match(path)
	if err != nil {
		return nil, err
	} else if m.Handler == nil {
		// Internal route without a handler
		return nil, fmt.Errorf("router: %w found for %s %s", ErrNoMatch, method, path)
	}
	match := &Match{
		Method:  method,
		Route:   m.Route,
		Path:    m.Path,
		Slots:   m.Slots,
		Handler: m.Handler,
	}
	// Find the layout handler
	match.Layout, err = rt.findLayout(method, m.Route.String())
	if err != nil {
		return nil, err
	}
	// Find the error handler
	match.Error, err = rt.findError(method, m.Route.String())
	if err != nil {
		return nil, err
	}
	return match, nil
}

// ServeFS serves a filesystem
func (rt *Router) ServeFS(path string, fsys fs.FS) error {
	return rt.set(http.MethodGet, path, http.FileServer(http.FS(fsys)))
}

// FindLayout finds the layout for a given path
func (rt *Router) findLayout(method, route string) (*LayoutHandler, error) {
	if method != http.MethodGet {
		return nil, nil
	}
	layout, err := rt.layouts.FindByPrefix(route)
	if err != nil {
		if !errors.Is(err, ErrNoMatch) {
			return nil, err
		}
		return nil, nil
	}
	return &LayoutHandler{
		Route:   layout.Route.String(),
		Handler: layout.Handler,
	}, nil
}

// Find the error handler
func (rt *Router) findError(method, route string) (*ErrorHandler, error) {
	if method != http.MethodGet {
		return nil, nil
	}
	errorHandler, err := rt.errors.FindByPrefix(route)
	if err != nil {
		if !errors.Is(err, ErrNoMatch) {
			return nil, err
		}
		return nil, nil
	}
	return &ErrorHandler{
		Route:   errorHandler.Route.String(),
		Handler: errorHandler.Handler,
	}, nil
}

var methodSort = map[string]int{
	http.MethodGet:    0,
	http.MethodPost:   1,
	http.MethodPut:    2,
	http.MethodPatch:  3,
	http.MethodDelete: 4,
}

// Set the route
func (rt *Router) set(method, route string, handler http.Handler) error {
	return rt.insert(method, path.Join(rt.base, route), handler)
}

// Insert the route into the method's radix tree
func (rt *Router) insert(method, route string, handler http.Handler) error {
	if _, ok := rt.methods[method]; !ok {
		rt.methods[method] = radix.New()
	}
	return rt.methods[method].Insert(route, handler)
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
