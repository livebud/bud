package slots_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/livebud/bud/pkg/slots"
	"github.com/matryer/is"
)

func TestSyncSlots(t *testing.T) {
	is := is.New(t)
	s1 := slots.New()
	inner, err := io.ReadAll(s1)
	is.NoErr(err)
	is.Equal(string(inner), "")
	s1.Write([]byte("from s1"))
	inner, err = io.ReadAll(s1)
	is.NoErr(err)
	is.Equal(string(inner), "")
	s1.Close()
	s2 := s1.Next()
	inner, err = io.ReadAll(s2)
	is.NoErr(err)
	is.Equal(string(inner), "from s1")
	s2.Write([]byte("from s2"))
	inner, err = io.ReadAll(s2)
	is.NoErr(err)
	is.Equal(string(inner), "")
	s2.Close()
	s3 := s2.Next()
	inner, err = io.ReadAll(s3)
	is.NoErr(err)
	is.Equal(string(inner), "from s2")
	inner, err = io.ReadAll(s3)
	is.NoErr(err)
	is.Equal(string(inner), "")
	s3.Close()
}

func TestReadAfterClose(t *testing.T) {
	is := is.New(t)
	s1 := slots.New()
	s1.Write([]byte("from s1"))
	s2 := s1.Next()
	s1.Close()
	s2.Close()
	inner, err := io.ReadAll(s2)
	is.NoErr(err)
	is.Equal(string(inner), "from s1")
}

func TestAsyncSlots(t *testing.T) {
	is := is.New(t)
	s1 := slots.New()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer s1.Close()
		inner, err := io.ReadAll(s1)
		is.NoErr(err)
		is.Equal(string(inner), "")
		s1.Write([]byte("from s1"))
		inner, err = io.ReadAll(s1)
		is.NoErr(err)
		is.Equal(string(inner), "")
	}()
	s2 := s1.Next()
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer s2.Close()
		inner, err := io.ReadAll(s2)
		is.NoErr(err)
		is.Equal(string(inner), "from s1")
		s2.Write([]byte("from s2"))
		inner, err = io.ReadAll(s2)
		is.NoErr(err)
		is.Equal(string(inner), "")
	}()
	s3 := s2.Next()
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer s3.Close()
		inner, err := io.ReadAll(s3)
		is.NoErr(err)
		is.Equal(string(inner), "from s2")
		inner, err = io.ReadAll(s3)
		is.NoErr(err)
		is.Equal(string(inner), "")
		s3.Write([]byte("from s3"))
	}()
	s4 := s3.Next()
	inner, err := io.ReadAll(s4)
	is.NoErr(err)
	is.Equal(string(inner), "from s3")
	wg.Wait()
}

func TestSyncNestedSlots(t *testing.T) {
	is := is.New(t)
	s1 := slots.New()
	inner, err := io.ReadAll(s1)
	is.NoErr(err)
	is.Equal(string(inner), "")
	s1.Write([]byte("from s1"))
	inner, err = io.ReadAll(s1)
	is.NoErr(err)
	is.Equal(string(inner), "")
	s1a := s1.Slot("s1a")
	s1a.Write([]byte("from s1a"))
	s1.Close()
	s2 := s1.Next()
	s2a := s1.Slot("s1a")
	s2a.Write([]byte(" from s2a"))
	inner, err = io.ReadAll(s2)
	is.NoErr(err)
	is.Equal(string(inner), "from s1")
	s2.Write([]byte("from s2"))
	inner, err = io.ReadAll(s2)
	is.NoErr(err)
	is.Equal(string(inner), "")
	s2.Close()
	s3 := s2.Next()
	inner, err = io.ReadAll(s3)
	is.NoErr(err)
	is.Equal(string(inner), "from s2")
	inner, err = io.ReadAll(s3)
	is.NoErr(err)
	is.Equal(string(inner), "")
	s3a := s3.Slot("s1a")
	inner, err = io.ReadAll(s3a)
	is.NoErr(err)
	is.Equal(string(inner), "from s1a from s2a")
	s3.Close()
}

func TestAsyncNestedSlots(t *testing.T) {
	is := is.New(t)
	s1 := slots.New()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer s1.Close()
		inner, err := io.ReadAll(s1)
		is.NoErr(err)
		is.Equal(string(inner), "")
		s1.Write([]byte("from s1"))
		inner, err = io.ReadAll(s1)
		is.NoErr(err)
		is.Equal(string(inner), "")
		s1a := s1.Slot("s1a")
		s1a.Write([]byte("from s1a"))
	}()
	s2 := s1.Next()
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer s2.Close()
		inner, err := io.ReadAll(s2)
		is.NoErr(err)
		is.Equal(string(inner), "from s1")
		s2.Write([]byte("from s2"))
		inner, err = io.ReadAll(s2)
		is.NoErr(err)
		is.Equal(string(inner), "")
		s2a := s2.Slot("s1a")
		s2a.Write([]byte(" from s2a"))
	}()
	s3 := s2.Next()
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer s3.Close()
		inner, err := io.ReadAll(s3)
		is.NoErr(err)
		is.Equal(string(inner), "from s2")
		inner, err = io.ReadAll(s3)
		is.NoErr(err)
		is.Equal(string(inner), "")
		s3.Write([]byte("from s3"))
	}()
	s4 := s3.Next()
	inner, err := io.ReadAll(s4)
	is.NoErr(err)
	is.Equal(string(inner), "from s3")
	s4a := s4.Slot("s1a")
	inner, err = io.ReadAll(s4a)
	is.NoErr(err)
	is.Equal(string(inner), "from s1a from s2a")
	wg.Wait()
	s4.Close()
}

type responseWriter struct {
	http.ResponseWriter
	slots io.Writer
}

func (w *responseWriter) Write(p []byte) (n int, err error) {
	return w.slots.Write(p)
}

func TestChain(t *testing.T) {
	is := is.New(t)
	handlers := []func(w http.ResponseWriter, r *http.Request){
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			w.Write([]byte("<view>"))
			w.Write(inner)
			w.Write([]byte("</view>"))
		},
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			w.Write([]byte("<frame1>"))
			w.Write(inner)
			w.Write([]byte("</frame1>"))
		},
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			w.Write([]byte("<frame2>"))
			w.Write(inner)
			w.Write([]byte("</frame2>"))
		},
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			w.Write([]byte("<layout>"))
			w.Write(inner)
			w.Write([]byte("</layout>"))
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := slots.New()
		for i, handler := range handlers {
			r = r.Clone(slots.ToContext(r.Context(), s))
			if i == len(handlers)-1 {
				handler(w, r)
				continue
			}
			w := &responseWriter{w, s}
			handler(w, r)
			s.Close()
			s = s.Next()
		}
	})
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusOK)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "<layout><frame2><frame1><view></view></frame1></frame2></layout>")
}

func TestChainNestedSlots(t *testing.T) {
	is := is.New(t)
	handlers := []func(w http.ResponseWriter, r *http.Request){
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			headSlot := slots.Slot("head")
			headSlot.Write([]byte("<title>some title</title>"))
			w.Write([]byte("<view>"))
			w.Write(inner)
			w.Write([]byte("</view>"))
		},
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			headSlot := slots.Slot("head")
			headSlot.Write([]byte("<meta name='description' content='some description'/>"))
			w.Write([]byte("<frame1>"))
			w.Write(inner)
			w.Write([]byte("</frame1>"))
		},
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			w.Write([]byte("<frame2>"))
			w.Write(inner)
			w.Write([]byte("</frame2>"))
		},
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			headSlot := slots.Slot("head")
			head, err := io.ReadAll(headSlot)
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			w.Write([]byte("<layout>"))
			w.Write([]byte("<head>"))
			w.Write(head)
			w.Write([]byte("</head>"))
			w.Write(inner)
			w.Write([]byte("</layout>"))
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := slots.New()
		for i, handler := range handlers {
			r = r.Clone(slots.ToContext(r.Context(), s))
			if i == len(handlers)-1 {
				handler(w, r)
				continue
			}
			w := &responseWriter{w, s}
			handler(w, r)
			s.Close()
			s = s.Next()
		}
	})
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusOK)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "<layout><head><title>some title</title><meta name='description' content='some description'/></head><frame2><frame1><view></view></frame1></frame2></layout>")
}

func TestBatch(t *testing.T) {
	is := is.New(t)
	handlers := []func(w http.ResponseWriter, r *http.Request){
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			w.Write([]byte("<view>"))
			w.Write(inner)
			w.Write([]byte("</view>"))
		},
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			w.Write([]byte("<frame1>"))
			w.Write(inner)
			w.Write([]byte("</frame1>"))
		},
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			w.Write([]byte("<frame2>"))
			w.Write(inner)
			w.Write([]byte("</frame2>"))
		},
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			w.Write([]byte("<layout>"))
			w.Write(inner)
			w.Write([]byte("</layout>"))
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg := new(sync.WaitGroup)
		s := slots.New()
		for i, handler := range handlers {
			i := i
			handler := handler
			r := r.Clone(slots.ToContext(r.Context(), s))
			if i == len(handlers)-1 {
				wg.Add(1)
				go func(slot *slots.Slots) {
					defer wg.Done()
					defer slot.Close()
					handler(w, r)
				}(s)
				continue
			}
			w := &responseWriter{w, s}
			wg.Add(1)
			go func(slot *slots.Slots) {
				defer wg.Done()
				defer slot.Close()
				handler(w, r)
			}(s)
			s = s.Next()
		}
		wg.Wait()
	})
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusOK)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "<layout><frame2><frame1><view></view></frame1></frame2></layout>")
}

func TestBatchNestedSlots(t *testing.T) {
	is := is.New(t)
	handlers := []func(w http.ResponseWriter, r *http.Request){
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			headSlot := slots.Slot("head")
			headSlot.Write([]byte("<title>some title</title>"))
			w.Write([]byte("<view>"))
			w.Write(inner)
			w.Write([]byte("</view>"))
		},
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			w.Write([]byte("<frame1>"))
			w.Write(inner)
			w.Write([]byte("</frame1>"))
		},
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			headSlot := slots.Slot("head")
			headSlot.Write([]byte("<meta name='description' content='some description'/>"))
			w.Write([]byte("<frame2>"))
			w.Write(inner)
			w.Write([]byte("</frame2>"))
		},
		func(w http.ResponseWriter, r *http.Request) {
			slots, err := slots.FromContext(r.Context())
			is.NoErr(err)
			inner, err := io.ReadAll(slots)
			is.NoErr(err)
			headSlot := slots.Slot("head")
			head, err := io.ReadAll(headSlot)
			is.NoErr(err)
			w.Write([]byte("<layout>"))
			w.Write([]byte("<head>"))
			w.Write(head)
			w.Write([]byte("</head>"))
			w.Write(inner)
			w.Write([]byte("</layout>"))
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg := new(sync.WaitGroup)
		s := slots.New()
		for i, handler := range handlers {
			i := i
			handler := handler
			r := r.Clone(slots.ToContext(r.Context(), s))
			if i == len(handlers)-1 {
				wg.Add(1)
				go func(slot *slots.Slots) {
					defer wg.Done()
					defer slot.Close()
					handler(w, r)
				}(s)
				continue
			}
			w := &responseWriter{w, s}
			wg.Add(1)
			go func(slot *slots.Slots) {
				defer wg.Done()
				defer slot.Close()
				handler(w, r)
			}(s)
			s = s.Next()
		}
		wg.Wait()
	})
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusOK)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "<layout><head><title>some title</title><meta name='description' content='some description'/></head><frame2><frame1><view></view></frame1></frame2></layout>")
}
