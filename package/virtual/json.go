package virtual

import (
	"encoding/json"
	"io/fs"
	"time"
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
	Mode    fs.FileMode
	ModTime time.Time
	Sys     interface{}
	Entries []*DirEntry
}

func (f *jsonEntry) Open() fs.File {
	if f.Mode.IsDir() {
		entries := make([]fs.DirEntry, len(f.Entries))
		for i, entry := range f.Entries {
			entries[i] = entry
		}
		return &openDir{&File{
			Path:    f.Path,
			Mode:    f.Mode,
			ModTime: f.ModTime,
			Entries: entries,
		}, 0}
	}
	return &openFile{&File{
		Path:    f.Path,
		Data:    f.Data,
		Mode:    f.Mode,
		ModTime: f.ModTime,
	}, 0}
}

func UnmarshalJSON(file []byte) (fs.File, error) {
	var entry jsonEntry
	err := json.Unmarshal(file, &entry)
	if err != nil {
		return nil, err
	}
	return entry.Open(), nil
}
