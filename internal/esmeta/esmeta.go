package esmeta

import (
	"encoding/json"
	"strings"
)

func Parse(metafile string) (*File, error) {
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
		for input := range output.Inputs {
			idx := strings.IndexByte(input, ':')
			if idx >= 0 {
				continue
			}
			deps = append(deps, input)
		}
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
