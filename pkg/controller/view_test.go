package controller_test

// func viewEqual(t testing.TB, h http.Handler, request, expect string) {
// 	t.Helper()
// 	rec := httptest.NewRecorder()
// 	parts := strings.Split(request, " ")
// 	if len(parts) != 2 {
// 		t.Fatalf("invalid request: %s", request)
// 	}
// 	u, err := url.Parse(parts[1])
// 	if err != nil {
// 		t.Fatalf("invalid request: %s", request)
// 	}
// 	req := httptest.NewRequest(parts[0], u.Path, nil)
// 	req.URL.RawQuery = u.RawQuery
// 	h.ServeHTTP(rec, req)
// 	res := rec.Result()
// 	dump, err := httputil.DumpResponse(res, true)
// 	if err != nil {
// 		if err.Error() != expect {
// 			t.Fatalf("unexpected error: %v", err)
// 		}
// 		return
// 	}
// 	diff.TestHTTP(t, expect, string(dump))
// }

// func TestHtmlErrorOk(t *testing.T) {
// 	is := is.New(t)
// 	fsys := fstest.MapFS{
// 		"index.html": &fstest.MapFile{Data: []byte(`All good`)},
// 		"error.html": &fstest.MapFile{Data: []byte(`Error: {{.Error}}`)},
// 	}
// 	vf := view.New(fsys, map[string]view.Renderer{
// 		".html": gohtml.New(),
// 	})
// 	page, err := vf.Find("index")
// 	is.NoErr(err)
// 	h, err := rpc.View(page, func() error {
// 		return nil
// 	})
// 	is.NoErr(err)
// 	viewEqual(t, h, "GET /", `
// 		HTTP/1.1 200 OK
// 		Connection: close
// 		Content-Type: text/html

// 		All good
// 	`)
// }

// func TestHtmlErrorNotOk(t *testing.T) {
// 	is := is.New(t)
// 	fsys := fstest.MapFS{
// 		"index.html": &fstest.MapFile{Data: []byte(`All good`)},
// 		"error.html": &fstest.MapFile{Data: []byte(`Error: {{.Error}}`)},
// 	}
// 	vf := view.New(fsys, map[string]view.Renderer{
// 		".html": gohtml.New(),
// 	})
// 	page, err := vf.Find("index")
// 	is.NoErr(err)
// 	h, err := rpc.View(page, func() error {
// 		return errors.New("Oh noz")
// 	})
// 	is.NoErr(err)
// 	viewEqual(t, h, "GET /", `
// 		HTTP/1.1 500 Internal Server Error
// 		Connection: close
// 		Content-Type: text/html

// 		Error: Oh noz
// 	`)
// }
