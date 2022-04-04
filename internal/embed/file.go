package embed

import "strings"

type File struct {
	Path string
	Data Data
}

type Data []byte

const lowerHex = "0123456789abcdef"

// Based on:
// https://github.com/go-bindata/go-bindata/blob/26949cc13d95310ffcc491c325da869a5aafce8f/stringwriter.go#L18-L36
func (data Data) String() string {
	if len(data) == 0 {
		return ""
	}
	s := new(strings.Builder)
	buf := []byte(`\x00`)
	for _, b := range data {
		buf[2] = lowerHex[b/16]
		buf[3] = lowerHex[b%16]
		s.Write(buf)
	}
	return s.String()
}
