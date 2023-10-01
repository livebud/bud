package controller

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"reflect"

	"github.com/livebud/bud/mux"
	"github.com/livebud/bud/pkg/controller/internal/request"
	"github.com/livebud/bud/pkg/view"
	"github.com/livebud/bud/pkg/view/slots"
	"golang.org/x/sync/errgroup"
)

type viewWriter struct {
	router       *mux.Router
	page         *view.Page
	handler      reflect.Value
	viewHandlers map[string]http.Handler
}

var _ writer = (*viewWriter)(nil)

func (v *viewWriter) WriteEmpty(w http.ResponseWriter, r *http.Request) {
	v.WriteOutput(w, r, nil)
}

type props struct {
	Error string `json:"error,omitempty"`
}

type responseWriter struct {
	w    http.ResponseWriter
	slot view.Slot
}

func (r *responseWriter) Header() http.Header {
	// TODO: only write certain headers (e.g. exclude Content-Type)
	return r.w.Header()
}

func (r *responseWriter) Write(b []byte) (int, error) {
	return r.slot.Write(b)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	// r.w.WriteHeader(statusCode)
}

func (v *viewWriter) WriteOutput(w http.ResponseWriter, r *http.Request, out any) {
	w.Header().Set("Content-Type", "text/html")
	b := new(bytes.Buffer)
	slot := slots.New(b)
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		v.WriteError(w, r, err)
		return
	}
	eg, ctx := errgroup.WithContext(ctx)
	if v.page.Layout != nil {
		innerSlot := slot
		eg.Go(func() (err error) {
			defer innerSlot.Close(&err)
			if handler, ok := v.viewHandlers[v.page.Layout.Key()]; ok {
				handler.ServeHTTP(&responseWriter{w, innerSlot}, request.Clone(ctx, r, body))
				return nil
			}
			err = v.page.Layout.Render(ctx, innerSlot, map[string]interface{}{})
			return err
		})
		slot = slot.New()
	}
	// Render the frames
	for _, frame := range v.page.Frames {
		frame := frame
		innerSlot := slot
		eg.Go(func() (err error) {
			// TODO: lookup the frame handler from the router
			defer innerSlot.Close(&err)
			if handler, ok := v.viewHandlers[frame.Key()]; ok {
				handler.ServeHTTP(&responseWriter{w, innerSlot}, request.Clone(ctx, r, body))
				return nil
			}
			if err := frame.Render(ctx, innerSlot, map[string]interface{}{}); err != nil {
				return err
			}
			return nil
		})
		slot = slot.New()
	}
	// Render the page
	eg.Go(func() (err error) {
		defer slot.Close(&err)
		if out == nil {
			out = map[string]interface{}{}
		}
		err = v.page.View.Render(ctx, slot, out)
		return err
	})
	// Wait for the results of all the renders
	if err := eg.Wait(); err != nil {
		v.WriteError(w, request.Clone(ctx, r, body), err)
		return
	}
	if _, err := b.WriteTo(w); err != nil {
		v.WriteError(w, request.Clone(ctx, r, body), err)
		return
	}
}

func (v *viewWriter) WriteError(w http.ResponseWriter, r *http.Request, originalError error) {
	w.Header().Set("Content-Type", "text/html")
	status := getStatus(http.StatusInternalServerError, originalError)
	w.WriteHeader(status)
	b := new(bytes.Buffer)
	slot := slots.New(b)
	ctx := r.Context()
	eg := new(errgroup.Group)
	if v.page.Layout != nil {
		innerSlot := slot
		eg.Go(func() (err error) {
			defer innerSlot.Close(&err)
			// TODO: lookup the layout handler from the router
			return v.page.Layout.Render(ctx, innerSlot, map[string]interface{}{})
		})
		slot = slot.New()
	}
	// Render the frames
	for _, frame := range v.page.Frames {
		frame := frame
		innerSlot := slot
		eg.Go(func() (err error) {
			// TODO: lookup the frame handler from the router
			defer innerSlot.Close(&err)
			return frame.Render(ctx, innerSlot, map[string]interface{}{})
		})
		slot = slot.New()
	}
	// Render the error page
	eg.Go(func() (err error) {
		defer slot.Close(&err)
		if v.page.Error == nil {
			w.Write([]byte(originalError.Error()))
			return
		}
		props := &props{
			Error: originalError.Error(),
		}
		return v.page.Error.Render(ctx, slot, props)
	})
	// Wait for the results of all the renders
	if err2 := eg.Wait(); err2 != nil {
		v.renderOnlyErrorPage(ctx, w, originalError)
		return
	}
	if _, err := b.WriteTo(w); err != nil {
		v.renderOnlyErrorPage(ctx, w, originalError)
		return
	}
}

func (v *viewWriter) renderOnlyErrorPage(ctx context.Context, w io.Writer, err error) {
	props := &props{
		Error: err.Error(),
	}
	slot := slots.New(w)
	defer slot.Close(&err)
	if v.page.Error == nil {
		w.Write([]byte(err.Error()))
		return
	}
	if err := v.page.Error.Render(ctx, slot, props); err != nil {
		slot.Write([]byte(err.Error()))
	}
}
