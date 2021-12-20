package valid

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Dir validates that the name matches a valid directory
func Dir(name string) bool {
	// Ignore dirs with capital letters and files that start with an underscore or dot
	if len(name) == 0 || name[0] == '_' || name[0] == '.' || strings.ToLower(name) != name {
		return false
	}
	// Ignore _
	if name[0] == '_' {
		return false
	}
	return true
}

// ViewEntry validates that name matches a valid view entrypoint
func ViewEntry(name string) bool {
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
