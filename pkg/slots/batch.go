package slots

import "net/http"

// Batch a list of handlers to be called in parallel.
func Batch(handlers ...http.Handler) http.Handler {
	// TODO: finish me
	return handlers[0]
}

// func Batch(handlers ...http.Handler) http.Handler {
// 	if len(handlers) == 1 {
// 		return handlers[0]
// 	}
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// slot := New()
// 		// eg := new(errgroup.Group)
// 		// status := 0
// 		// headers := w.Header()
// 		// mu := sync.Mutex{}
// 		// for i := 0; i < len(handlers); i++ {
// 		// 	handler := handlers[i]
// 		// 	r := request.Clone(r)
// 		// 	innerSlot := slot
// 		// 	eg.Go(func() (err error) {
// 		// 		defer func() { err = innerSlot.Close() }()
// 		// 		innerHeaders := http.Header{}
// 		// 		w := httpsnoop.Wrap(w, httpsnoop.Hooks{
// 		// 			WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
// 		// 				return func(code int) {
// 		// 					mu.Lock()
// 		// 					status = code
// 		// 					mu.Unlock()
// 		// 				}
// 		// 			},
// 		// 			Write: func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
// 		// 				return innerSlot.Write
// 		// 			},
// 		// 			Header: func(next httpsnoop.HeaderFunc) httpsnoop.HeaderFunc {
// 		// 				return func() http.Header {
// 		// 					return innerHeaders
// 		// 				}
// 		// 			},
// 		// 		})
// 		// 		handler.ServeHTTP(w, r)
// 		// 		mu.Lock()
// 		// 		for key := range innerHeaders {
// 		// 			headers.Set(key, innerHeaders.Get(key))
// 		// 		}
// 		// 		mu.Unlock()
// 		// 		return err
// 		// 	})
// 		// 	innerSlot = slot.Next()
// 		// }
// 		// if err := eg.Wait(); err != nil {
// 		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 		// 	return
// 		// }

// 		// state, err := readState(pipeline)
// 		// if err != nil {
// 		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 		// 	return
// 		// }
// 		// w.Write([]byte(state.Data))
// 	})
// }

// // func newResponseWriter(w http.ResponseWriter, slot io.Writer) http.ResponseWriter {

// // 	return &responseWriter{w, slot, w.Header(), 200}
// // }

// // type responseWriter struct {
// // 	http.ResponseWriter
// // 	slot   io.Writer
// // 	header http.Header
// // 	status int
// // }

// // func (w *responseWriter) Write(p []byte) (int, error) {
// // 	return w.slot.Write(p)
// // }

// // func (w *responseWriter) Header() http.Header {
// // 	return w.header
// // }

// // func (w *responseWriter) WriteHeader(status int) {
// // 	w.status = status
// // }
