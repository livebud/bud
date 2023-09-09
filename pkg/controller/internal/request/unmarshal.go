package request

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"

	"github.com/ajg/form"
)

// Unmarshal the request data into v
func Unmarshal(r *http.Request, v interface{}) error {
	err := unmarshalBody(r, v)
	if err != nil {
		return err
	}
	err = unmarshalURL(r.URL, v)
	if err != nil {
		return err
	}
	return nil
}

func unmarshalBody(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return nil
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}
	switch mediaType {
	case "application/json":
		return unmarshalJSON(r.Body, v)
	case "application/x-www-form-urlencoded":
		return unmarshalForm(r, v)
	}
	return nil
}

func unmarshalURL(u *url.URL, v interface{}) error {
	dec := form.NewDecoder(nil)
	dec.IgnoreCase(true)
	dec.IgnoreUnknownKeys(true)
	return dec.DecodeValues(v, u.Query())
}

func unmarshalForm(r *http.Request, v interface{}) error {
	if r.PostForm == nil {
		r.ParseForm()
	}
	dec := form.NewDecoder(nil)
	dec.IgnoreCase(true)
	dec.IgnoreUnknownKeys(true)
	return dec.DecodeValues(v, r.PostForm)
}

func unmarshalJSON(r io.Reader, v interface{}) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}
