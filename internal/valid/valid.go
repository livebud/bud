package valid

import (
	"unicode"
	"unicode/utf8"
)

func Dir(name string) bool {
	// Ignore capitalized files and files that start with an underscore or dot
	if len(name) == 0 || name[0] == '_' || name[0] == '.' || unicode.IsUpper(firstRune(name)) {
		return false
	}
	// Ignore _
	if name[0] == '_' {
		return false
	}
	return true
}

func firstRune(s string) rune {
	r, _ := utf8.DecodeRuneInString(s)
	return r
}
