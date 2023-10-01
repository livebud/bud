package log

import (
	"encoding/json"
	"io"
)

func Json(w io.Writer) Handler {
	return &jsonHandler{json.NewEncoder(w)}
}

type jsonHandler struct {
	enc *json.Encoder
}

func (h *jsonHandler) Log(entry *Entry) error {
	return h.enc.Encode(entry)
}
