package virtual

import (
	"encoding/json"
	"io/fs"
)

func MarshalJSON(file fs.File) ([]byte, error) {
	entry, err := From(file)
	if err != nil {
		return nil, err
	}
	return json.Marshal(entry)
}

type jsonEntry struct {
	Path    string
	Data    []byte
	Entries []*DirEntry
}

func (f *jsonEntry) Open() fs.File {
	if f.Entries != nil {
		entries := make([]fs.DirEntry, len(f.Entries))
		for i, entry := range f.Entries {
			entries[i] = entry
		}
		return &Dir{
			Path:    f.Path,
			Entries: entries,
		}
	}
	return &File{
		Path: f.Path,
		Data: f.Data,
	}
}

func UnmarshalJSON(file []byte) (fs.File, error) {
	var entry jsonEntry
	err := json.Unmarshal(file, &entry)
	if err != nil {
		return nil, err
	}
	return entry.Open(), nil
}
