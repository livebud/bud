package esmeta

import (
	"encoding/json"
	"path/filepath"
	"strings"
)

func Load(metafile string) (*File, error) {
	file := new(File)
	if err := json.Unmarshal([]byte(metafile), file); err != nil {
		return nil, err
	}
	return file, nil
}

type File struct {
	Inputs  map[string]*Input  `json:"inputs,omitempty"`
	Outputs map[string]*Output `json:"outputs,omitempty"`
}

func (f *File) Dependencies() (deps []string) {
	for _, output := range f.Outputs {
		_ = output
	}
	return deps
}

type Input struct {
	Bytes   int       `json:"bytes,omitempty"`
	Imports []*Import `json:"imports,omitempty"`
}

type Output struct {
	Bytes      int                     `json:"bytes,omitempty"`
	Inputs     map[string]*OutputInput `json:"inputs,omitempty"`
	Imports    []*Import               `json:"imports,omitempty"`
	Exports    []string                `json:"exports,omitempty"`
	EntryPoint *string                 `json:"entryPoint,omitempty"`
}

type OutputInput struct {
	BytesInOutput int `json:"bytesInOutput,omitempty"`
}

type Import struct {
	Path string `json:"path,omitempty"`
	Kind string `json:"kind,omitempty"`
}

func Clean(s string) string {
	// Trim off any namespaces
	idx := strings.IndexByte(s, ':')
	if idx >= 0 {
		s = s[idx+1:]
	}
	return filepath.Clean(s)
}

// type linkable interface {
// 	Link(from, to string)
// }

// // Link the inputs to the outputs. Link isn't a representation of which file
// // required what, it's just saying that this file contributed to this output.
// func Link(l linkable, from, metafile string) error {
// 	file, err := Load(metafile)
// 	if err != nil {
// 		return err
// 	}
// 	for _, output := range file.Outputs {
// 		for to := range output.Inputs {
// 			// Ignore virtual files
// 			if strings.HasSuffix(to, "!") {
// 				continue
// 			}
// 			// Trim off any namespaces
// 			to = trimNamespace(to)
// 			// Remove ./
// 			to = filepath.Clean(to)
// 			// Link the file
// 			l.Link(from, to)
// 		}
// 	}
// 	return nil
// }

// func Link2(l linkable, metafile string) error {
// 	file, err := Load(metafile)
// 	if err != nil {
// 		return err
// 	}
// 	for from, input := range file.Inputs {
// 		from = Clean(from)
// 		// fmt.Println(from)
// 		for _, im := range input.Imports {
// 			// fmt.Println("  ", im.Path)
// 			l.Link(from, Clean(im.Path))
// 		}
// 	}
// 	for from, output := range file.Outputs {
// 		from = Clean(from)
// 		// fmt.Println(from, output.Exports)
// 		for to := range output.Inputs {
// 			l.Link(from, Clean(to))
// 		}
// 	}
// 	return nil
// }

// func trimNamespace(s string) string {
// 	// Trim off any namespaces
// 	idx := strings.IndexByte(s, ':')
// 	if idx >= 0 {
// 		s = s[idx+1:]
// 	}
// 	return s
// }
