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

func TestBatchMainSlot(t *testing.T) {
	t.Skip("TODO: fix this test")
	is := is.New(t)
	view := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		is.NoErr(err)
		inner, err := io.ReadAll(slot)
		is.NoErr(err)
		fmt.Fprintf(w, "<view>%s</view>", inner)
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
		fmt.Fprintf(w, "<layout>%s</layout>", inner)
	})
	handler := slots.Batch(view, frame1, frame2, layout)
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"hi"}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusOK)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "<layout><frame2><frame1><view></view></frame1></frame2></layout>")
}
