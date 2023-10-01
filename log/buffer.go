package log

func Buffer() *bufferHandler {
	return &bufferHandler{}
}

type bufferHandler struct {
	Entries []*Entry
}

var _ Handler = (*bufferHandler)(nil)

func (h *bufferHandler) Log(entry *Entry) error {
	h.Entries = append(h.Entries, entry)
	return nil
}
