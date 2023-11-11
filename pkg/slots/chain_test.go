package slots_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/livebud/bud/pkg/slots"
	"github.com/matryer/is"
)

func TestChainOneHandler(t *testing.T) {
	is := is.New(t)
	h1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		is.NoErr(err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(body)
	})
	handler := slots.Chain(h1)
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"hi"}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusCreated)
	is.Equal(res.Header.Get("Content-Type"), "application/json")
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), `{"name":"hi"}`)
}

func TestChainMainSlot(t *testing.T) {
	is := is.New(t)
	view := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		is.NoErr(err)
		inner, err := io.ReadAll(slot)
		is.NoErr(err)
		w.Header().Add("X-View", "V")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "<view>%s</view>", inner)
	})
	frame1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		is.NoErr(err)
		inner, err := io.ReadAll(slot)
		is.NoErr(err)
		w.Header().Add("X-Frame", "F1")
		fmt.Fprintf(w, "<frame1>%s</frame1>", inner)
	})
	frame2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		is.NoErr(err)
		inner, err := io.ReadAll(slot)
		is.NoErr(err)
		w.Header().Add("X-Frame", "F2")
		fmt.Fprintf(w, "<frame2>%s</frame2>", inner)
	})
	layout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		is.NoErr(err)
		inner, err := io.ReadAll(slot)
		is.NoErr(err)
		w.Header().Add("X-Layout", "L")
		fmt.Fprintf(w, "<layout>%s</layout>", inner)
	})
	handler := slots.Chain(view, frame1, frame2, layout)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusCreated)
	is.Equal(res.Header.Get("X-View"), "V")
	is.Equal(res.Header.Get("X-Frame"), "F2")
	is.Equal(res.Header.Get("X-Layout"), "L")
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "<layout><frame2><frame1><view></view></frame1></frame2></layout>")
}

func TestChainOtherSlots(t *testing.T) {
	is := is.New(t)
	view := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		is.NoErr(err)
		inner, err := io.ReadAll(slot)
		is.NoErr(err)
		fmt.Fprintf(w, "<view>%s</view>", inner)
		slot.Slot("script").Write([]byte(`<script src='module.js'></script>`))
	})
	frame1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		is.NoErr(err)
		inner, err := io.ReadAll(slot)
		is.NoErr(err)
		fmt.Fprintf(w, "<frame1>%s</frame1>", inner)
	})
	frame2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		is.NoErr(err)
		inner, err := io.ReadAll(slot)
		is.NoErr(err)
		fmt.Fprintf(w, "<frame2>%s</frame2>", inner)
	})
	layout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		is.NoErr(err)
		inner, err := io.ReadAll(slot)
		is.NoErr(err)
		script, err := io.ReadAll(slot.Slot("script"))
		is.NoErr(err)
		fmt.Fprintf(w, "<layout>%s%s</layout>", script, inner)
	})
	handler := slots.Chain(view, frame1, frame2, layout)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusOK)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "<layout><script src='module.js'></script><frame2><frame1><view></view></frame1></frame2></layout>")
}

func TestChainErrorPriority(t *testing.T) {
	is := is.New(t)
	view := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		is.NoErr(err)
		inner, err := io.ReadAll(slot)
		is.NoErr(err)
		w.Header().Add("X-View", "V")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<view>%s</view>", inner)
	})
	frame1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		is.NoErr(err)
		inner, err := io.ReadAll(slot)
		is.NoErr(err)
		w.Header().Add("X-Frame", "F1")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "<frame1>%s</frame1>", inner)
	})
	frame2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		is.NoErr(err)
		inner, err := io.ReadAll(slot)
		is.NoErr(err)
		w.Header().Add("X-Frame", "F2")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "<frame2>%s</frame2>", inner)
	})
	layout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		is.NoErr(err)
		inner, err := io.ReadAll(slot)
		is.NoErr(err)
		w.Header().Add("X-Layout", "L")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "<layout>%s</layout>", inner)
	})
	handler := slots.Chain(view, frame1, frame2, layout)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusInternalServerError)
	is.Equal(res.Header.Get("X-View"), "V")
	is.Equal(res.Header.Get("X-Frame"), "F2")
	is.Equal(res.Header.Get("X-Layout"), "L")
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "<layout><frame2><frame1><view></view></frame1></frame2></layout>")
}
