package token

import (
	"strconv"
	"strings"
)

type Type string

type Token struct {
	Type  Type
	Text  string
	Start int
	Line  int
}

func (t *Token) String() string {
	s := new(strings.Builder)
	s.WriteString(string(t.Type))
	if t.Text != "" && t.Text != string(t.Type) {
		s.WriteString(":")
		s.WriteString(strconv.Quote(t.Text))
	}
	return s.String()
}

const (
	End        Type = "end"
	Error      Type = "error"
	Regexp     Type = "regexp"
	Path       Type = "path"
	Slot       Type = "slot"
	Slash      Type = "/"
	OpenCurly  Type = "{"
	CloseCurly Type = "}"
	Question   Type = "?"
	Star       Type = "*"
	Pipe       Type = "|"
)
