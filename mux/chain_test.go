package mux_test

// func TestChain(t *testing.T) {
// 	is := is.New(t)
// 	handler := mux.Chain(
// 		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.Header().Set("a", "aa")
// 			body, err := io.ReadAll(r.Body)
// 			is.NoErr(err)
// 			w.Write([]byte("<a>" + string(body) + "</a>"))
// 		}),
// 		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.Header().Set("b", "bb")
// 			body, err := io.ReadAll(r.Body)
// 			is.NoErr(err)
// 			w.Write([]byte("<b>" + string(body) + "</b>"))
// 		}),
// 		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.Header().Set("b", "cc")
// 			body, err := io.ReadAll(r.Body)
// 			is.NoErr(err)
// 			w.Write([]byte("<c>" + string(body) + "</c>"))
// 		}),
// 	)
// 	r := httptest.NewRequest("GET", "/", nil)
// 	w := httptest.NewRecorder()
// 	handler.ServeHTTP(w, r)
// 	res := w.Result()
// 	is.Equal(res.StatusCode, http.StatusOK)
// 	is.Equal(res.Header.Get("a"), "aa")
// 	is.Equal(res.Header.Get("b"), "bb")
// 	body, err := io.ReadAll(res.Body)
// 	is.NoErr(err)
// 	is.Equal(string(body), "<c><b><a></a></b></c>")
// }
