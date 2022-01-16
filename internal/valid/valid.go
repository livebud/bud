package valid

import (
	"path"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Dir validates that the name matches a valid directory
func Dir(name string) bool {
	return !invalidDir(name)
}

// Invalid dir check
func invalidDir(name string) bool {
	return len(name) == 0 || // Empty string
		name[0] == '_' || // Starts with _
		name[0] == '.' || // Starts with .
		name == "bud" || // Named bud (reserved)
		strings.ToLower(name) != name // Has uppercase letters
}

// PluginDir validates that the name is a valid plugin directory
func PluginDir(name string) bool {
	return len(name) == 0 || // Empty string
		name[0] == '_' || // Starts with _
		name[0] == '.' || // Starts with .
		strings.ToLower(name) != name // Has uppercase letters
}

// ViewEntry validates that name matches a valid view entrypoint
func ViewEntry(name string) bool {
	return !invalidViewEntry(name)
}

// Invalid view entry check
func invalidViewEntry(name string) bool {
	return len(name) == 0 || // Empty string
		name[0] == '_' || // Starts with _
		name[0] == '.' || // Starts with .
		name == "bud" || // Named bud (reserved)
		unicode.IsUpper(firstRune(name)) // Starts with a capital letter
}

func ActionFile(name string) bool {
	return !invalidActionFile(name)
}

func invalidActionFile(name string) bool {
	return len(name) == 0 || // Empty string
		path.Ext(name) != ".go" ||
		name[0] == '_' || // Starts with _
		name[0] == '.' || // Starts with .
		name == "bud.go" // Named bud (reserved)
}

func firstRune(s string) rune {
	r, _ := utf8.DecodeRuneInString(s)
	return r
}
