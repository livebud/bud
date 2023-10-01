package controller

import (
	"fmt"
	"net/http"
	"path"
	"reflect"

	"github.com/livebud/bud/mux"
	"github.com/livebud/bud/pkg/controller/internal/viewkey"
	"github.com/livebud/bud/pkg/view"
	"github.com/matthewmueller/text"
)

func New(vf view.Finder, router *mux.Router) *Router {
	return &Router{vf, router, map[string]http.Handler{}}
}

type Router struct {
	vf           view.Finder
	router       *mux.Router
	viewHandlers map[string]http.Handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) Register(route string, controller any) error {
	return r.register(route, reflect.ValueOf(controller))
}

func (r *Router) Get(route string, handler any) error {
	return r.registerFunction(http.MethodGet, route, reflect.ValueOf(handler))
}

func isStructOrPtrStruct(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct
}

var reservedMethods = map[string]bool{
	"Layout": true,
	"Frame":  true,
	"Error":  true,
}

type layoutFrameErrorWriter struct {
	key  string
	view view.View
}

func (v *layoutFrameErrorWriter) WriteEmpty(w http.ResponseWriter, r *http.Request) {
	v.WriteOutput(w, r, nil)
}

func (v *layoutFrameErrorWriter) WriteOutput(w http.ResponseWriter, r *http.Request, out any) {
	rs, ok := w.(*responseWriter)
	if !ok {
		v.WriteError(w, r, fmt.Errorf("expected *responseWriter, got %T", out))
		return
	}
	if err := v.view.Render(r.Context(), rs.slot, out); err != nil {
		v.WriteError(w, r, err)
		return
	}
}

func (v *layoutFrameErrorWriter) WriteError(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Println("writing error", v.key, err)
}

func (r *Router) register(route string, recv reflect.Value) error {
	t := recv.Type()
	if !isStructOrPtrStruct(t) {
		return fmt.Errorf("rpc: expected a struct, got %v", t)
	}
	numMethods := t.NumMethod()
	for i := 0; i < numMethods; i++ {
		m := t.Method(i)
		if reservedMethods[m.Name] {
			key, err := viewkey.Infer(routePath(route, m.Name))
			if err != nil {
				return err
			}
			view, err := r.vf.FindView(key)
			if err != nil {
				return err
			}
			viewHandler, err := wrapValue(&defaultReader{}, &layoutFrameErrorWriter{key, view}, m.Func, recv)
			if err != nil {
				return err
			}
			r.viewHandlers[key] = viewHandler
			continue
		}
		if err := r.registerMethod(httpMethod(m.Name), routePath(route, m.Name), recv, m.Func); err != nil {
			return err
		}
	}
	return nil
}

func (r *Router) htmlWriter(httpMethod, route string, fn reflect.Value) (writer, error) {
	switch httpMethod {
	case http.MethodGet:
		key, err := viewkey.Infer(route)
		if err != nil {
			return nil, err
		}
		page, err := r.vf.FindPage(key)
		if err != nil {
			return nil, err
		}
		return &viewWriter{r.router, page, fn, r.viewHandlers}, nil
	case http.MethodPost, http.MethodPatch, http.MethodDelete:
		return &formWriter{fn}, nil
	default:
		return nil, fmt.Errorf("rpc: unhandled writer for HTTP %s method", httpMethod)
	}
}

func maybeHandler(fn reflect.Value, recv ...reflect.Value) (http.Handler, bool) {
	offset := len(recv)
	fnType := fn.Type()
	if fnType.NumOut() != 0 || fnType.NumIn()-offset != 2 {
		return nil, false
	} else if !isResponseWriter(fnType.In(0 + offset)) {
		return nil, false
	} else if !isRequest(fnType.In(1 + offset)) {
		return nil, false
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn.Call(append(recv, reflect.ValueOf(w), reflect.ValueOf(r)))
	}), true
}

func (r *Router) registerMethod(httpMethod, route string, recv, methodFn reflect.Value) error {
	if handler, ok := maybeHandler(methodFn, recv); ok {
		return r.router.Set(httpMethod, route, handler)
	}
	htmlWriter, err := r.htmlWriter(httpMethod, route, methodFn)
	if err != nil {
		return err
	}
	jsonWriter := newJsonWriter()
	accepts := &accept{
		HTML: htmlWriter,
		JSON: jsonWriter,
	}
	handler, err := wrapValue(&defaultReader{}, accepts, methodFn, recv)
	if err != nil {
		return err
	}
	return r.router.Set(httpMethod, route, handler)
}

func (r *Router) registerFunction(httpMethod, route string, fn reflect.Value) error {
	htmlWriter, err := r.htmlWriter(httpMethod, route, fn)
	if err != nil {
		return err
	}
	jsonWriter := newJsonWriter()
	accepts := &accept{
		HTML: htmlWriter,
		JSON: jsonWriter,
	}
	handler, err := wrapValue(&defaultReader{}, accepts, fn)
	if err != nil {
		return err
	}
	return r.router.Set(httpMethod, route, handler)
}

const (
	methodIndex  = "Index"
	methodShow   = "Show"
	methodNew    = "New"
	methodEdit   = "Edit"
	methodCreate = "Create"
	methodUpdate = "Update"
	methodDelete = "Delete"
)

// Route to the action
func routePath(route, methodName string) string {
	switch methodName {
	case methodShow, methodUpdate, methodDelete:
		return path.Join(route, "{id}")
	case methodNew:
		return path.Join(route, "new")
	case methodEdit:
		return path.Join(route, "{id}", "edit")
	case methodIndex, methodCreate:
		return route
	default:
		return path.Join(route, text.Lower(text.Snake(methodName)))
	}
}

// Method is the HTTP method for this controller
func httpMethod(actionName string) string {
	switch actionName {
	case "Create":
		return http.MethodPost
	case "Update":
		return http.MethodPatch
	case "Delete":
		return http.MethodDelete
	default:
		return http.MethodGet
	}
}
